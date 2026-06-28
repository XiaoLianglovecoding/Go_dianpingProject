package service

import (
	"context"
	"log"
	"strconv"
	"time"

	"hmdp-go/internal/model"
	"hmdp-go/internal/pkg/constants"
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
	if voucher.ShopID <= 0 {
		return result.Fail("invalid shop id")
	}
	if voucher.Stock <= 0 {
		return result.Fail("stock must be greater than 0")
	}
	if voucher.BeginTime == nil || voucher.EndTime == nil {
		return result.Fail("beginTime and endTime can not be empty")
	}
	if !voucher.EndTime.After(*voucher.BeginTime) {
		return result.Fail("endTime must be after beginTime")
	}
	// 1. 标记为秒杀券
	voucher.Type = 1

	// 2. 构建秒杀券关联模型
	seckillVoucher := &model.SeckillVoucher{
		Stock:     voucher.Stock,
		BeginTime: *voucher.BeginTime,
		EndTime:   *voucher.EndTime,
	}

	// 3. 【核心优化】：调用 Repository 的事务方法，同时保存两张表
	if err := s.voucherRepo.SaveSeckillVoucherTx(ctx, &voucher, seckillVoucher); err != nil {
		log.Printf("[AddSeckillVoucher] DB transaction failed: %v", err)
		return result.Fail("数据库保存失败，请重试")
	}

	// 4. 预热 Redis 库存。
	stockKey := constants.SeckillStockKey + strconv.FormatInt(voucher.ID, 10)
	orderKey := constants.SeckillOrdersKey + strconv.FormatInt(voucher.ID, 10)

	// 可以让库存 key 在秒杀结束后自动过期，避免 Redis 里长期堆旧数据。
	ttl := time.Until(*voucher.EndTime)
	if ttl <= 0 {
		return result.Fail("seckill voucher has expired")
	}

	pipe := s.redisClient.TxPipeline()
	pipe.Set(ctx, stockKey, voucher.Stock, ttl)
	pipe.Del(ctx, orderKey)

	// 5. 【核心优化】：捕获 Redis 异常，给出特定降级提示
	if _, err := pipe.Exec(ctx); err != nil {
		log.Printf("[AddSeckillVoucher] redis preheat failed for voucherID=%d: %v", voucher.ID, err)
		// 此时数据库已经有数据了，向前端返回特殊的 Fail 提示，引导管理员重试
		return result.Fail("保存成功，但Redis库存预热失败，需要重试预热")
	}

	return result.OKWithData(voucher.ID)
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
