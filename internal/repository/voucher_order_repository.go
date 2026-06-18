package repository

import (
	"context"

	"hmdp-go/internal/model"

	"gorm.io/gorm"
)

type VoucherOrderRepository interface {
	// CountByUserAndVoucher 统计某个用户是否已经买过某张券。
	CountByUserAndVoucher(ctx context.Context, userID int64, voucherID int64) (int64, error)
	// SaveVoucherOrder 保存秒杀订单。
	SaveVoucherOrder(ctx context.Context, order *model.VoucherOrder) error
}

type voucherOrderRepository struct {
	// db 是 GORM 数据库连接对象。
	db *gorm.DB
}

// NewVoucherOrderRepository 创建订单 Repository。
func NewVoucherOrderRepository(db *gorm.DB) VoucherOrderRepository {
	return &voucherOrderRepository{db: db}
}

// CountByUserAndVoucher 后面用于实现“一人一单”。
func (r *voucherOrderRepository) CountByUserAndVoucher(ctx context.Context, userID int64, voucherID int64) (int64, error) {
	// TODO: Count existing voucher orders to enforce one user one order.
	return 0, nil
}

// SaveVoucherOrder 后面用于把抢购成功的订单写入 tb_voucher_order。
func (r *voucherOrderRepository) SaveVoucherOrder(ctx context.Context, order *model.VoucherOrder) error {
	// TODO: Insert voucher order into tb_voucher_order.
	return nil
}
