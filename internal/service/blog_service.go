package service

import (
	"context"

	"hmdp-go/internal/model"
	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/repository"

	"errors"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// BlogService 定义“博客/探店笔记”相关业务能力。
type BlogService interface {
	// SaveBlog 发布博客。
	SaveBlog(ctx context.Context, blog model.Blog) result.Result
	// QueryByID 查询博客详情。
	QueryByID(ctx context.Context, id int64) result.Result
	// LikeBlog 点赞或取消点赞。
	LikeBlog(ctx context.Context, id int64) result.Result
	// QueryMyBlog 查询当前登录用户自己的博客。
	QueryMyBlog(ctx context.Context, current int) result.Result
	// QueryHotBlog 查询首页热门博客。
	QueryHotBlog(ctx context.Context, current int) result.Result
	// QueryBlogLikes 查询一篇博客的点赞用户列表。
	QueryBlogLikes(ctx context.Context, id int64) result.Result
	// QueryBlogByUserID 查询指定用户发布的博客。
	QueryBlogByUserID(ctx context.Context, userID int64, current int) result.Result
}

type blogService struct {
	// blogRepo 负责 tb_blog 的数据库操作。
	blogRepo repository.BlogRepository
	// userRepo 用来查作者昵称、头像。
	userRepo repository.UserRepository
	// redisClient 后续点赞、Feed 流会用 Redis。
	redisClient *redis.Client
}

// NewBlogService 创建博客 Service。
func NewBlogService(blogRepo repository.BlogRepository, userRepo repository.UserRepository, redisClient *redis.Client) BlogService {
	return &blogService{blogRepo: blogRepo, userRepo: userRepo, redisClient: redisClient}
}

// SaveBlog 发布博客。
//
// 后面要做的事情：
// 1. 从登录上下文拿当前用户 id；
// 2. 保存博客到 tb_blog；
// 3. 把博客推送到粉丝的 Feed 流。
func (s *blogService) SaveBlog(ctx context.Context, blog model.Blog) result.Result {
	// TODO: Get current user id from context, save blog, and push feed to followers.
	return result.Fail("TODO: save blog")
}

// QueryByID 查询博客详情。
//
// 后面要补充作者信息、当前用户是否点赞等字段。
func (s *blogService) QueryByID(ctx context.Context, id int64) result.Result {
	if id <= 0 {
		return result.Fail("invalid blog id")
	}
	// 2. 查博客表
	blog, err := s.blogRepo.FindBlogByID(ctx, id)
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
	blog.IsLike = false
	return result.OKWithData(blog)
}

// LikeBlog 点赞/取消点赞。
//
// Java 版使用 Redis ZSet 记录点赞用户，并同步更新数据库 liked 数量。
func (s *blogService) LikeBlog(ctx context.Context, id int64) result.Result {
	// TODO: Toggle like status with blog:liked:{id} sorted set and update liked count.
	return result.Fail("TODO: like blog")
}

// QueryMyBlog 查询当前登录用户发布的博客。
func (s *blogService) QueryMyBlog(ctx context.Context, current int) result.Result {
	// TODO: Query current user's blogs.
	return result.Fail("TODO: query my blog")
}

// QueryHotBlog 查询首页热门博客。
//
// 这一步不只是查 tb_blog：
// 1. 先按点赞数查出博客列表；
// 2. 再根据每篇博客的 user_id 查作者；
// 3. 把作者昵称 name、头像 icon 填回 Blog；
// 4. 返回给前端渲染首页列表。
func (s *blogService) QueryHotBlog(ctx context.Context, current int) result.Result {
	blogs, err := s.blogRepo.FindBlogsByHot(ctx, current)
	if err != nil {
		return result.Fail("query hot blog failed")
	}

	// 用下标 i 遍历，是因为我们要修改 blogs[i] 本身。
	// 如果写成 for _, blog := range blogs，blog 是副本，改了不会影响原切片。
	for i := range blogs {
		user, err := s.userRepo.FindUserByID(ctx, blogs[i].UserID)
		if err != nil {
			// 有些测试数据可能找不到作者；这里先跳过作者补充，不让整个首页失败。
			continue
		}

		// Name/Icon 是 Blog 结构体里的 gorm:"-" 字段，不存在数据库中，只用于返回给前端。
		blogs[i].Name = user.NickName
		blogs[i].Icon = user.Icon
		// 还没实现登录和点赞状态，所以先统一返回 false。
		blogs[i].IsLike = false
	}
	return result.OKWithData(blogs)
}

// QueryBlogLikes 查询点赞用户列表。
func (s *blogService) QueryBlogLikes(ctx context.Context, id int64) result.Result {
	// TODO: Query top liked users from blog:liked:{id}.
	return result.Fail("TODO: query blog likes")
}

// QueryBlogByUserID 查询指定用户发布的博客列表。
func (s *blogService) QueryBlogByUserID(ctx context.Context, userID int64, current int) result.Result {
	if userID <= 0 {
		return result.Fail("invalid userID")
	}
	blogs, err := s.blogRepo.FindBlogsByUserID(ctx, userID, current)
	if err != nil {
		return result.Fail("query userblogs failed")
	}
	if len(blogs) == 0 {
		return result.OKWithData(make([]model.Blog, 0))
	}

	for i := range blogs {
		// 还没实现登录和点赞状态，所以先统一返回 false。
		blogs[i].IsLike = false
	}
	return result.OKWithData(blogs)
}
