package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"hmdp-go/internal/dto"
	"hmdp-go/internal/model"
	"hmdp-go/internal/pkg/constants"
	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/repository"

	"errors"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// BlogService 定义“博客/探店笔记”相关业务能力。
type BlogService interface {
	// SaveBlog 发布博客。
	SaveBlog(ctx context.Context, blog model.Blog, currentUserID int64) result.Result
	// QueryByID 查询博客详情。
	QueryByID(ctx context.Context, id int64, userID int64) result.Result
	// LikeBlog 点赞或取消点赞。
	LikeBlog(ctx context.Context, id int64, userID int64) result.Result
	// QueryHotBlog 查询首页热门博客。
	QueryHotBlog(ctx context.Context, current int, userID int64) result.Result
	// QueryBlogLikes 查询一篇博客的点赞用户列表。
	QueryBlogLikes(ctx context.Context, id int64) result.Result
	// QueryBlogByUserID 查询指定用户发布的博客。
	QueryBlogByUserID(ctx context.Context, authorID int64, current int, viewerID int64) result.Result
	//增加查询关注流
	QueryBlogOfFollow(ctx context.Context, max int64, offset int, currentUserID int64) result.Result
}

type blogService struct {
	// blogRepo 负责 tb_blog 的数据库操作。
	blogRepo repository.BlogRepository
	// userRepo 用来查作者昵称、头像。
	userRepo   repository.UserRepository
	followRepo repository.FollowRepository
	// redisClient 后续点赞、Feed 流会用 Redis。
	redisClient *redis.Client
}

// NewBlogService 创建博客 Service。
func NewBlogService(blogRepo repository.BlogRepository, userRepo repository.UserRepository, followRepo repository.FollowRepository, redisClient *redis.Client) BlogService {
	return &blogService{blogRepo: blogRepo, userRepo: userRepo, followRepo: followRepo, redisClient: redisClient}
}

// SaveBlog 发布博客。
//
// 后面要做的事情：
// 1. 从登录上下文拿当前用户 id；
// 2. 保存博客到 tb_blog；
// 3. 把博客推送到粉丝的 Feed 流。
func (s *blogService) SaveBlog(ctx context.Context, blog model.Blog, currentUserID int64) result.Result {
	// 1. 补全博客信息
	blog.UserID = currentUserID

	// 2. 保存博客到 MySQL 数据库
	if err := s.blogRepo.SaveBlog(ctx, &blog); err != nil {
		log.Printf("保存博客失败: %v", err)
		return result.Fail("发布博客失败")
	}
	// 注意：此时通过 GORM 的特性，blog.ID 已经被自动赋上了数据库生成的自增 ID
	// 3. 查询当前用户的所有粉丝 ID 列表
	// 假设你的 followRepo 提供了批量查粉丝 ID 的方法，返回 []int64
	followerIDs, err := s.followRepo.FindFollowerIDs(ctx, currentUserID)
	if err != nil {
		log.Printf("获取粉丝列表失败: %v", err)
		// 核心数据已落盘，为了用户体验，学习阶段可以打印日志并继续放行，或者提示成功
		return result.OKWithData(blog.ID)
	}
	if len(followerIDs) == 0 {
		return result.OKWithData(blog.ID)
	}

	now := time.Now().UnixMilli()
	blogIDStr := strconv.FormatInt(blog.ID, 10)
	pipe := s.redisClient.Pipeline()
	for _, followerID := range followerIDs {
		key := constants.FeedKey + strconv.FormatInt(followerID, 10)
		pipe.ZAdd(ctx, key, redis.Z{
			Score:  float64(now),
			Member: blogIDStr,
		})
	}

	if _, err := pipe.Exec(ctx); err != nil {
		log.Printf("推送 Feed 失败: %v", err)
	}
	return result.OKWithData(blog.ID)
}

// QueryByID 查询博客详情。
//
// 后面要补充作者信息、当前用户是否点赞等字段。
func (s *blogService) QueryByID(ctx context.Context, blogId int64, userid int64) result.Result {
	if blogId <= 0 {
		return result.Fail("invalid blog id")
	}
	// 2. 查博客表
	blog, err := s.blogRepo.FindBlogByID(ctx, blogId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return result.Fail("blog not found")
	}
	if err != nil {
		return result.Fail("query blog failed")
	}

	// 3. 查用户表 (使用 blog.UserID，并且传入 ctx)
	user, err := s.userRepo.FindUserByID(ctx, blog.UserID)
	if err == nil {
		blog.Name = user.NickName
		blog.Icon = user.Icon
	}

	// 4.复用工具方法！
	// 完美解锁“游客不查 Redis 节省 IO”和“统一异常捕获”双重隐藏成就
	s.setBlogIsLike(ctx, blog, userid)

	return result.OKWithData(blog)
}

// LikeBlog 点赞/取消点赞。
//
// Java 版使用 Redis ZSet 记录点赞用户，并同步更新数据库 liked 数量。
func (s *blogService) LikeBlog(ctx context.Context, blogId int64, userId int64) result.Result {
	if blogId <= 0 {
		return result.Fail("无效的博客记录")
	}
	if userId <= 0 {
		return result.Fail("无效的用户状态，请重新登录")
	}
	// 1. 拼接 Redis Key: "blog:liked:10"
	// 注意：Go 里面字符串不能直接加整型，必须用 strconv 转换
	key := constants.BlogLikedKey + strconv.FormatInt(blogId, 10)
	member := strconv.FormatInt(userId, 10)

	// 2. 检查该用户是否已经点赞过
	// ZScore 命令：如果 member 存在，返回 score；如果不存在，返回 redis.Nil 错误。时间复杂度 O(1)，极快！
	_, err := s.redisClient.ZScore(ctx, key, member).Result()

	if errors.Is(err, redis.Nil) {
		// 场景 A：未点赞 -> 执行【点赞】逻辑
		// 3.1 数据库点赞数 +1 ( 这里要调用Repository，下面给出了原子操作的写法)
		updateErr := s.blogRepo.UpdateBlogLiked(ctx, blogId, 1)
		if updateErr != nil {
			return result.Fail("点赞失败，数据库异常")
		}

		// 3.2 将用户加入 Redis ZSet，Score 存当前时间戳 (为了以后按点赞时间排序)
		//  注意：提取出 Err() 来进行严格校验
		zaddErr := s.redisClient.ZAdd(ctx, key, redis.Z{
			Score:  float64(time.Now().UnixMilli()),
			Member: member,
		}).Err()

		// 3.3 终极防线：如果 Redis 失败，必须回滚 MySQL！
		if zaddErr != nil {
			log.Printf("Redis ZAdd 失败，准备回滚 MySQL: %v", zaddErr)

			// 补偿操作：把刚才加的 1 减回去
			rollbackErr := s.blogRepo.UpdateBlogLiked(ctx, blogId, -1)
			if rollbackErr != nil {
				// 极端情况：Redis 挂了，回滚时 MySQL 也挂了（发生概率极低，通常需人工介入）
				log.Printf("严重一致性灾难:MySQL 回滚点赞数失败! blogId: %d", blogId)
			}
			return result.Fail("点赞失败，请稍后重试")
		}

		return result.OKWithData("点赞成功")

	} else if err == nil {
		// 场景 B：已点赞 -> 执行【取消点赞】逻辑
		// 4.1 数据库点赞数 -1
		updateErr := s.blogRepo.UpdateBlogLiked(ctx, blogId, -1)
		if updateErr != nil {
			return result.Fail("取消点赞失败，数据库异常")
		}

		// 4.2 将用户从 Redis ZSet 中移除
		zremErr := s.redisClient.ZRem(ctx, key, member).Err()

		// 4.3 终极防线：如果 Redis 失败，必须回滚 MySQL！
		if zremErr != nil {
			log.Printf("Redis ZRem 失败，准备回滚 MySQL: %v", zremErr)

			// 补偿操作：把刚才减的 1 加回来
			rollbackErr := s.blogRepo.UpdateBlogLiked(ctx, blogId, 1)
			if rollbackErr != nil {
				log.Printf("严重一致性灾难：MySQL 回滚取消点赞数失败! blogId: %d", blogId)
			}
			return result.Fail("取消点赞失败，请稍后重试")
		}

		return result.OKWithData("取消点赞成功")

	} else {
		// 场景 C：Redis 真的挂了或者网络抖动
		log.Printf(" 查询 Redis ZScore 异常: %v", err)
		return result.Fail("系统异常，请稍后重试")
	}
}

// QueryHotBlog 查询首页热门博客。
//
// 这一步不只是查 tb_blog：
// 1. 先按点赞数查出博客列表；
// 2. 再根据每篇博客的 user_id 查作者；
// 3. 把作者昵称 name、头像 icon 填回 Blog；
// 4. 返回给前端渲染首页列表。
func (s *blogService) QueryHotBlog(ctx context.Context, current int, userID int64) result.Result {
	blogs, err := s.blogRepo.FindBlogsByHot(ctx, current)
	if err != nil {
		return result.Fail("query hot blog failed")
	}

	// 用下标 i 遍历，是因为我们要修改 blogs[i] 本身。
	// 如果写成 for _, blog := range blogs，blog 是副本，改了不会影响原切片。
	for i := range blogs {
		user, err := s.userRepo.FindUserByID(ctx, blogs[i].UserID)
		if err == nil {
			// Name/Icon 是 Blog 结构体里的 gorm:"-" 字段，不存在数据库中，只用于返回给前端。
			blogs[i].Name = user.NickName
			blogs[i].Icon = user.Icon
		}

		// 复用工具方法，判断点赞状态,实现点赞状态同步
		s.setBlogIsLike(ctx, &blogs[i], userID)

	}
	return result.OKWithData(blogs)
}

// QueryBlogLikes 查询点赞用户列表。
func (s *blogService) QueryBlogLikes(ctx context.Context, blogId int64) result.Result {
	// 拼接 Redis Key: "blog:liked:10"
	// 注意：Go 里面字符串不能直接加整型，必须用 strconv 转换
	key := constants.BlogLikedKey + strconv.FormatInt(blogId, 10)

	// 1. 查询 Top 5 的点赞用户 (利用 ZSet 默认按 Score 时间戳升序，取出前 5 个)
	// 等价于 Redis 命令：ZRANGE blog:liked:5 0 4
	top5, err := s.redisClient.ZRange(ctx, key, 0, 4).Result()
	if err != nil {
		log.Printf("Redis 命令:ZRANGE出错 ")
		return result.Fail("redis error")
	}
	if len(top5) == 0 {
		// 没人点赞，返回空列表，不要返回 nil，防止前端崩溃
		return result.OKWithData(make([]dto.UserDTO, 0))
	}

	// 2. 根据取出的用户 ID 去数据库查询对应的昵称和头像
	var userDTOs []dto.UserDTO
	for _, idStr := range top5 {
		uid, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			// 记录报警日志，但不阻断整个列表的渲染，直接跳过这个脏数据
			log.Printf("警告: Redis 中存在无法解析的点赞用户 ID [%s]: %v", idStr, err)
			continue
		}
		user, err := s.userRepo.FindUserByID(ctx, uid)
		if err == nil {
			userDTOs = append(userDTOs, dto.UserDTO{
				ID:       user.ID,
				NickName: user.NickName,
				Icon:     user.Icon,
			})
		}
	}

	// 3. 返回给前端渲染头像列表
	return result.OKWithData(userDTOs)
}

// QueryBlogByUserID 查询指定用户发布的博客列表。
func (s *blogService) QueryBlogByUserID(ctx context.Context, authorID int64, current int, viewerID int64) result.Result {
	if authorID <= 0 {
		return result.Fail("invalid authorID")
	}

	// 1. 使用 authorID 查出主页主人的所有博客
	blogs, err := s.blogRepo.FindBlogsByUserID(ctx, authorID, current)
	if err != nil {
		return result.Fail("query userblogs failed")
	}
	if len(blogs) == 0 {
		return result.OKWithData(make([]model.Blog, 0))
	}

	// 2. 遍历博客，检查点赞状态
	for i := range blogs {
		// 核心修复：这里传入 viewerID！
		// 判断的是“当前屏幕前的用户”，有没有给这些篇博客点过赞
		s.setBlogIsLike(ctx, &blogs[i], viewerID)
	}
	return result.OKWithData(blogs)
}

// 抽离“点赞状态判断”工具
// setBlogIsLike 判断当前用户是否点赞了该博客，并给 IsLike 赋值
func (s *blogService) setBlogIsLike(ctx context.Context, blog *model.Blog, userID int64) {
	if userID == 0 {
		// 游客未登录，直接显示未点赞
		blog.IsLike = false
		return
	}

	// 拼接 Redis Key
	key := constants.BlogLikedKey + strconv.FormatInt(blog.ID, 10)
	member := strconv.FormatInt(userID, 10)

	// 使用 ZScore 查询，如果 err == nil 说明查到了分数，代表已点赞
	_, err := s.redisClient.ZScore(ctx, key, member).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Printf("Redis ZScore查询 error!")
	}
	blog.IsLike = err == nil
}

func (s *blogService) QueryBlogOfFollow(ctx context.Context, max int64, offset int, currentUserID int64) result.Result {
	if currentUserID <= 0 {
		return result.Fail("用户未登录")
	}
	if max <= 0 {
		max = time.Now().UnixMilli()
	}
	if offset < 0 {
		offset = 0
	}

	key := constants.FeedKey + strconv.FormatInt(currentUserID, 10)

	zs, err := s.redisClient.ZRevRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
		Max:    strconv.FormatInt(max, 10),
		Min:    "0",
		Offset: int64(offset),
		Count:  2,
	}).Result()
	if err != nil {
		return result.Fail("查询关注流失败")
	}

	if len(zs) == 0 {
		return result.OKWithData(dto.ScrollResult{
			List:    []model.Blog{},
			MinTime: 0,
			Offset:  0,
		})
	}

	ids := make([]int64, 0, len(zs))
	minTime := int64(0)
	nextOffset := 0

	for _, z := range zs {
		score := int64(z.Score)

		if score == minTime {
			nextOffset++
		} else {
			minTime = score
			nextOffset = 1
		}

		idStr := fmt.Sprint(z.Member)
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Printf("Feed 中存在非法 blogId: %v", z.Member)
			continue
		}
		ids = append(ids, id)
	}

	blogs, err := s.blogRepo.FindBlogsByIDs(ctx, ids)
	if err != nil {
		return result.Fail("查询博客失败")
	}

	blogMap := make(map[int64]model.Blog, len(blogs))
	for _, blog := range blogs {
		blogMap[blog.ID] = blog
	}

	orderedBlogs := make([]model.Blog, 0, len(ids))
	for _, id := range ids {
		blog, ok := blogMap[id]
		if !ok {
			continue
		}
		orderedBlogs = append(orderedBlogs, blog)
	}

	for i := range orderedBlogs {
		user, err := s.userRepo.FindUserByID(ctx, orderedBlogs[i].UserID)
		if err == nil {
			orderedBlogs[i].Name = user.NickName
			orderedBlogs[i].Icon = user.Icon
		}
		s.setBlogIsLike(ctx, &orderedBlogs[i], currentUserID)
	}

	return result.OKWithData(dto.ScrollResult{
		List:    orderedBlogs,
		MinTime: minTime,
		Offset:  nextOffset,
	})
}
