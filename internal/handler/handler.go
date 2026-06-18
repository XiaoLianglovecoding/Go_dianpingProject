package handler

import "hmdp-go/internal/service"

// Handlers 集中保存所有 Handler。
//
// Handler 层对应 Java 里的 Controller：
// 它只处理 HTTP 相关的事情，例如取参数、绑定 JSON、返回响应。
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

// NewHandlers 根据 Services 创建所有 Handler。
//
// 依赖方向是 Handler -> Service -> Repository。
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
