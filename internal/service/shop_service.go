package service

import (
	"context"

	"hmdp-go/internal/model"
	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/repository"

	"github.com/redis/go-redis/v9"
)

type ShopService interface {
	// QueryByID 查询店铺详情。
	QueryByID(ctx context.Context, id int64) result.Result
	// SaveShop 新增店铺。
	SaveShop(ctx context.Context, shop model.Shop) result.Result
	// UpdateShop 更新店铺。
	UpdateShop(ctx context.Context, shop model.Shop) result.Result
	// QueryByType 根据分类查询店铺列表。
	QueryByType(ctx context.Context, typeID int64, current int) result.Result
	// QueryByName 根据关键词搜索店铺。
	QueryByName(ctx context.Context, name string, current int) result.Result
}

type shopService struct {
	shopRepo    repository.ShopRepository // shopRepo 负责 tb_shop 数据库操作。
	redisClient *redis.Client             // redisClient 后面用于店铺缓存和 GEO 查询。
}

// NewShopService 创建店铺 Service。
func NewShopService(shopRepo repository.ShopRepository, redisClient *redis.Client) ShopService {
	return &shopService{shopRepo: shopRepo, redisClient: redisClient}
}

// QueryByID 查询店铺详情。
//
// 后面会重点学习 Redis 缓存穿透、缓存击穿、逻辑过期等内容。
func (s *shopService) QueryByID(ctx context.Context, id int64) result.Result {
	// TODO: Query shop with Redis cache pass-through or logical expiration.
	return result.Fail("TODO: query shop by id")
}

// SaveShop 新增店铺，保存成功后返回店铺 id。
func (s *shopService) SaveShop(ctx context.Context, shop model.Shop) result.Result {
	// TODO: Save shop to MySQL and return generated id.
	return result.Fail("TODO: save shop")
}

// UpdateShop 更新店铺。
//
// 更新 MySQL 后要删除 Redis 缓存，避免前端看到旧数据。
func (s *shopService) UpdateShop(ctx context.Context, shop model.Shop) result.Result {
	// TODO: Update shop in MySQL and delete cache:shop:{id}.
	return result.Fail("TODO: update shop")
}

// QueryByType 根据店铺分类分页查询。
func (s *shopService) QueryByType(ctx context.Context, typeID int64, current int) result.Result {
	// TODO: Query shops by type with pagination; later support GEO query if x/y are provided.
	return result.Fail("TODO: query shop by type")
}

// QueryByName 根据店铺名称搜索。
func (s *shopService) QueryByName(ctx context.Context, name string, current int) result.Result {
	// TODO: Query shops by name keyword with pagination.
	return result.Fail("TODO: query shop by name")
}
