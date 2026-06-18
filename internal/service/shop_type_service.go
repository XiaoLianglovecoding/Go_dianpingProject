package service

import (
	"context"

	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/repository"

	"github.com/redis/go-redis/v9"
)

type ShopTypeService interface {
	List(ctx context.Context) result.Result
}

type shopTypeService struct {
	shopTypeRepo repository.ShopTypeRepository
	redisClient  *redis.Client
}

func NewShopTypeService(shopTypeRepo repository.ShopTypeRepository, redisClient *redis.Client) ShopTypeService {
	return &shopTypeService{shopTypeRepo: shopTypeRepo, redisClient: redisClient}
}

func (s *shopTypeService) List(ctx context.Context) result.Result {
	// TODO: Query shop type list from Redis cache first, then MySQL ordered by sort.
	return result.Fail("TODO: shop type list")
}
