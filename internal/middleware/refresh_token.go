package middleware

import "github.com/gin-gonic/gin"

func RefreshTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Refresh login token TTL when a valid token exists.
		c.Next()
	}
}
