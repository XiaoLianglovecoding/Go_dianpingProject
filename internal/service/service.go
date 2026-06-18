package service

import (
	"hmdp-go/internal/repository"

	"github.com/redis/go-redis/v9"
)

// Services 集中保存所有 Service。
//
// Service 层负责“业务逻辑”，比如登录校验、缓存策略、拼装返回数据。
// 它会调用 Repository 查数据库，也会调用 Redis 做缓存、点赞、签到等。
type Services struct {
	User         UserService
	Shop         ShopService
	ShopType     ShopTypeService
	Blog         BlogService
	Follow       FollowService
	Upload       UploadService
	Voucher      VoucherService
	VoucherOrder VoucherOrderService
}

// NewServices 创建所有 Service，并把它们需要的依赖注入进去。
//
// repos: 数据库访问层集合。
// redisClient: Redis 客户端，后续做缓存、登录 token、点赞、签到时会用。
func NewServices(repos *repository.Repositories, redisClient *redis.Client) *Services {
	return &Services{
		User:         NewUserService(repos.User, redisClient),
		Shop:         NewShopService(repos.Shop, redisClient),
		ShopType:     NewShopTypeService(repos.ShopType, redisClient),
		Blog:         NewBlogService(repos.Blog, repos.User, redisClient),
		Follow:       NewFollowService(repos.Follow, redisClient),
		Upload:       NewUploadService(),
		Voucher:      NewVoucherService(repos.Voucher, redisClient),
		VoucherOrder: NewVoucherOrderService(repos.VoucherOrder, repos.Voucher, redisClient),
	}
}
