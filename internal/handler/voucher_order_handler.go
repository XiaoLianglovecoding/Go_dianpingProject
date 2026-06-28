package handler

import (
	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/pkg/userutils"
	"hmdp-go/internal/service"
	"net/http"

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
	// 2. 获取当前登录用户 ID (必须要登录才能抢购)
	userDTO, err := userutils.GetUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, result.Fail("未登录或登录已过期"))
		return
	}
	writeResult(c, h.voucherOrderService.SeckillVoucher(c.Request.Context(), id, userDTO.ID))
}
