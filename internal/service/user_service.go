package service

import (
	"context"
	"errors"
	"fmt"
	"hmdp-go/internal/dto"
	"hmdp-go/internal/model"
	"hmdp-go/internal/pkg/constants"
	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/repository"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid" // 用于生成 Token
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserService interface {
	// SendCode 发送登录验证码。
	SendCode(ctx context.Context, phone string) result.Result
	// Login 使用手机号和验证码登录，成功后返回 token。
	Login(ctx context.Context, form dto.LoginFormDTO) result.Result
	// Logout 退出登录。
	Logout(ctx context.Context, token string) result.Result
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
	if len(phone) == 0 {
		return result.Fail("phone is empty")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06d", r.Int31n(1000000))

	key := constants.LoginCodeKey + phone

	ttl := time.Duration(constants.LoginCodeTTLMinutes) * time.Minute

	err := s.redisClient.Set(ctx, key, code, ttl).Err()
	if err != nil {
		return result.Fail("have code failed,please try it later")
	}

	fmt.Printf("验证码:%s\n", code)
	return result.OK()
}

// Login 负责登录流程。
//
// 后续实现步骤：校验验证码 -> 查/建用户 -> 生成 token -> 用户信息写入 Redis Hash。
func (s *userService) Login(ctx context.Context, form dto.LoginFormDTO) result.Result {
	phone := form.Phone
	code := form.Code

	// 1. 校验手机号和验证码是否为空
	if len(phone) == 0 || len(code) == 0 {
		return result.Fail("phone or code is empty")
	}

	// 2. 从 Redis 取出验证码 (对应 Java 的 opsForValue().get)
	key := constants.LoginCodeKey + phone
	rediscode, err := s.redisClient.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return result.Fail("验证码已过期或不存在")
	} else if err != nil {
		return result.Fail("系统异常")
	}

	// 3. 判断验证码是否一致
	if code != rediscode {
		return result.Fail("验证码错误")
	}

	// 4. 根据 phone 查询 tb_user
	user, err := s.userRepo.FindUserByPhone(ctx, phone)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 5. 如果用户不存在，创建新用户
		// 这里的 model.User 结构需要根据我们的实际字段调整
		user = &model.User{
			Phone:    phone,
			NickName: "user_" + uuid.New().String()[:8],
		}
		createErr := s.userRepo.CreateUser(ctx, user)
		if createErr != nil {
			return result.Fail("创建新用户失败")
		}
	} else if err != nil {
		return result.Fail("查询用户异常")
	}
	// 6. 生成 Token
	token := uuid.New().String()

	// 7. 把 UserDTO 准备好并保存到 Redis Hash
	tokenKey := constants.LoginUserKey + token

	// 在 Go 中存入 Redis Hash，最好将数据转成 map[string]interface{}
	// 并且注意 ID 等整型最好转为字符串，避免序列化兼容问题
	// 存入 Hash (对应 Java 的 opsForHash().putAll)
	// Redis 3.x 不支持使用 HSET 一次性存多个键值对，必须用 HMSET！
	userMap := map[string]interface{}{
		"id":       strconv.FormatInt(user.ID, 10),
		"nickName": user.NickName,
		"icon":     user.Icon,
	}

	// 存入 Hash：将 HSet 改为 HMSet (Hash Multiple Set)
	err = s.redisClient.HMSet(ctx, tokenKey, userMap).Err()
	if err != nil {
		log.Printf("❌ Redis HSet 保存登录状态失败: %v", err)
		return result.Fail("保存登陆状态失败")
	}

	// 8. 设置 Token TTL 为 30 分钟
	ttl := time.Duration(constants.LoginUserTTLMinutes) * time.Minute
	if err := s.redisClient.Expire(ctx, tokenKey, ttl).Err(); err != nil {
		// 强依赖：如果设置过期时间失败，必须记录日志并返回错误，防止 Token 永久驻留
		log.Printf("failed to set expire for token %s: %v", tokenKey, err) // 建议加上日志
		return result.Fail("系统内部错误，登录状态保存失败")                              // 假设你的 result 包有类似 Fail/Error 的方法
	}

	//9.登录成功后删除验证码
	if err := s.redisClient.Del(ctx, key).Err(); err != nil {
		// 弱依赖：清理操作失败通常不影响用户此次登录，记录 Warn/Error 日志即可，不阻塞返回
		log.Printf("failed to delete captcha key %s: %v", key, err)
	}
	return result.OKWithData(token)
}

// Logout 负责退出登录。
func (s *userService) Logout(ctx context.Context, token string) result.Result {
	if token == "" {
		return result.OK()
	}
	tokenkey := constants.LoginUserKey + token
	deletecount, err := s.redisClient.Del(ctx, tokenkey).Result()
	if err != nil {
		log.Printf(" Redis 删除 Token 系统异常 %s: %v", tokenkey, err)
		return result.Fail("删除token失败,Redis 报错")
	}

	// 5. 验证是否是重复退出（可选：纯粹为了日志观测，不影响返回给前端的结果）
	if deletecount == 0 {
		// key 本来就不存在（可能已经过期自然死亡，或者用户手抖狂点退出按钮）
		// 按照我们的幂等性原则，这也算作退出成功！
		log.Printf(" Token 已失效或不存在，无需重复删除: %s", tokenkey)
	}

	return result.OK()
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
