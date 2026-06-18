package handler

import (
	"strconv"

	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/service"

	"github.com/gin-gonic/gin"
)

type FollowHandler struct {
	followService service.FollowService
}

func NewFollowHandler(followService service.FollowService) *FollowHandler {
	return &FollowHandler{followService: followService}
}

func (h *FollowHandler) IsFollow(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	writeResult(c, h.followService.IsFollow(c.Request.Context(), id))
}

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

func (h *FollowHandler) Common(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	writeResult(c, h.followService.Common(c.Request.Context(), id))
}
