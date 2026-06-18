package service

import (
	"context"

	"hmdp-go/internal/model"
	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/repository"

	"github.com/redis/go-redis/v9"
)

type BlogService interface {
	SaveBlog(ctx context.Context, blog model.Blog) result.Result
	QueryByID(ctx context.Context, id int64) result.Result
	LikeBlog(ctx context.Context, id int64) result.Result
	QueryMyBlog(ctx context.Context, current int) result.Result
	QueryHotBlog(ctx context.Context, current int) result.Result
	QueryBlogLikes(ctx context.Context, id int64) result.Result
	QueryBlogByUserID(ctx context.Context, userID int64, current int) result.Result
}

type blogService struct {
	blogRepo    repository.BlogRepository
	userRepo    repository.UserRepository
	redisClient *redis.Client
}

func NewBlogService(blogRepo repository.BlogRepository, userRepo repository.UserRepository, redisClient *redis.Client) BlogService {
	return &blogService{blogRepo: blogRepo, userRepo: userRepo, redisClient: redisClient}
}

func (s *blogService) SaveBlog(ctx context.Context, blog model.Blog) result.Result {
	// TODO: Get current user id from context, save blog, and push feed to followers.
	return result.Fail("TODO: save blog")
}

func (s *blogService) QueryByID(ctx context.Context, id int64) result.Result {
	// TODO: Query blog by id, attach author info, and mark whether current user liked it.
	return result.Fail("TODO: query blog by id")
}

func (s *blogService) LikeBlog(ctx context.Context, id int64) result.Result {
	// TODO: Toggle like status with blog:liked:{id} sorted set and update liked count.
	return result.Fail("TODO: like blog")
}

func (s *blogService) QueryMyBlog(ctx context.Context, current int) result.Result {
	// TODO: Query current user's blogs.
	return result.Fail("TODO: query my blog")
}

func (s *blogService) QueryHotBlog(ctx context.Context, current int) result.Result {
	// TODO: Query hot blogs ordered by liked count.
	return result.Fail("TODO: query hot blog")
}

func (s *blogService) QueryBlogLikes(ctx context.Context, id int64) result.Result {
	// TODO: Query top liked users from blog:liked:{id}.
	return result.Fail("TODO: query blog likes")
}

func (s *blogService) QueryBlogByUserID(ctx context.Context, userID int64, current int) result.Result {
	// TODO: Query blogs by target user id.
	return result.Fail("TODO: query blog by user")
}
