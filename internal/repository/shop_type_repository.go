package repository

import (
	"context"

	"hmdp-go/internal/model"

	"gorm.io/gorm"
)

type ShopTypeRepository interface {
	// FindShopTypes 查询所有店铺分类，按 sort 字段升序返回。
	FindShopTypes(ctx context.Context) ([]model.ShopType, error)
}

type shopTypeRepository struct {
	// db 是 GORM 数据库连接对象，Repository 层通过它操作 MySQL。
	db *gorm.DB
}

// NewShopTypeRepository 创建店铺分类 Repository。
func NewShopTypeRepository(db *gorm.DB) ShopTypeRepository {
	return &shopTypeRepository{db: db}
}

// FindShopTypes 查询 tb_shop_type 表。
//
// 对应 SQL 大概是:
// SELECT * FROM tb_shop_type ORDER BY sort ASC;
func (r *shopTypeRepository) FindShopTypes(ctx context.Context) ([]model.ShopType, error) {
	// shopTypes 是接收查询结果的切片，因为分类有多条。
	var shopTypes []model.ShopType

	// WithContext(ctx) 让数据库查询跟随请求生命周期；
	// Order("sort ASC") 表示按 sort 从小到大排序；
	// Find(&shopTypes) 表示查询多条记录并写入 shopTypes。
	err := r.db.WithContext(ctx).
		Order("sort ASC").
		Find(&shopTypes).Error

	if err != nil {
		return nil, err
	}
	return shopTypes, nil
}
