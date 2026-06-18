package repository

import (
	"context"

	"hmdp-go/internal/model"

	"gorm.io/gorm"
)

type VoucherRepository interface {
	SaveVoucher(ctx context.Context, voucher *model.Voucher) error
	SaveSeckillVoucher(ctx context.Context, voucher *model.SeckillVoucher) error
	FindVouchersByShopID(ctx context.Context, shopID int64) ([]model.Voucher, error)
	FindSeckillVoucherByID(ctx context.Context, voucherID int64) (*model.SeckillVoucher, error)
}

type voucherRepository struct {
	db *gorm.DB
}

func NewVoucherRepository(db *gorm.DB) VoucherRepository {
	return &voucherRepository{db: db}
}

func (r *voucherRepository) SaveVoucher(ctx context.Context, voucher *model.Voucher) error {
	// TODO: Insert voucher into tb_voucher.
	return nil
}

func (r *voucherRepository) SaveSeckillVoucher(ctx context.Context, voucher *model.SeckillVoucher) error {
	// TODO: Insert seckill metadata into tb_seckill_voucher.
	return nil
}

func (r *voucherRepository) FindVouchersByShopID(ctx context.Context, shopID int64) ([]model.Voucher, error) {
	// TODO: Query vouchers by shop_id and join seckill fields when needed.
	return nil, nil
}

func (r *voucherRepository) FindSeckillVoucherByID(ctx context.Context, voucherID int64) (*model.SeckillVoucher, error) {
	// TODO: Query tb_seckill_voucher by voucher_id.
	return nil, nil
}
