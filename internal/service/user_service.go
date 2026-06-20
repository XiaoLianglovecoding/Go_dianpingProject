package service

import (
	"context"

	"errors"
	"hmdp-go/internal/dto"
	"hmdp-go/internal/model"
	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/repository"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserService interface {
	// SendCode 发送登录验证码。
	SendCode(ctx context.Context, phone string) result.Result
	// Login 使用手机号和验证码登录，成功后返回 token。
	Login(ctx context.Context, form dto.LoginFormDTO) result.Result
	// Logout 退出登录。
	Logout(ctx context.Context) result.Result
	// Me 查询当前登录用户。
	Me(ctx context.Context) result.Result
	// QueryUserByID 根据 id 查询用户基础信息。
	QueryUserByID(ctx context.Context, id int64) result.Result
	// QueryUserInfo 查询用户扩展资料。
	QueryUserInfo(ctx context.Context, id int64) result.Result
	// Sign 今日签到。
	Sign(ctx context.Context) result.Result
	// SignCount 查询连续签到天数。
	SignCount(ctx context.Context) result.Result
}

type userService struct {
	userRepo    repository.UserRepository // userRepo 负责用户表相关数据库操作。
	redisClient *redis.Client             // redisClient 后面用于验证码、token、签到 bitmap。
}

// NewUserService 创建用户 Service。
func NewUserService(userRepo repository.UserRepository, redisClient *redis.Client) UserService {
	return &userService{userRepo: userRepo, redisClient: redisClient}
}

// SendCode 负责验证码发送流程。
//
// 后续实现步骤：校验手机号 -> 生成 6 位验证码 -> 写入 Redis -> 发送/打印验证码。
func (s *userService) SendCode(ctx context.Context, phone string) result.Result {
	// TODO: Validate phone, generate code, store login:code:{phone} in Redis, and send/log code.
	return result.Fail("TODO: send user code")
}

// Login 负责登录流程。
//
// 后续实现步骤：校验验证码 -> 查/建用户 -> 生成 token -> 用户信息写入 Redis Hash。
func (s *userService) Login(ctx context.Context, form dto.LoginFormDTO) result.Result {
	// TODO: Validate code, create user if needed, store login:token:{token} hash in Redis, and return token.
	return result.Fail("TODO: user login")
}

// Logout 负责退出登录。
func (s *userService) Logout(ctx context.Context) result.Result {
	// TODO: Delete login token from Redis.
	return result.Fail("TODO: user logout")
}

// Me 返回当前登录用户。
//
// 依赖 AuthMiddleware/RefreshTokenMiddleware 先把用户信息放入 context。
func (s *userService) Me(ctx context.Context) result.Result {
	// TODO: Return current user from request context after auth middleware is implemented.
	return result.Fail("TODO: current user")
}

// QueryUserByID 查询用户基础信息，返回时要转换成 UserDTO，避免泄露密码/手机号。
func (s *userService) QueryUserByID(ctx context.Context, id int64) result.Result {
	if id <= 0 {
		return result.Fail("invalid user id")
	}
	user, err := s.userRepo.FindUserByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return result.Fail("user not found")
	}
	if err != nil {
		return result.Fail("query user failed")
	}
	userDTO := dto.UserDTO{
		ID:       user.ID,
		NickName: user.NickName,
		Icon:     user.Icon,
	}
	return result.OKWithData(userDTO)
}

// QueryUserInfo 查询用户扩展资料。
func (s *userService) QueryUserInfo(ctx context.Context, id int64) result.Result {
	if id <= 0 {
		return result.Fail("invalid user id")
	}
	info, err := s.userRepo.FindUserInfoByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return result.OKWithData(&model.UserInfo{})
	}
	if err != nil {
		return result.Fail("query user info failed")
	}
	infoDTO := &dto.UserInfoDTO{
		City:      info.City,
		Introduce: info.Introduce,
		Fans:      info.Fans,
		Followee:  info.Followee,
		Gender:    info.Gender,
		Credits:   info.Credits,
		Level:     info.Level,
	}
	if info.Birthday != nil {
		infoDTO.Birthday = info.Birthday.Format("2006-01-02")
	}
	return result.OKWithData(infoDTO)
}

// Sign 今日签到。
//
// Java 版使用 Redis Bitmap：sign:{userId}:yyyyMM 的第 day-1 位设置为 1。
func (s *userService) Sign(ctx context.Context) result.Result {
	// TODO: Write today's sign bit with SETBIT sign:{userId}:yyyyMM.
	return result.Fail("TODO: user sign")
}

// SignCount 查询连续签到天数。
func (s *userService) SignCount(ctx context.Context) result.Result {
	// TODO: Count continuous sign days with Redis BITFIELD.
	return result.Fail("TODO: user sign count")
}
