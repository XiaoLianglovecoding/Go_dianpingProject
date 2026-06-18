package service

import (
	"context"

	"hmdp-go/internal/model"
	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/repository"

	"github.com/redis/go-redis/v9"
)

type VoucherService interface {
	AddVoucher(ctx context.Context, voucher model.Voucher) result.Result
	AddSeckillVoucher(ctx context.Context, voucher model.Voucher) result.Result
	QueryVoucherOfShop(ctx context.Context, shopID int64) result.Result
}

type voucherService struct {
	voucherRepo repository.VoucherRepository
	redisClient *redis.Client
}

func NewVoucherService(voucherRepo repository.VoucherRepository, redisClient *redis.Client) VoucherService {
	return &voucherService{voucherRepo: voucherRepo, redisClient: redisClient}
}

func (s *voucherService) AddVoucher(ctx context.Context, voucher model.Voucher) result.Result {
	// TODO: Save normal voucher to tb_voucher.
	return result.Fail("TODO: add voucher")
}

func (s *voucherService) AddSeckillVoucher(ctx context.Context, voucher model.Voucher) result.Result {
	// TODO: Save voucher and seckill stock/time, then preload seckill:stock:{id} in Redis.
	return result.Fail("TODO: add seckill voucher")
}

func (s *voucherService) QueryVoucherOfShop(ctx context.Context, shopID int64) result.Result {
	// TODO: Query vouchers for shop detail page.
	return result.Fail("TODO: query voucher of shop")
}
