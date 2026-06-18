package service

import (
	"context"

	"hmdp-go/internal/model"
	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/repository"

	"github.com/redis/go-redis/v9"
)

type ShopService interface {
	QueryByID(ctx context.Context, id int64) result.Result
	SaveShop(ctx context.Context, shop model.Shop) result.Result
	UpdateShop(ctx context.Context, shop model.Shop) result.Result
	QueryByType(ctx context.Context, typeID int64, current int) result.Result
	QueryByName(ctx context.Context, name string, current int) result.Result
}

type shopService struct {
	shopRepo    repository.ShopRepository
	redisClient *redis.Client
}

func NewShopService(shopRepo repository.ShopRepository, redisClient *redis.Client) ShopService {
	return &shopService{shopRepo: shopRepo, redisClient: redisClient}
}

func (s *shopService) QueryByID(ctx context.Context, id int64) result.Result {
	// TODO: Query shop with Redis cache pass-through or logical expiration.
	return result.Fail("TODO: query shop by id")
}

func (s *shopService) SaveShop(ctx context.Context, shop model.Shop) result.Result {
	// TODO: Save shop to MySQL and return generated id.
	return result.Fail("TODO: save shop")
}

func (s *shopService) UpdateShop(ctx context.Context, shop model.Shop) result.Result {
	// TODO: Update shop in MySQL and delete cache:shop:{id}.
	return result.Fail("TODO: update shop")
}

func (s *shopService) QueryByType(ctx context.Context, typeID int64, current int) result.Result {
	// TODO: Query shops by type with pagination; later support GEO query if x/y are provided.
	return result.Fail("TODO: query shop by type")
}

func (s *shopService) QueryByName(ctx context.Context, name string, current int) result.Result {
	// TODO: Query shops by name keyword with pagination.
	return result.Fail("TODO: query shop by name")
}
