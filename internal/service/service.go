package service

import (
	"hmdp-go/internal/repository"

	"github.com/redis/go-redis/v9"
)

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
