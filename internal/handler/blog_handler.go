package handler

import (
	"hmdp-go/internal/dto"
	"hmdp-go/internal/model"
	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/pkg/userutils"
	"hmdp-go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BlogHandler struct {
	// blogService 负责博客业务逻辑。
	blogService service.BlogService
}

// NewBlogHandler 创建博客 Handler。
func NewBlogHandler(blogService service.BlogService) *BlogHandler {
	return &BlogHandler{blogService: blogService}
}

// SaveBlog 处理 POST /blog。
//
// ShouldBindJSON 会把请求体 JSON 绑定到 model.Blog 结构体。
func (h *BlogHandler) SaveBlog(c *gin.Context) {
	var blog model.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		writeResult(c, result.Fail("invalid blog request body"))
		return
	}
	writeResult(c, h.blogService.SaveBlog(c.Request.Context(), blog))
}

// QueryByID 处理 GET /blog/:id，查询博客详情。
func (h *BlogHandler) QueryByID(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	// 【温柔提取】：不报错，查不到就是游客 (0)
	var userID int64 = 0
	if userObj, exists := c.Get("user"); exists {
		userID = userObj.(dto.UserDTO).ID
	}
	writeResult(c, h.blogService.QueryByID(c.Request.Context(), id, userID))
}

// LikeBlog 处理 PUT /blog/like/:id，点赞或取消点赞。
func (h *BlogHandler) LikeBlog(c *gin.Context) {
	// 1. 获取路径上的博客 ID
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}

	userDTO, err := userutils.GetUser(c)
	// 2. 严谨的错误处理：如果拿不到人，直接打回 401报错并 return 中断请求
	if err != nil {
		c.JSON(http.StatusUnauthorized, result.Fail(err.Error()))
		return
	}

	writeResult(c, h.blogService.LikeBlog(c.Request.Context(), id, userDTO.ID))
}

// QueryMyBlog 处理 GET /blog/of/me，查询我的博客。
func (h *BlogHandler) QueryMyBlog(c *gin.Context) {
	current := parseIntQuery(c, "current", 1)
	writeResult(c, h.blogService.QueryMyBlog(c.Request.Context(), current))
}

// QueryHotBlog 处理 GET /blog/hot?current=1，首页热门博客列表。
func (h *BlogHandler) QueryHotBlog(c *gin.Context) {
	current := parseIntQuery(c, "current", 1)

	// 【温柔提取】：不报错，查不到就是游客 (0)
	var userID int64 = 0
	if userObj, exists := c.Get("user"); exists {
		userID = userObj.(dto.UserDTO).ID
	}
	writeResult(c, h.blogService.QueryHotBlog(c.Request.Context(), current, userID))
}

// QueryBlogLikes 处理 GET /blog/likes/:id，查询点赞用户列表。
func (h *BlogHandler) QueryBlogLikes(c *gin.Context) {
	blogID, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	writeResult(c, h.blogService.QueryBlogLikes(c.Request.Context(), blogID))
}

// QueryBlogByUserID 处理 GET /blog/of/user?id=xxx，查询某个用户的博客。
func (h *BlogHandler) QueryBlogByUserID(c *gin.Context) {
	// 1. 获取目标用户的 ID (主页主人 authorID)
	authorID := int64(parseIntQuery(c, "id", 0))
	current := parseIntQuery(c, "current", 1)

	// 2. 【温柔提取】当前登录用户 ID (观看者 viewerID，没登录就是 0)
	var viewerID int64 = 0
	if userObj, exists := c.Get("user"); exists {
		viewerID = userObj.(dto.UserDTO).ID
	}
	writeResult(c, h.blogService.QueryBlogByUserID(c.Request.Context(), authorID, current, viewerID))
}
