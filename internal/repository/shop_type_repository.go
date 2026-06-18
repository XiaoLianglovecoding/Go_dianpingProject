package repository

import (
	"context"

	"hmdp-go/internal/model"

	"gorm.io/gorm"
)

type ShopTypeRepository interface {
	FindShopTypes(ctx context.Context) ([]model.ShopType, error)
}

type shopTypeRepository struct {
	db *gorm.DB
}

func NewShopTypeRepository(db *gorm.DB) ShopTypeRepository {
	return &shopTypeRepository{db: db}
}

func (r *shopTypeRepository) FindShopTypes(ctx context.Context) ([]model.ShopType, error) {
	// TODO: Query tb_shop_type ordered by sort.
	return nil, nil
}
