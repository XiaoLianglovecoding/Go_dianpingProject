package handler

import (
	"hmdp-go/internal/model"
	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/service"

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
	writeResult(c, h.blogService.QueryByID(c.Request.Context(), id))
}

// LikeBlog 处理 PUT /blog/like/:id，点赞或取消点赞。
func (h *BlogHandler) LikeBlog(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	writeResult(c, h.blogService.LikeBlog(c.Request.Context(), id))
}

// QueryMyBlog 处理 GET /blog/of/me，查询我的博客。
func (h *BlogHandler) QueryMyBlog(c *gin.Context) {
	current := parseIntQuery(c, "current", 1)
	writeResult(c, h.blogService.QueryMyBlog(c.Request.Context(), current))
}

// QueryHotBlog 处理 GET /blog/hot?current=1，首页热门博客列表。
func (h *BlogHandler) QueryHotBlog(c *gin.Context) {
	current := parseIntQuery(c, "current", 1)
	writeResult(c, h.blogService.QueryHotBlog(c.Request.Context(), current))
}

// QueryBlogLikes 处理 GET /blog/likes/:id，查询点赞用户列表。
func (h *BlogHandler) QueryBlogLikes(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	writeResult(c, h.blogService.QueryBlogLikes(c.Request.Context(), id))
}

// QueryBlogByUserID 处理 GET /blog/of/user?id=xxx，查询某个用户的博客。
func (h *BlogHandler) QueryBlogByUserID(c *gin.Context) {
	id := int64(parseIntQuery(c, "id", 0))
	current := parseIntQuery(c, "current", 1)
	writeResult(c, h.blogService.QueryBlogByUserID(c.Request.Context(), id, current))
}
