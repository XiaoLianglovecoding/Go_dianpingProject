package repository

import (
	"context"
	"errors"
	"time"

	"hmdp-go/internal/model"

	"gorm.io/gorm"
)

type VoucherOrderRepository interface {
	// CountByUserAndVoucher 统计某个用户是否已经买过某张券。
	CountByUserAndVoucher(ctx context.Context, userID int64, voucherID int64) (int64, error)
	// SaveVoucherOrder 保存秒杀订单。
	SaveVoucherOrder(ctx context.Context, order *model.VoucherOrder) error
	//创建订单方法
	CreateSeckillOrder(ctx context.Context, order *model.VoucherOrder) error
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
	var count int64

	// 对应 SQL: SELECT count(*) FROM tb_voucher_order WHERE user_id = ? AND voucher_id = ?
	err := r.db.WithContext(ctx).
		Model(&model.VoucherOrder{}). // 告诉 GORM 我们要查哪张表（根据结构体映射）
		Where("user_id = ? AND voucher_id = ?", userID, voucherID).
		Count(&count). // 把查到的数量塞进 count 变量里
		Error

	return count, err
}

// SaveVoucherOrder 后面用于把抢购成功的订单写入 tb_voucher_order。
func (r *voucherOrderRepository) SaveVoucherOrder(ctx context.Context, order *model.VoucherOrder) error {
	// 使用 GORM 创建记录
	return r.db.WithContext(ctx).Create(order).Error
}

// 异步处理,MySQL创建订单
func (r *voucherOrderRepository) CreateSeckillOrder(ctx context.Context, order *model.VoucherOrder) error {
	if order == nil {
		return errors.New("empty voucher order")
	}

	now := time.Now()
	if order.CreateTime.IsZero() {
		order.CreateTime = now
	}
	if order.UpdateTime.IsZero() {
		order.UpdateTime = now
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. MySQL 兜底检查一人一单
		var count int64
		if err := tx.Model(&model.VoucherOrder{}).
			Where("user_id = ? AND voucher_id = ?", order.UserID, order.VoucherID).
			Count(&count).Error; err != nil {
			return err
		}

		// 如果消息重复消费，数据库里已经有订单了，直接返回 nil，让外层 ACK。
		if count > 0 {
			return nil
		}

		// 2. MySQL 乐观锁扣库存
		res := tx.Model(&model.SeckillVoucher{}).
			Where("voucher_id = ? AND stock > 0", order.VoucherID).
			UpdateColumn("stock", gorm.Expr("stock - 1"))

		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("stock not enough")
		}

		// 3. 创建订单
		return tx.Create(order).Error
	})
}
