package userutils

import (
	"errors"
	"hmdp-go/internal/dto"

	"github.com/gin-gonic/gin"
)

// GetUser 从 Gin Context 中安全获取当前登录用户。
//
// 作用等同于 Java 中的 UserHolder.getUser()。
// 凡是经过 LoginInterceptor 拦截器的受保护请求，调用此方法必定能拿到合法的 UserDTO。
func GetUser(c *gin.Context) (dto.UserDTO, error) {
	// 1. 安全提取，防止 Panic
	userObj, exists := c.Get("user")
	if !exists {
		// 正常情况下，如果有保安(LoginInterceptor)守着，这里是不会触发的
		// 但为了极端的安全性，这里依然要做防御性返回
		return dto.UserDTO{}, errors.New("用户未登录或状态已失效")
	}

	// 2. 类型断言
	userDTO, ok := userObj.(dto.UserDTO)
	if !ok {
		return dto.UserDTO{}, errors.New("系统内部错误：用户信息解析异常")
	}

	return userDTO, nil
}

// GetUserID 尝试获取当前登录用户的 ID。
//
// 适用于“可选登录”的公开接口。如果未登录或类型解析失败，均安全地返回 0（代表游客）。
func GetUserID(c *gin.Context) int64 {
	// 直接复用 GetUser 的安全校验逻辑
	user, err := GetUser(c)
	if err != nil {
		return 0 // 获取失败或未登录，直接降级为游客
	}
	return user.ID
}
