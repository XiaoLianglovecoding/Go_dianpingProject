package handler

import (
	"hmdp-go/internal/model"
	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/service"

	"github.com/gin-gonic/gin"
)

type BlogHandler struct {
	blogService service.BlogService
}

func NewBlogHandler(blogService service.BlogService) *BlogHandler {
	return &BlogHandler{blogService: blogService}
}

func (h *BlogHandler) SaveBlog(c *gin.Context) {
	var blog model.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		writeResult(c, result.Fail("invalid blog request body"))
		return
	}
	writeResult(c, h.blogService.SaveBlog(c.Request.Context(), blog))
}

func (h *BlogHandler) QueryByID(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	writeResult(c, h.blogService.QueryByID(c.Request.Context(), id))
}

func (h *BlogHandler) LikeBlog(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	writeResult(c, h.blogService.LikeBlog(c.Request.Context(), id))
}

func (h *BlogHandler) QueryMyBlog(c *gin.Context) {
	current := parseIntQuery(c, "current", 1)
	writeResult(c, h.blogService.QueryMyBlog(c.Request.Context(), current))
}

func (h *BlogHandler) QueryHotBlog(c *gin.Context) {
	current := parseIntQuery(c, "current", 1)
	writeResult(c, h.blogService.QueryHotBlog(c.Request.Context(), current))
}

func (h *BlogHandler) QueryBlogLikes(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	writeResult(c, h.blogService.QueryBlogLikes(c.Request.Context(), id))
}

func (h *BlogHandler) QueryBlogByUserID(c *gin.Context) {
	id := int64(parseIntQuery(c, "id", 0))
	current := parseIntQuery(c, "current", 1)
	writeResult(c, h.blogService.QueryBlogByUserID(c.Request.Context(), id, current))
}
