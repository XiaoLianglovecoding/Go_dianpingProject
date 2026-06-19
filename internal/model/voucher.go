package model

import "time"

// Voucher 对应数据库表 tb_voucher，表示优惠券基础信息。
type Voucher struct {
	ID          int64      `json:"id" gorm:"column:id;primaryKey"`
	ShopID      int64      `json:"shopId" gorm:"column:shop_id"`           // 所属店铺 id。
	Title       string     `json:"title" gorm:"column:title"`              // 优惠券标题。
	SubTitle    string     `json:"subTitle" gorm:"column:sub_title"`       // 副标题。
	Rules       string     `json:"rules" gorm:"column:rules"`              // 使用规则。
	PayValue    int64      `json:"payValue" gorm:"column:pay_value"`       // 支付金额，单位分。
	ActualValue int64      `json:"actualValue" gorm:"column:actual_value"` // 抵扣金额，单位分。
	Type        int        `json:"type" gorm:"column:type"`                // 类型：0 普通券，1 秒杀券。
	Status      int        `json:"status" gorm:"column:status"`            // 状态：上架/下架/过期等。
	Stock       int        `json:"stock,omitempty" gorm:"-"`               // 秒杀库存来自 tb_seckill_voucher，不在 tb_voucher 表里。
	BeginTime   *time.Time `json:"beginTime,omitempty" gorm:"-"`           // 秒杀开始时间，查询时手动补充。
	EndTime     *time.Time `json:"endTime,omitempty" gorm:"-"`             // 秒杀结束时间，查询时手动补充。

	// 👇 【新增这一行】告诉 GORM 这两张表是一对一关系，方便我们用 Preload 查数据
	// json:"-" 的作用是：返回给前端时隐藏这个嵌套的结构，防止前端报错。
	SeckillVoucher *SeckillVoucher `json:"-" gorm:"foreignKey:VoucherID;references:ID"`
	TimeFields
}

// TableName 告诉 GORM：Voucher 对应 tb_voucher 表。
func (Voucher) TableName() string {
	return "tb_voucher"
}

// SeckillVoucher 对应数据库表 tb_seckill_voucher，保存秒杀券额外信息。
type SeckillVoucher struct {
	VoucherID int64     `json:"voucherId" gorm:"column:voucher_id;primaryKey"`
	Stock     int       `json:"stock" gorm:"column:stock"`          // 秒杀库存。
	BeginTime time.Time `json:"beginTime" gorm:"column:begin_time"` // 秒杀开始时间。
	EndTime   time.Time `json:"endTime" gorm:"column:end_time"`     // 秒杀结束时间。
	TimeFields
}

// TableName 告诉 GORM：SeckillVoucher 对应 tb_seckill_voucher 表。
func (SeckillVoucher) TableName() string {
	return "tb_seckill_voucher"
}

// VoucherOrder 对应数据库表 tb_voucher_order，表示用户下单记录。
type VoucherOrder struct {
	ID         int64      `json:"id" gorm:"column:id;primaryKey"`
	UserID     int64      `json:"userId" gorm:"column:user_id"`       // 下单用户 id。
	VoucherID  int64      `json:"voucherId" gorm:"column:voucher_id"` // 被购买的优惠券 id。
	PayType    int        `json:"payType" gorm:"column:pay_type"`     // 支付方式。
	Status     int        `json:"status" gorm:"column:status"`        // 订单状态。
	CreateTime time.Time  `json:"createTime" gorm:"column:create_time"`
	PayTime    *time.Time `json:"payTime" gorm:"column:pay_time"`
	UseTime    *time.Time `json:"useTime" gorm:"column:use_time"`
	RefundTime *time.Time `json:"refundTime" gorm:"column:refund_time"`
	UpdateTime time.Time  `json:"updateTime" gorm:"column:update_time"`
}

// TableName 告诉 GORM：VoucherOrder 对应 tb_voucher_order 表。
func (VoucherOrder) TableName() string {
	return "tb_voucher_order"
}
