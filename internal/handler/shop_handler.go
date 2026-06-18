package handler

import (
	"hmdp-go/internal/model"
	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/service"

	"github.com/gin-gonic/gin"
)

type ShopHandler struct {
	// shopService 负责店铺业务逻辑。
	shopService service.ShopService
}

// NewShopHandler 创建店铺 Handler。
func NewShopHandler(shopService service.ShopService) *ShopHandler {
	return &ShopHandler{shopService: shopService}
}

// QueryShopByID 处理 GET /shop/:id，查询店铺详情。
func (h *ShopHandler) QueryShopByID(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	writeResult(c, h.shopService.QueryByID(c.Request.Context(), id))
}

// SaveShop 处理 POST /shop，新增店铺。
func (h *ShopHandler) SaveShop(c *gin.Context) {
	var shop model.Shop
	if err := c.ShouldBindJSON(&shop); err != nil {
		writeResult(c, result.Fail("invalid shop request body"))
		return
	}
	writeResult(c, h.shopService.SaveShop(c.Request.Context(), shop))
}

// UpdateShop 处理 PUT /shop，更新店铺。
func (h *ShopHandler) UpdateShop(c *gin.Context) {
	var shop model.Shop
	if err := c.ShouldBindJSON(&shop); err != nil {
		writeResult(c, result.Fail("invalid shop request body"))
		return
	}
	writeResult(c, h.shopService.UpdateShop(c.Request.Context(), shop))
}

// QueryShopByType 处理 GET /shop/of/type?typeId=1&current=1。
func (h *ShopHandler) QueryShopByType(c *gin.Context) {
	typeID := int64(parseIntQuery(c, "typeId", 0))
	current := parseIntQuery(c, "current", 1)
	writeResult(c, h.shopService.QueryByType(c.Request.Context(), typeID, current))
}

// QueryShopByName 处理 GET /shop/of/name?name=xxx&current=1。
func (h *ShopHandler) QueryShopByName(c *gin.Context) {
	name := c.Query("name")
	current := parseIntQuery(c, "current", 1)
	writeResult(c, h.shopService.QueryByName(c.Request.Context(), name, current))
}
