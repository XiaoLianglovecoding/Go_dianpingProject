package handler

import (
	"hmdp-go/internal/model"
	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/service"

	"github.com/gin-gonic/gin"
)

type VoucherHandler struct {
	// voucherService 负责优惠券业务逻辑。
	voucherService service.VoucherService
}

// NewVoucherHandler 创建优惠券 Handler。
func NewVoucherHandler(voucherService service.VoucherService) *VoucherHandler {
	return &VoucherHandler{voucherService: voucherService}
}

// AddVoucher 处理 POST /voucher，新增普通优惠券。
func (h *VoucherHandler) AddVoucher(c *gin.Context) {
	var voucher model.Voucher
	if err := c.ShouldBindJSON(&voucher); err != nil {
		writeResult(c, result.Fail("invalid voucher request body"))
		return
	}
	writeResult(c, h.voucherService.AddVoucher(c.Request.Context(), voucher))
}

// AddSeckillVoucher 处理 POST /voucher/seckill，新增秒杀券。
func (h *VoucherHandler) AddSeckillVoucher(c *gin.Context) {
	var voucher model.Voucher
	if err := c.ShouldBindJSON(&voucher); err != nil {
		writeResult(c, result.Fail("invalid voucher request body"))
		return
	}
	writeResult(c, h.voucherService.AddSeckillVoucher(c.Request.Context(), voucher))
}

// QueryVoucherOfShop 处理 GET /voucher/list/:shopId。
func (h *VoucherHandler) QueryVoucherOfShop(c *gin.Context) {
	shopID, ok := parseInt64Param(c, "shopId")
	if !ok {
		return
	}
	writeResult(c, h.voucherService.QueryVoucherOfShop(c.Request.Context(), shopID))
}
