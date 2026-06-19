package service

import (
	"context"

	"hmdp-go/internal/model"
	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/repository"

	"github.com/redis/go-redis/v9"
)

type VoucherService interface {
	// AddVoucher 新增普通优惠券。
	AddVoucher(ctx context.Context, voucher model.Voucher) result.Result
	// AddSeckillVoucher 新增秒杀券。
	AddSeckillVoucher(ctx context.Context, voucher model.Voucher) result.Result
	// QueryVoucherOfShop 查询某店铺的优惠券列表。
	QueryVoucherOfShop(ctx context.Context, shopID int64) result.Result
}

type voucherService struct {
	voucherRepo repository.VoucherRepository // voucherRepo 负责优惠券表操作。
	redisClient *redis.Client                // redisClient 后面用于秒杀库存。
}

// NewVoucherService 创建优惠券 Service。
func NewVoucherService(voucherRepo repository.VoucherRepository, redisClient *redis.Client) VoucherService {
	return &voucherService{voucherRepo: voucherRepo, redisClient: redisClient}
}

// AddVoucher 保存普通优惠券。
func (s *voucherService) AddVoucher(ctx context.Context, voucher model.Voucher) result.Result {
	// TODO: Save normal voucher to tb_voucher.
	return result.Fail("TODO: add voucher")
}

// AddSeckillVoucher 保存秒杀券，并把库存预热到 Redis。
func (s *voucherService) AddSeckillVoucher(ctx context.Context, voucher model.Voucher) result.Result {
	// TODO: Save voucher and seckill stock/time, then preload seckill:stock:{id} in Redis.
	return result.Fail("TODO: add seckill voucher")
}

// QueryVoucherOfShop 查询店铺详情页展示的优惠券。
func (s *voucherService) QueryVoucherOfShop(ctx context.Context, shopID int64) result.Result {
	if shopID <= 0 {
		return result.Fail("invalid shop id")
	}
	vouchers, err := s.voucherRepo.FindVouchersByShopID(ctx, shopID)
	if err != nil {
		return result.Fail("query voucher of shop failed")
	}
	for i := range vouchers {
		if vouchers[i].SeckillVoucher != nil {
			vouchers[i].Stock = vouchers[i].SeckillVoucher.Stock
			vouchers[i].BeginTime = &vouchers[i].SeckillVoucher.BeginTime
			vouchers[i].EndTime = &vouchers[i].SeckillVoucher.EndTime
		}
	}
	return result.OKWithData(vouchers)
}
