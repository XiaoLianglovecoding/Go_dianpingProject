package service

import (
	"context"

	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/repository"

	"github.com/redis/go-redis/v9"
)

type FollowService interface {
	// IsFollow 查询是否已关注某用户。
	IsFollow(ctx context.Context, followUserID int64) result.Result
	// Follow 关注或取消关注。
	Follow(ctx context.Context, followUserID int64, isFollow bool) result.Result
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
func (s *followService) IsFollow(ctx context.Context, followUserID int64) result.Result {
	// TODO: Check whether current user follows followUserID.
	return result.Fail("TODO: follow or not")
}

// Follow 根据 isFollow 参数决定关注或取消关注。
func (s *followService) Follow(ctx context.Context, followUserID int64, isFollow bool) result.Result {
	// TODO: Create or delete follow relation and sync Redis set follows:{userId}.
	return result.Fail("TODO: update follow")
}

// Common 查询当前用户和另一个用户的共同关注。
func (s *followService) Common(ctx context.Context, otherUserID int64) result.Result {
	// TODO: Find common follows with Redis set intersection or SQL.
	return result.Fail("TODO: common follows")
}
