package handler

import (
	"hmdp-go/internal/model"
	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/service"

	"github.com/gin-gonic/gin"
)

type VoucherHandler struct {
	voucherService service.VoucherService
}

func NewVoucherHandler(voucherService service.VoucherService) *VoucherHandler {
	return &VoucherHandler{voucherService: voucherService}
}

func (h *VoucherHandler) AddVoucher(c *gin.Context) {
	var voucher model.Voucher
	if err := c.ShouldBindJSON(&voucher); err != nil {
		writeResult(c, result.Fail("invalid voucher request body"))
		return
	}
	writeResult(c, h.voucherService.AddVoucher(c.Request.Context(), voucher))
}

func (h *VoucherHandler) AddSeckillVoucher(c *gin.Context) {
	var voucher model.Voucher
	if err := c.ShouldBindJSON(&voucher); err != nil {
		writeResult(c, result.Fail("invalid voucher request body"))
		return
	}
	writeResult(c, h.voucherService.AddSeckillVoucher(c.Request.Context(), voucher))
}

func (h *VoucherHandler) QueryVoucherOfShop(c *gin.Context) {
	shopID, ok := parseInt64Param(c, "shopId")
	if !ok {
		return
	}
	writeResult(c, h.voucherService.QueryVoucherOfShop(c.Request.Context(), shopID))
}
