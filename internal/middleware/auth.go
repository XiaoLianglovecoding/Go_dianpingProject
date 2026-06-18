package middleware

import "github.com/gin-gonic/gin"

// AuthMiddleware 是登录校验中间件。
//
// 中间件会在 Handler 之前执行。
// 后续实现时，这里会从请求头 authorization 读取 token，
// 再去 Redis 查 login:token:{token}，查到才允许继续访问。
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Read authorization header, load login:token:{token} from Redis, and put UserDTO into context.
		c.Next()
	}
}
