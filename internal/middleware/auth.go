package middleware

import "github.com/gin-gonic/gin"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Read authorization header, load login:token:{token} from Redis, and put UserDTO into context.
		c.Next()
	}
}
