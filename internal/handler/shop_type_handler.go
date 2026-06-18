package handler

import (
	"hmdp-go/internal/service"

	"github.com/gin-gonic/gin"
)

type ShopTypeHandler struct {
	// shopTypeService 负责店铺分类业务逻辑。
	shopTypeService service.ShopTypeService
}

// NewShopTypeHandler 创建店铺分类 Handler。
func NewShopTypeHandler(shopTypeService service.ShopTypeService) *ShopTypeHandler {
	return &ShopTypeHandler{shopTypeService: shopTypeService}
}

// List 处理 GET /shop-type/list。
//
// Handler 只负责接 HTTP 请求，然后把工作交给 Service。
func (h *ShopTypeHandler) List(c *gin.Context) {
	writeResult(c, h.shopTypeService.List(c.Request.Context()))
}
