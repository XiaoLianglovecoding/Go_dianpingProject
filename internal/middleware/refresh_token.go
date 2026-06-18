package middleware

import "github.com/gin-gonic/gin"

// RefreshTokenMiddleware 是 token 刷新中间件。
//
// Java 版有类似 RefreshTokenInterceptor：
// 如果请求带了合法 token，就刷新 Redis 里的 token 过期时间，
// 并把用户信息放到请求上下文，后续 Handler/Service 可以取。
func RefreshTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Refresh login token TTL when a valid token exists.
		c.Next()
	}
}
