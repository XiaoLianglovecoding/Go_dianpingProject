package handler

import (
	"hmdp-go/internal/dto"
	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) SendCode(c *gin.Context) {
	phone := c.Query("phone")
	writeResult(c, h.userService.SendCode(c.Request.Context(), phone))
}

func (h *UserHandler) Login(c *gin.Context) {
	var form dto.LoginFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		writeResult(c, result.Fail("invalid login request body"))
		return
	}
	writeResult(c, h.userService.Login(c.Request.Context(), form))
}

func (h *UserHandler) Logout(c *gin.Context) {
	writeResult(c, h.userService.Logout(c.Request.Context()))
}

func (h *UserHandler) Me(c *gin.Context) {
	writeResult(c, h.userService.Me(c.Request.Context()))
}

func (h *UserHandler) QueryUserByID(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	writeResult(c, h.userService.QueryUserByID(c.Request.Context(), id))
}

func (h *UserHandler) QueryUserInfo(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	writeResult(c, h.userService.QueryUserInfo(c.Request.Context(), id))
}

func (h *UserHandler) Sign(c *gin.Context) {
	writeResult(c, h.userService.Sign(c.Request.Context()))
}

func (h *UserHandler) SignCount(c *gin.Context) {
	writeResult(c, h.userService.SignCount(c.Request.Context()))
}
