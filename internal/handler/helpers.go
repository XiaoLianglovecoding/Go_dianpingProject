package handler

import (
	"net/http"
	"strconv"

	"hmdp-go/internal/pkg/result"

	"github.com/gin-gonic/gin"
)

// writeResult 统一把 result.Result 写成 JSON 响应。
//
// 当前项目为了兼容前端，业务成功/失败都返回 HTTP 200，
// 真正是否成功由 JSON 里的 success 字段判断。
func writeResult(c *gin.Context, res result.Result) {
	c.JSON(http.StatusOK, res)
}

// parseInt64Param 读取路径参数并转换成 int64。
//
// 例如 /blog/123 中的 123，可以用 parseInt64Param(c, "id") 读取。
// 第二个返回值 bool 表示是否解析成功。
func parseInt64Param(c *gin.Context, name string) (int64, bool) {
	value := c.Param(name)
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		writeResult(c, result.Fail("invalid path parameter: "+name))
		return 0, false
	}
	return id, true
}

// parseIntQuery 读取 query 参数并转换成 int。
//
// 例如 /blog/hot?current=1 中的 current。
// 如果没传或传错，就返回 defaultValue。
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
