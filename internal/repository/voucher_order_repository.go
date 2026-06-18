package repository

import (
	"context"

	"hmdp-go/internal/model"

	"gorm.io/gorm"
)

type VoucherOrderRepository interface {
	CountByUserAndVoucher(ctx context.Context, userID int64, voucherID int64) (int64, error)
	SaveVoucherOrder(ctx context.Context, order *model.VoucherOrder) error
}

type voucherOrderRepository struct {
	db *gorm.DB
}

func NewVoucherOrderRepository(db *gorm.DB) VoucherOrderRepository {
	return &voucherOrderRepository{db: db}
}

func (r *voucherOrderRepository) CountByUserAndVoucher(ctx context.Context, userID int64, voucherID int64) (int64, error) {
	// TODO: Count existing voucher orders to enforce one user one order.
	return 0, nil
}

func (r *voucherOrderRepository) SaveVoucherOrder(ctx context.Context, order *model.VoucherOrder) error {
	// TODO: Insert voucher order into tb_voucher_order.
	return nil
}
