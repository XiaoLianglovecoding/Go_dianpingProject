package repository

import (
	"context"

	"hmdp-go/internal/model"

	"gorm.io/gorm"
)

type ShopRepository interface {
	FindShopByID(ctx context.Context, id int64) (*model.Shop, error)
	FindShopsByType(ctx context.Context, typeID int64, current int) ([]model.Shop, error)
	FindShopsByName(ctx context.Context, name string, current int) ([]model.Shop, error)
	SaveShop(ctx context.Context, shop *model.Shop) error
	UpdateShop(ctx context.Context, shop *model.Shop) error
}

type shopRepository struct {
	db *gorm.DB
}

func NewShopRepository(db *gorm.DB) ShopRepository {
	return &shopRepository{db: db}
}

func (r *shopRepository) FindShopByID(ctx context.Context, id int64) (*model.Shop, error) {
	// TODO: Query tb_shop by id.
	return nil, nil
}

func (r *shopRepository) FindShopsByType(ctx context.Context, typeID int64, current int) ([]model.Shop, error) {
	// TODO: Query tb_shop by type_id with pagination.
	return nil, nil
}

func (r *shopRepository) FindShopsByName(ctx context.Context, name string, current int) ([]model.Shop, error) {
	// TODO: Query tb_shop by name keyword with pagination.
	return nil, nil
}

func (r *shopRepository) SaveShop(ctx context.Context, shop *model.Shop) error {
	// TODO: Insert a new shop into tb_shop.
	return nil
}

func (r *shopRepository) UpdateShop(ctx context.Context, shop *model.Shop) error {
	// TODO: Update tb_shop and later evict cache:shop:{id}.
	return nil
}
