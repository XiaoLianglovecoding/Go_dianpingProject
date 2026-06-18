package handler

import (
	"hmdp-go/internal/service"

	"github.com/gin-gonic/gin"
)

type VoucherOrderHandler struct {
	// voucherOrderService 负责优惠券订单业务逻辑。
	voucherOrderService service.VoucherOrderService
}

// NewVoucherOrderHandler 创建优惠券订单 Handler。
func NewVoucherOrderHandler(voucherOrderService service.VoucherOrderService) *VoucherOrderHandler {
	return &VoucherOrderHandler{voucherOrderService: voucherOrderService}
}

// SeckillVoucher 处理 POST /voucher-order/seckill/:id。
func (h *VoucherOrderHandler) SeckillVoucher(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	writeResult(c, h.voucherOrderService.SeckillVoucher(c.Request.Context(), id))
}
