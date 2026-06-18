package service

import (
	"context"

	"hmdp-go/internal/dto"
	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/repository"

	"github.com/redis/go-redis/v9"
)

type UserService interface {
	SendCode(ctx context.Context, phone string) result.Result
	Login(ctx context.Context, form dto.LoginFormDTO) result.Result
	Logout(ctx context.Context) result.Result
	Me(ctx context.Context) result.Result
	QueryUserByID(ctx context.Context, id int64) result.Result
	QueryUserInfo(ctx context.Context, id int64) result.Result
	Sign(ctx context.Context) result.Result
	SignCount(ctx context.Context) result.Result
}

type userService struct {
	userRepo    repository.UserRepository
	redisClient *redis.Client
}

func NewUserService(userRepo repository.UserRepository, redisClient *redis.Client) UserService {
	return &userService{userRepo: userRepo, redisClient: redisClient}
}

func (s *userService) SendCode(ctx context.Context, phone string) result.Result {
	// TODO: Validate phone, generate code, store login:code:{phone} in Redis, and send/log code.
	return result.Fail("TODO: send user code")
}

func (s *userService) Login(ctx context.Context, form dto.LoginFormDTO) result.Result {
	// TODO: Validate code, create user if needed, store login:token:{token} hash in Redis, and return token.
	return result.Fail("TODO: user login")
}

func (s *userService) Logout(ctx context.Context) result.Result {
	// TODO: Delete login token from Redis.
	return result.Fail("TODO: user logout")
}

func (s *userService) Me(ctx context.Context) result.Result {
	// TODO: Return current user from request context after auth middleware is implemented.
	return result.Fail("TODO: current user")
}

func (s *userService) QueryUserByID(ctx context.Context, id int64) result.Result {
	// TODO: Query user by id and convert to UserDTO.
	return result.Fail("TODO: query user by id")
}

func (s *userService) QueryUserInfo(ctx context.Context, id int64) result.Result {
	// TODO: Query user info by user id.
	return result.Fail("TODO: query user info")
}

func (s *userService) Sign(ctx context.Context) result.Result {
	// TODO: Write today's sign bit with SETBIT sign:{userId}:yyyyMM.
	return result.Fail("TODO: user sign")
}

func (s *userService) SignCount(ctx context.Context) result.Result {
	// TODO: Count continuous sign days with Redis BITFIELD.
	return result.Fail("TODO: user sign count")
}
