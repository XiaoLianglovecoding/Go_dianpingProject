package handler

import (
	"net/http"
	"strconv"

	"hmdp-go/internal/pkg/result"

	"github.com/gin-gonic/gin"
)

func writeResult(c *gin.Context, res result.Result) {
	c.JSON(http.StatusOK, res)
}

func parseInt64Param(c *gin.Context, name string) (int64, bool) {
	value := c.Param(name)
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		writeResult(c, result.Fail("invalid path parameter: "+name))
		return 0, false
	}
	return id, true
}

func parseIntQuery(c *gin.Context, name string, defaultValue int) int {
	value := c.Query(name)
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}
