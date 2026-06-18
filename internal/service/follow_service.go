package service

import (
	"context"

	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/repository"

	"github.com/redis/go-redis/v9"
)

type FollowService interface {
	IsFollow(ctx context.Context, followUserID int64) result.Result
	Follow(ctx context.Context, followUserID int64, isFollow bool) result.Result
	Common(ctx context.Context, otherUserID int64) result.Result
}

type followService struct {
	followRepo  repository.FollowRepository
	redisClient *redis.Client
}

func NewFollowService(followRepo repository.FollowRepository, redisClient *redis.Client) FollowService {
	return &followService{followRepo: followRepo, redisClient: redisClient}
}

func (s *followService) IsFollow(ctx context.Context, followUserID int64) result.Result {
	// TODO: Check whether current user follows followUserID.
	return result.Fail("TODO: follow or not")
}

func (s *followService) Follow(ctx context.Context, followUserID int64, isFollow bool) result.Result {
	// TODO: Create or delete follow relation and sync Redis set follows:{userId}.
	return result.Fail("TODO: update follow")
}

func (s *followService) Common(ctx context.Context, otherUserID int64) result.Result {
	// TODO: Find common follows with Redis set intersection or SQL.
	return result.Fail("TODO: common follows")
}
