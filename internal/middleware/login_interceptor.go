package middleware

import (
	"hmdp-go/internal/pkg/result" // 替换为你自己的 result 包路径
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginInterceptor 第二层：登录校验拦截器
// 专门挂载到需要登录的路由上，如果不包含 user 信息，直接拦截报错
func LoginInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从 Context 中获取刚才 RefreshTokenMiddleware 放进去的 "user"
		_, exists := c.Get("user")

		if !exists {
			// 如果没拿到，说明这个人要么没带 Token，要么 Token 已经过期了
			// 直接中止请求 (相当于 Spring 拦截器 return false)，并返回 401 状态码
			c.AbortWithStatusJSON(http.StatusUnauthorized, result.Fail("用户未登录或身份已过期"))
			return
		}

		// 存在用户信息，说明鉴权通过，放行去执行具体的 Controller(Handler)
		c.Next()
	}
}
