package repository

import (
	"context"

	"hmdp-go/internal/model"

	"gorm.io/gorm"
)

type VoucherRepository interface {
	// SaveVoucher 保存普通优惠券。
	SaveVoucher(ctx context.Context, voucher *model.Voucher) error
	// SaveSeckillVoucher 保存秒杀券额外信息。
	SaveSeckillVoucher(ctx context.Context, voucher *model.SeckillVoucher) error
	// FindVouchersByShopID 查询某个店铺下的优惠券。
	FindVouchersByShopID(ctx context.Context, shopID int64) ([]model.Voucher, error)
	// FindSeckillVoucherByID 查询秒杀券详情。
	FindSeckillVoucherByID(ctx context.Context, voucherID int64) (*model.SeckillVoucher, error)
}

type voucherRepository struct {
	// db 是 GORM 数据库连接对象。
	db *gorm.DB
}

// NewVoucherRepository 创建优惠券 Repository。
func NewVoucherRepository(db *gorm.DB) VoucherRepository {
	return &voucherRepository{db: db}
}

// SaveVoucher 后面用于插入 tb_voucher。
func (r *voucherRepository) SaveVoucher(ctx context.Context, voucher *model.Voucher) error {
	// TODO: Insert voucher into tb_voucher.
	return nil
}

// SaveSeckillVoucher 后面用于插入 tb_seckill_voucher。
func (r *voucherRepository) SaveSeckillVoucher(ctx context.Context, voucher *model.SeckillVoucher) error {
	// TODO: Insert seckill metadata into tb_seckill_voucher.
	return nil
}

// FindVouchersByShopID 后面用于店铺详情页展示优惠券列表。
func (r *voucherRepository) FindVouchersByShopID(ctx context.Context, shopID int64) ([]model.Voucher, error) {
	var vouchers []model.Voucher
	err := r.db.WithContext(ctx).
		Preload("SeckillVoucher").
		Where("shop_id = ?", shopID).
		Find(&vouchers).Error
	if err != nil {
		return nil, err
	}
	return vouchers, nil
}

// FindSeckillVoucherByID 后面用于秒杀下单前检查库存和时间。
func (r *voucherRepository) FindSeckillVoucherByID(ctx context.Context, voucherID int64) (*model.SeckillVoucher, error) {
	// TODO: Query tb_seckill_voucher by voucher_id.
	return nil, nil
}
