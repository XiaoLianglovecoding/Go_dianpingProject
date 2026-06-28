package repository

import (
	"context"
	"errors"

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

	//DeductStock扣库存
	DeductStock(ctx context.Context, voucherID int64) error

	// SaveSeckillVoucherTx 开启事务，同时保存普通券和秒杀券信息
	SaveSeckillVoucherTx(ctx context.Context, voucher *model.Voucher, seckillVoucher *model.SeckillVoucher) error
}

type voucherRepository struct {
	// db 是 GORM 数据库连接对象。
	db *gorm.DB
}

// NewVoucherRepository 创建优惠券 Repository。
func NewVoucherRepository(db *gorm.DB) VoucherRepository {
	return &voucherRepository{db: db}
}

// SaveVoucher 用于插入 tb_voucher。
func (r *voucherRepository) SaveVoucher(ctx context.Context, voucher *model.Voucher) error {
	if voucher == nil {
		return errors.New("empty voucher")
	}
	return r.db.WithContext(ctx).Create(voucher).Error
}

// SaveSeckillVoucher 用于插入 tb_seckill_voucher。
func (r *voucherRepository) SaveSeckillVoucher(ctx context.Context, voucher *model.SeckillVoucher) error {
	if voucher == nil {
		return errors.New("empty seckill voucher")
	}
	return r.db.WithContext(ctx).Create(voucher).Error
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

// FindSeckillVoucherByID 用于秒杀下单前检查库存和时间。
func (r *voucherRepository) FindSeckillVoucherByID(ctx context.Context, voucherID int64) (*model.SeckillVoucher, error) {
	var voucher model.SeckillVoucher
	err := r.db.WithContext(ctx).
		Where("voucher_id = ?", voucherID).
		First(&voucher).Error
	if err != nil {
		return nil, err
	}
	return &voucher, nil
}

// DeductStock 扣减库存（带乐观锁防超卖）
func (r *voucherRepository) DeductStock(ctx context.Context, voucherID int64) error {
	// 对应 SQL: UPDATE tb_seckill_voucher SET stock = stock - 1 WHERE voucher_id = ? AND stock > 0

	// 1. 执行 Update，并将返回结果暂存到 res 变量中
	res := r.db.WithContext(ctx).
		Table("tb_seckill_voucher").                      // 注意：这里写你真实的表名，或者使用 Model(&model.SeckillVoucher{})
		Where("voucher_id = ? AND stock > 0", voucherID). // 【核心魔法】：加上 stock > 0 的条件
		Update("stock", gorm.Expr("stock - 1"))
	// 2. 判断有没有发生系统级别的 SQL 错误（比如数据库断开、语法错误）
	if res.Error != nil {
		return res.Error
	}
	// 3. 【最重要的一步】：判断受影响的行数 (RowsAffected)
	if res.RowsAffected == 0 {
		// 如果 RowsAffected == 0，说明 SQL 执行成功了，但是没有找到满足条件的数据！
		// 为什么找不到？因为 stock 已经被别人扣到 0 了，不满足 stock > 0 的条件了。
		// 这时候我们要人为抛出一个错误，告诉 Service 层：“扣减失败了，没库存了！”
		return errors.New("扣减库存失败，手慢了，库存已被抢空")
	}
	return nil
}

// SaveSeckillVoucherTx 实现事务级保存
func (r *voucherRepository) SaveSeckillVoucherTx(ctx context.Context, voucher *model.Voucher, seckillVoucher *model.SeckillVoucher) error {
	// GORM 的 Transaction 方法会自动开启事务、处理 Commit 和 Rollback
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 插入 tb_voucher，插入后 GORM 会将自增 ID 赋值给 voucher.ID
		if err := tx.Create(voucher).Error; err != nil {
			return err // 返回错误，事务自动回滚
		}

		// 2. 将刚生成的 voucher.ID 赋给 seckillVoucher
		seckillVoucher.VoucherID = voucher.ID

		// 3. 插入 tb_seckill_voucher
		if err := tx.Create(seckillVoucher).Error; err != nil {
			return err // 返回错误，事务自动回滚，连带第一步的 tb_voucher 一起撤销
		}

		// 两个都成功，返回 nil，事务自动提交
		return nil
	})
}
