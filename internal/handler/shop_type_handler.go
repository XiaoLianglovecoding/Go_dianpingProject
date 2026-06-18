package handler

import (
	"hmdp-go/internal/service"

	"github.com/gin-gonic/gin"
)

type ShopTypeHandler struct {
	shopTypeService service.ShopTypeService
}

func NewShopTypeHandler(shopTypeService service.ShopTypeService) *ShopTypeHandler {
	return &ShopTypeHandler{shopTypeService: shopTypeService}
}

func (h *ShopTypeHandler) List(c *gin.Context) {
	writeResult(c, h.shopTypeService.List(c.Request.Context()))
}
