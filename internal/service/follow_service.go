package service

import (
	"context"
	"strconv"

	"errors"
	"hmdp-go/internal/dto"
	"hmdp-go/internal/model"
	"hmdp-go/internal/pkg/constants"
	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/repository"
	"log"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type FollowService interface {
	// IsFollow 查询是否已关注某用户。
	IsFollow(ctx context.Context, targetUserID int64, userID int64) result.Result
	// Follow 关注或取消关注。
	Follow(ctx context.Context, userID int64, targetUserID int64, isFollow bool) result.Result
	// Common 查询共同关注。
	Common(ctx context.Context, currentUserID int64, otherUserID int64) result.Result
}

type followService struct {
	followRepo  repository.FollowRepository // followRepo 负责 tb_follow 数据库操作。
	userRepo    repository.UserRepository   // userRepo 负责 tb_user 数据库操作。
	redisClient *redis.Client               // redisClient 后面用于关注集合、共同关注交集。
}

// NewFollowService 创建关注 Service。
func NewFollowService(followRepo repository.FollowRepository, userRepo repository.UserRepository, redisClient *redis.Client) FollowService {
	return &followService{followRepo: followRepo, userRepo: userRepo, redisClient: redisClient}
}

// IsFollow 查询当前用户是否关注了目标用户。
func (s *followService) IsFollow(ctx context.Context, targetUserID int64, currentUserID int64) result.Result {
	// 1. 拦截非法传参
	if targetUserID <= 0 || currentUserID <= 0 {
		return result.Fail("用户ID无效")
	}
	if currentUserID == targetUserID {
		return result.OKWithData(false)
	}
	// 拼接 Redis Key
	key := constants.FollowsKey + strconv.FormatInt(currentUserID, 10)

	// 只有这个用户的关注集合存在时，才信 Redis 的 false。
	exists, err := s.redisClient.Exists(ctx, key).Result()
	if err == nil && exists > 0 {
		isMember, err := s.redisClient.SIsMember(ctx, key, strconv.FormatInt(targetUserID, 10)).Result()
		if err == nil {
			return result.OKWithData(isMember)
		}
		log.Printf("[FollowService.IsFollow] Redis SIsMember failed: %v", err)
	}

	// redis中不存在，再调用 Repo 层查询记录 (注意参数顺序：谁 关注了 谁)
	_, err = s.followRepo.FindFollow(ctx, currentUserID, targetUserID)
	// 核心契约解析：把底层错误翻译成前端看得懂的 true / false
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 场景 A：Repo 告诉我们没查到数据，这就代表【未关注】
		return result.OKWithData(false)
	}
	if err != nil {
		// 场景 B：发生了真正的数据库异常（比如断网）
		return result.Fail("查询关注状态失败")
	}

	// 4. 场景 C：没有任何错误，说明查到了这条记录，代表【已关注】
	return result.OKWithData(true)
}

// Follow 根据 isFollow 参数决定关注或取消关注。
func (s *followService) Follow(ctx context.Context, currentUserID int64, targetUserID int64, isFollow bool) result.Result {
	// 1. 参数校验
	if targetUserID <= 0 || currentUserID <= 0 {
		return result.Fail("用户ID无效")
	}
	if currentUserID == targetUserID {
		return result.Fail("Donot follow yourself")
	}
	key := constants.FollowsKey + strconv.FormatInt(currentUserID, 10)
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

		// MySQL 成功后，写入 Redis Set
		// 这里把 targetUserID 传入，Redis 会自动转成字符串存储
		err = s.redisClient.SAdd(ctx, key, targetUserID).Err()
		if err != nil {
			// 直接返回失败便于排查，生产环境通常记录日志或异步补偿
			// 记录错误日志，方便排查告警
			log.Printf("Redis 写入失败, 请重试. userID: %d", currentUserID)
		}
	} else {
		// 【取消关注逻辑】：直接把两个 ID 传给 DeleteFollow 方法
		err := s.followRepo.DeleteFollow(ctx, currentUserID, targetUserID)
		if err != nil {
			return result.Fail("取消关注失败")
		}
		// MySQL 成功后，移除Redis Set
		// 这里把 targetUserID 传入，Redis 会自动转成字符串存储
		err = s.redisClient.SRem(ctx, key, targetUserID).Err()
		if err != nil {
			return result.Fail("同步取消关注缓存失败")
		}
	}
	return result.OK()
}

// Common 查询当前用户和另一个用户的共同关注。
func (s *followService) Common(ctx context.Context, currentUserID int64, otherUserID int64) result.Result {
	// 1. 参数校验
	if otherUserID <= 0 || currentUserID <= 0 {
		return result.Fail("用户ID无效")
	}
	if currentUserID == otherUserID {
		return result.OKWithData([]dto.UserDTO{})
	}
	// 2. 查 Redis 求交集 SInter
	currentKey := constants.FollowsKey + strconv.FormatInt(currentUserID, 10)
	otherKey := constants.FollowsKey + strconv.FormatInt(otherUserID, 10)
	commonFollowsStr, err := s.redisClient.SInter(ctx, currentKey, otherKey).Result()
	if err != nil {
		log.Printf("[FollowService.Common] Redis SInter failed: %v", err)
		// Redis 挂了，为了高可用，必须降级去查 MySQL
		return s.fallbackToMySQL(ctx, currentUserID, otherUserID)
	}

	// 3. 核心逻辑判断：如果 Redis 返回交集为空，需要区分原因
	if len(commonFollowsStr) == 0 {
		// 检查这两个用户的 Key 是否都在 Redis 中存在
		// Exists 会返回存在的 Key 的数量 (0, 1, 或 2)
		existCount, _ := s.redisClient.Exists(ctx, currentKey, otherKey).Result()

		if existCount == 2 {
			// 【场景 A】：两个人的关注列表都在缓存里，但交集为空。
			// 结论：他们【真的】没有共同关注！直接阻断，保护数据库！
			return result.OKWithData([]dto.UserDTO{})
		} else {
			// 【场景 B】：至少有一个人的关注列表不在缓存里（缓存丢失/过期）。
			// 结论：Redis 的数据是不完整的，必须降级查 MySQL！
			// (并且理论上查完 MySQL 后，应该顺手把数据重新回写进 Redis 补全缓存，这里为了简洁先省略)
			return s.fallbackToMySQL(ctx, currentUserID, otherUserID)
		}
	}

	// 4. Redis 成功命中，且有共同关注，转换类型
	commonIDs := make([]int64, 0, len(commonFollowsStr))
	for _, idStr := range commonFollowsStr {
		if id, parseErr := strconv.ParseInt(idStr, 10, 64); parseErr == nil {
			commonIDs = append(commonIDs, id)
		}
	}
	// 5. 根据 ID 批量查用户信息 (这里可以直接复用之前的批量查方法)
	users, err := s.userRepo.FindUsersByIDs(ctx, commonIDs)
	if err != nil {
		log.Printf("[FollowService.Common] 数据库批量查询用户失败: %v", err)
		return result.Fail("获取共同关注列表失败")
	}

	return result.OKWithData(toUserDTOs(users))
}

// 抽取出来的降级查 MySQL 的私有方法，保持主逻辑干净
func (s *followService) fallbackToMySQL(ctx context.Context, currentUserID, otherUserID int64) result.Result {
	log.Printf("触发降级：从 MySQL 查询用户 %d 和 %d 的共同关注", currentUserID, otherUserID)

	users, err := s.followRepo.FindCommonFollows(ctx, currentUserID, otherUserID)
	if err != nil {
		log.Printf("兜底查库失败: %v", err)
		return result.Fail("获取共同关注列表失败")
	}

	userDTOs := make([]dto.UserDTO, 0, len(users))
	for _, u := range users {
		userDTOs = append(userDTOs, dto.UserDTO{ID: u.ID, NickName: u.NickName, Icon: u.Icon})
	}

	return result.OKWithData(userDTOs)
}

func toUserDTOs(users []model.User) []dto.UserDTO {
	userDTOs := make([]dto.UserDTO, 0, len(users))
	for _, u := range users {
		userDTOs = append(userDTOs, dto.UserDTO{
			ID:       u.ID,
			NickName: u.NickName,
			Icon:     u.Icon,
		})
	}
	return userDTOs
}
