package repository

import (
	"context"

	"hmdp-go/internal/model"

	"gorm.io/gorm"
)

type ShopRepository interface {
	// FindShopByID 根据店铺 id 查询店铺详情。
	FindShopByID(ctx context.Context, id int64) (*model.Shop, error)
	// FindShopsByType 根据分类分页查询店铺。
	FindShopsByType(ctx context.Context, typeID int64, current int) ([]model.Shop, error)
	// FindShopsByName 根据关键词搜索店铺。
	FindShopsByName(ctx context.Context, name string, current int) ([]model.Shop, error)
	// SaveShop 新增店铺。
	SaveShop(ctx context.Context, shop *model.Shop) error
	// UpdateShop 更新店铺。
	UpdateShop(ctx context.Context, shop *model.Shop) error
}

type shopRepository struct {
	// db 是 GORM 数据库连接对象。
	db *gorm.DB
}

// NewShopRepository 创建店铺 Repository。
func NewShopRepository(db *gorm.DB) ShopRepository {
	return &shopRepository{db: db}
}

// FindShopByID 后面会查询 tb_shop，并配合 Redis 做缓存。
func (r *shopRepository) FindShopByID(ctx context.Context, id int64) (*model.Shop, error) {
	// TODO: Query tb_shop by id.
	return nil, nil
}

// FindShopsByType 后面用于点击首页分类后的店铺列表页。
func (r *shopRepository) FindShopsByType(ctx context.Context, typeID int64, current int) ([]model.Shop, error) {
	// TODO: Query tb_shop by type_id with pagination.
	return nil, nil
}

// FindShopsByName 后面用于顶部搜索框搜索店铺。
func (r *shopRepository) FindShopsByName(ctx context.Context, name string, current int) ([]model.Shop, error) {
	// TODO: Query tb_shop by name keyword with pagination.
	return nil, nil
}

// SaveShop 后面用于后台新增店铺。
func (r *shopRepository) SaveShop(ctx context.Context, shop *model.Shop) error {
	// TODO: Insert a new shop into tb_shop.
	return nil
}

// UpdateShop 后面更新店铺后，还需要删除 Redis 里的旧缓存。
func (r *shopRepository) UpdateShop(ctx context.Context, shop *model.Shop) error {
	// TODO: Update tb_shop and later evict cache:shop:{id}.
	return nil
}
