package repository

import "gorm.io/gorm"

// Repositories 集中保存所有 Repository。
//
// Repository 层只负责“怎么和数据库打交道”，不写业务规则。
// 这样 Service 层只需要依赖接口，不需要关心 GORM 的细节。
type Repositories struct {
	User         UserRepository
	Shop         ShopRepository
	ShopType     ShopTypeRepository
	Blog         BlogRepository
	Follow       FollowRepository
	Voucher      VoucherRepository
	VoucherOrder VoucherOrderRepository
}

// NewRepositories 把同一个 *gorm.DB 注入到每个 Repository。
//
// db 是数据库连接对象；每个 Repository 会用它查询自己的表。
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:         NewUserRepository(db),
		Shop:         NewShopRepository(db),
		ShopType:     NewShopTypeRepository(db),
		Blog:         NewBlogRepository(db),
		Follow:       NewFollowRepository(db),
		Voucher:      NewVoucherRepository(db),
		VoucherOrder: NewVoucherOrderRepository(db),
	}
}
