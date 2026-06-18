package service

import (
	"context"

	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/repository"

	"github.com/redis/go-redis/v9"
)

// ShopTypeService 定义“店铺分类”业务能力。
//
// Handler 只知道调用这个接口，不关心内部是查 MySQL 还是 Redis。
type ShopTypeService interface {
	// List 返回首页顶部店铺分类列表。
	List(ctx context.Context) result.Result
}

type shopTypeService struct {
	// shopTypeRepo 负责查询 tb_shop_type 表。
	shopTypeRepo repository.ShopTypeRepository
	// redisClient 后面可以用来缓存分类列表，当前先直接查 MySQL。
	redisClient *redis.Client
}

// NewShopTypeService 创建店铺分类 Service。
func NewShopTypeService(shopTypeRepo repository.ShopTypeRepository, redisClient *redis.Client) ShopTypeService {
	return &shopTypeService{shopTypeRepo: shopTypeRepo, redisClient: redisClient}
}

// List 是 /shop-type/list 的业务逻辑。
//
// 当前版本：直接从 MySQL 查询分类，再包装成统一 Result 返回。
// 后续优化：先从 Redis 读 cache:shopType，未命中再查 MySQL。
func (s *shopTypeService) List(ctx context.Context) result.Result {
	shopTypes, err := s.shopTypeRepo.FindShopTypes(ctx)
	if err != nil {
		return result.Fail("query shop type failed")
	}
	return result.OKWithData(shopTypes)
}
