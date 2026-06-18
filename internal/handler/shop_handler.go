package handler

import (
	"hmdp-go/internal/model"
	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/service"

	"github.com/gin-gonic/gin"
)

type ShopHandler struct {
	shopService service.ShopService
}

func NewShopHandler(shopService service.ShopService) *ShopHandler {
	return &ShopHandler{shopService: shopService}
}

func (h *ShopHandler) QueryShopByID(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	writeResult(c, h.shopService.QueryByID(c.Request.Context(), id))
}

func (h *ShopHandler) SaveShop(c *gin.Context) {
	var shop model.Shop
	if err := c.ShouldBindJSON(&shop); err != nil {
		writeResult(c, result.Fail("invalid shop request body"))
		return
	}
	writeResult(c, h.shopService.SaveShop(c.Request.Context(), shop))
}

func (h *ShopHandler) UpdateShop(c *gin.Context) {
	var shop model.Shop
	if err := c.ShouldBindJSON(&shop); err != nil {
		writeResult(c, result.Fail("invalid shop request body"))
		return
	}
	writeResult(c, h.shopService.UpdateShop(c.Request.Context(), shop))
}

func (h *ShopHandler) QueryShopByType(c *gin.Context) {
	typeID := int64(parseIntQuery(c, "typeId", 0))
	current := parseIntQuery(c, "current", 1)
	writeResult(c, h.shopService.QueryByType(c.Request.Context(), typeID, current))
}

func (h *ShopHandler) QueryShopByName(c *gin.Context) {
	name := c.Query("name")
	current := parseIntQuery(c, "current", 1)
	writeResult(c, h.shopService.QueryByName(c.Request.Context(), name, current))
}
