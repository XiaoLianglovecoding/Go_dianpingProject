package handler

import (
	"strconv"

	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/service"

	"github.com/gin-gonic/gin"
)

type FollowHandler struct {
	// followService 负责关注业务逻辑。
	followService service.FollowService
}

// NewFollowHandler 创建关注 Handler。
func NewFollowHandler(followService service.FollowService) *FollowHandler {
	return &FollowHandler{followService: followService}
}

// IsFollow 处理 GET /follow/or/not/:id，查询是否关注。
func (h *FollowHandler) IsFollow(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	writeResult(c, h.followService.IsFollow(c.Request.Context(), id))
}

// Follow 处理 PUT /follow/:id/:isFollow，关注或取消关注。
func (h *FollowHandler) Follow(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	isFollow, err := strconv.ParseBool(c.Param("isFollow"))
	if err != nil {
		writeResult(c, result.Fail("invalid path parameter: isFollow"))
		return
	}
	writeResult(c, h.followService.Follow(c.Request.Context(), id, isFollow))
}

// Common 处理 GET /follow/common/:id，查询共同关注。
func (h *FollowHandler) Common(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	writeResult(c, h.followService.Common(c.Request.Context(), id))
}
