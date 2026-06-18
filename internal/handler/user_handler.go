package handler

import (
	"hmdp-go/internal/dto"
	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	// userService 负责用户业务逻辑。
	userService service.UserService
}

// NewUserHandler 创建用户 Handler。
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// SendCode 处理 POST /user/code?phone=xxx。
func (h *UserHandler) SendCode(c *gin.Context) {
	phone := c.Query("phone")
	writeResult(c, h.userService.SendCode(c.Request.Context(), phone))
}

// Login 处理 POST /user/login。
//
// 前端会提交 JSON，比如 {"phone":"...","code":"..."}。
func (h *UserHandler) Login(c *gin.Context) {
	var form dto.LoginFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		writeResult(c, result.Fail("invalid login request body"))
		return
	}
	writeResult(c, h.userService.Login(c.Request.Context(), form))
}

// Logout 处理 POST /user/logout。
func (h *UserHandler) Logout(c *gin.Context) {
	writeResult(c, h.userService.Logout(c.Request.Context()))
}

// Me 处理 GET /user/me，查询当前登录用户。
func (h *UserHandler) Me(c *gin.Context) {
	writeResult(c, h.userService.Me(c.Request.Context()))
}

// QueryUserByID 处理 GET /user/:id。
func (h *UserHandler) QueryUserByID(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	writeResult(c, h.userService.QueryUserByID(c.Request.Context(), id))
}

// QueryUserInfo 处理 GET /user/info/:id。
func (h *UserHandler) QueryUserInfo(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	writeResult(c, h.userService.QueryUserInfo(c.Request.Context(), id))
}

// Sign 处理 POST /user/sign。
func (h *UserHandler) Sign(c *gin.Context) {
	writeResult(c, h.userService.Sign(c.Request.Context()))
}

// SignCount 处理 GET /user/sign/count。
func (h *UserHandler) SignCount(c *gin.Context) {
	writeResult(c, h.userService.SignCount(c.Request.Context()))
}
