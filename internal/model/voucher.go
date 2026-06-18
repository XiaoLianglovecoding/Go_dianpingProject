package model

import "time"

type Voucher struct {
	ID          int64      `json:"id" gorm:"column:id;primaryKey"`
	ShopID      int64      `json:"shopId" gorm:"column:shop_id"`
	Title       string     `json:"title" gorm:"column:title"`
	SubTitle    string     `json:"subTitle" gorm:"column:sub_title"`
	Rules       string     `json:"rules" gorm:"column:rules"`
	PayValue    int64      `json:"payValue" gorm:"column:pay_value"`
	ActualValue int64      `json:"actualValue" gorm:"column:actual_value"`
	Type        int        `json:"type" gorm:"column:type"`
	Status      int        `json:"status" gorm:"column:status"`
	Stock       int        `json:"stock,omitempty" gorm:"-"`
	BeginTime   *time.Time `json:"beginTime,omitempty" gorm:"-"`
	EndTime     *time.Time `json:"endTime,omitempty" gorm:"-"`
	TimeFields
}

func (Voucher) TableName() string {
	return "tb_voucher"
}

type SeckillVoucher struct {
	VoucherID int64     `json:"voucherId" gorm:"column:voucher_id;primaryKey"`
	Stock     int       `json:"stock" gorm:"column:stock"`
	BeginTime time.Time `json:"beginTime" gorm:"column:begin_time"`
	EndTime   time.Time `json:"endTime" gorm:"column:end_time"`
	TimeFields
}

func (SeckillVoucher) TableName() string {
	return "tb_seckill_voucher"
}

type VoucherOrder struct {
	ID         int64      `json:"id" gorm:"column:id;primaryKey"`
	UserID     int64      `json:"userId" gorm:"column:user_id"`
	VoucherID  int64      `json:"voucherId" gorm:"column:voucher_id"`
	PayType    int        `json:"payType" gorm:"column:pay_type"`
	Status     int        `json:"status" gorm:"column:status"`
	CreateTime time.Time  `json:"createTime" gorm:"column:create_time"`
	PayTime    *time.Time `json:"payTime" gorm:"column:pay_time"`
	UseTime    *time.Time `json:"useTime" gorm:"column:use_time"`
	RefundTime *time.Time `json:"refundTime" gorm:"column:refund_time"`
	UpdateTime time.Time  `json:"updateTime" gorm:"column:update_time"`
}

func (VoucherOrder) TableName() string {
	return "tb_voucher_order"
}
