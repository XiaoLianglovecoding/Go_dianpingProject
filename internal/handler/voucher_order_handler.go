package handler

import (
	"hmdp-go/internal/service"

	"github.com/gin-gonic/gin"
)

type VoucherOrderHandler struct {
	voucherOrderService service.VoucherOrderService
}

func NewVoucherOrderHandler(voucherOrderService service.VoucherOrderService) *VoucherOrderHandler {
	return &VoucherOrderHandler{voucherOrderService: voucherOrderService}
}

func (h *VoucherOrderHandler) SeckillVoucher(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	writeResult(c, h.voucherOrderService.SeckillVoucher(c.Request.Context(), id))
}
