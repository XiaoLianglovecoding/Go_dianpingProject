package service

import (
	"context"

	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/repository"

	"github.com/redis/go-redis/v9"
)

type VoucherOrderService interface {
	// SeckillVoucher 抢购秒杀券。
	SeckillVoucher(ctx context.Context, voucherID int64) result.Result
}

type voucherOrderService struct {
	orderRepo   repository.VoucherOrderRepository // orderRepo 负责订单表操作。
	voucherRepo repository.VoucherRepository      // voucherRepo 负责查询秒杀券信息。
	redisClient *redis.Client                     // redisClient 后面用于 Lua 扣库存、消息队列等。
}

// NewVoucherOrderService 创建优惠券订单 Service。
func NewVoucherOrderService(orderRepo repository.VoucherOrderRepository, voucherRepo repository.VoucherRepository, redisClient *redis.Client) VoucherOrderService {
	return &voucherOrderService{orderRepo: orderRepo, voucherRepo: voucherRepo, redisClient: redisClient}
}

// SeckillVoucher 是秒杀下单入口。
//
// 后面会实现：Lua 判断库存和一人一单 -> 返回订单 id -> 异步写订单。
func (s *voucherOrderService) SeckillVoucher(ctx context.Context, voucherID int64) result.Result {
	// TODO: Implement seckill flow with Lua stock check, one-user-one-order, async order creation, and Redis ID worker.
	return result.Fail("TODO: seckill voucher order")
}
