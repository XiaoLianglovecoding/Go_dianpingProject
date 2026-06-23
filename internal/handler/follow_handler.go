package handler

import (
	"net/http"
	"strconv"

	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/pkg/userutils"
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
	// 1. 获取路径参数：目标用户 ID (你要看谁的主页)
	targetUserID, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}

	userDTO, err := userutils.GetUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, result.Fail(err.Error()))
		return
	}
	currentUserID := userDTO.ID
	writeResult(c, h.followService.IsFollow(c.Request.Context(), targetUserID, currentUserID))
}

// Follow 处理 PUT /follow/:id/:isFollow，关注或取消关注。
func (h *FollowHandler) Follow(c *gin.Context) {
	// 1. 获取目标用户 ID
	targetUserID, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	// 2. 获取 isFollow 布尔值参数
	isFollow, err := strconv.ParseBool(c.Param("isFollow"))
	if err != nil {
		writeResult(c, result.Fail("invalid path parameter: isFollow"))
		return
	}

	//获取当前登录用户ID
	userDTO, err := userutils.GetUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, result.Fail(err.Error()))
		return
	}
	currentUserID := userDTO.ID
	writeResult(c, h.followService.Follow(c.Request.Context(), currentUserID, targetUserID, isFollow))
}

// Common 处理 GET /follow/common/:id，查询共同关注。
func (h *FollowHandler) Common(c *gin.Context) {
	otherUserID, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	//获取当前登录用户ID
	userDTO, err := userutils.GetUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, result.Fail(err.Error()))
		return
	}
	currentUserID := userDTO.ID
	writeResult(c, h.followService.Common(c.Request.Context(), currentUserID, otherUserID))
}
