package service

import (
	"context"

	"errors"
	"hmdp-go/internal/model"
	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/repository"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type FollowService interface {
	// IsFollow 查询是否已关注某用户。
	IsFollow(ctx context.Context, targetUserID int64, userID int64) result.Result
	// Follow 关注或取消关注。
	Follow(ctx context.Context, userID int64, isFollow bool, targetUserID int64) result.Result
	// Common 查询共同关注。
	Common(ctx context.Context, otherUserID int64) result.Result
}

type followService struct {
	followRepo  repository.FollowRepository // followRepo 负责 tb_follow 数据库操作。
	redisClient *redis.Client               // redisClient 后面用于关注集合、共同关注交集。
}

// NewFollowService 创建关注 Service。
func NewFollowService(followRepo repository.FollowRepository, redisClient *redis.Client) FollowService {
	return &followService{followRepo: followRepo, redisClient: redisClient}
}

// IsFollow 查询当前用户是否关注了目标用户。
func (s *followService) IsFollow(ctx context.Context, targetUserID int64, currentUserID int64) result.Result {
	// 1. 拦截非法传参
	if targetUserID <= 0 || currentUserID <= 0 {
		return result.Fail("用户ID无效")
	}
	if currentUserID == targetUserID {
		return result.Fail("Donot follow yourself")
	}
	// 2. 调用 Repo 层查询记录 (注意参数顺序：谁 关注了 谁)
	_, err := s.followRepo.FindFollow(ctx, currentUserID, targetUserID)
	// 3. 核心契约解析：把底层错误翻译成前端看得懂的 true / false
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 场景 A：Repo 告诉我们没查到数据，这就代表【未关注】
		return result.OKWithData(false)

	} else if err != nil {
		// 场景 B：发生了真正的数据库异常（比如断网）
		return result.Fail("查询关注状态失败")
	}

	// 4. 场景 C：没有任何错误，说明查到了这条记录，代表【已关注】
	return result.OKWithData(true)
}

// Follow 根据 isFollow 参数决定关注或取消关注。
func (s *followService) Follow(ctx context.Context, currentUserID int64, isFollow bool, targetUserID int64) result.Result {
	// 1. 参数校验
	if targetUserID <= 0 || currentUserID <= 0 {
		return result.Fail("用户ID无效")
	}
	if currentUserID == targetUserID {
		return result.Fail("Donot follow yourself")
	}
	// 2. 根据 isFollow 决定动作
	if isFollow {
		// 【关注逻辑】：实例化一个 Follow 结构体对象，传给 Repo
		// 注意：如果您的 model 里面有类似 TimeFields 的结构，CreateTime 会由 GORM 自动生成
		follow := &model.Follow{
			UserID:       currentUserID,
			FollowUserID: targetUserID,
		}
		err := s.followRepo.SaveFollow(ctx, follow)
		if err != nil {
			return result.Fail("关注失败")
		}
	} else {
		// 【取消关注逻辑】：直接把两个 ID 传给 DeleteFollow 方法
		err := s.followRepo.DeleteFollow(ctx, currentUserID, targetUserID)
		if err != nil {
			return result.Fail("取消关注失败")
		}
	}
	return result.OK()
}

// Common 查询当前用户和另一个用户的共同关注。
func (s *followService) Common(ctx context.Context, otherUserID int64) result.Result {
	// TODO: Find common follows with Redis set intersection or SQL.
	return result.Fail("TODO: common follows")
}
