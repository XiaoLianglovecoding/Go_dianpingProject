package service

import (
	"context"

	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/repository"

	"github.com/redis/go-redis/v9"
)

type VoucherOrderService interface {
	SeckillVoucher(ctx context.Context, voucherID int64) result.Result
}

type voucherOrderService struct {
	orderRepo   repository.VoucherOrderRepository
	voucherRepo repository.VoucherRepository
	redisClient *redis.Client
}

func NewVoucherOrderService(orderRepo repository.VoucherOrderRepository, voucherRepo repository.VoucherRepository, redisClient *redis.Client) VoucherOrderService {
	return &voucherOrderService{orderRepo: orderRepo, voucherRepo: voucherRepo, redisClient: redisClient}
}

func (s *voucherOrderService) SeckillVoucher(ctx context.Context, voucherID int64) result.Result {
	// TODO: Implement seckill flow with Lua stock check, one-user-one-order, async order creation, and Redis ID worker.
	return result.Fail("TODO: seckill voucher order")
}
