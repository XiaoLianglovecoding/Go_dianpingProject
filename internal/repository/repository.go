package repository

import "gorm.io/gorm"

type Repositories struct {
	User         UserRepository
	Shop         ShopRepository
	ShopType     ShopTypeRepository
	Blog         BlogRepository
	Follow       FollowRepository
	Voucher      VoucherRepository
	VoucherOrder VoucherOrderRepository
}

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
