package handler

import "hmdp-go/internal/service"

type Handlers struct {
	User         *UserHandler
	Shop         *ShopHandler
	ShopType     *ShopTypeHandler
	Blog         *BlogHandler
	Follow       *FollowHandler
	Upload       *UploadHandler
	Voucher      *VoucherHandler
	VoucherOrder *VoucherOrderHandler
}

func NewHandlers(services *service.Services) *Handlers {
	return &Handlers{
		User:         NewUserHandler(services.User),
		Shop:         NewShopHandler(services.Shop),
		ShopType:     NewShopTypeHandler(services.ShopType),
		Blog:         NewBlogHandler(services.Blog),
		Follow:       NewFollowHandler(services.Follow),
		Upload:       NewUploadHandler(services.Upload),
		Voucher:      NewVoucherHandler(services.Voucher),
		VoucherOrder: NewVoucherOrderHandler(services.VoucherOrder),
	}
}
