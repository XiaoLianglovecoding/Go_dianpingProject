package handler

import (
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
	//获取当前登录用户ID
	userDTO, err := userutils.GetUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, result.Fail(err.Error()))
		return
	}
	currentUserID := userDTO.ID
	writeResult(c, h.blogService.SaveBlog(c.Request.Context(), blog, currentUserID))
}

// QueryByID 处理 GET /blog/:id，查询博客详情。
func (h *BlogHandler) QueryByID(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	// (观看者 viewerID，没登录就是 0)安全提取与兜底，彻底消除 Panic 隐患
	userID := userutils.GetUserID(c)

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

	// 【强登录接口】：使用 GetUser 严格校验，失败直接返回 401 阻断
	userDTO, err := userutils.GetUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, result.Fail(err.Error()))
		return
	}

	// 看自己的主页，作者和观看者都是自己
	userID := userDTO.ID
	writeResult(c, h.blogService.QueryBlogByUserID(c.Request.Context(), userID, current, userID))
}

// QueryHotBlog 处理 GET /blog/hot?current=1，首页热门博客列表。
func (h *BlogHandler) QueryHotBlog(c *gin.Context) {
	current := parseIntQuery(c, "current", 1)

	// 安全提取游客 ID(观看者 viewerID，没登录就是 0)
	userID := userutils.GetUserID(c)
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

	// 2. 温柔提取当前登录用户 ID (观看者 viewerID，没登录就是 0)
	// 安全提取观看者 viewerID
	viewerID := userutils.GetUserID(c)

	writeResult(c, h.blogService.QueryBlogByUserID(c.Request.Context(), authorID, current, viewerID))
}

// Handler 增加查询关注流
func (h *BlogHandler) QueryBlogOfFollow(c *gin.Context) {
	lastID := parseInt64Query(c, "lastId", 0)
	offset := parseIntQuery(c, "offset", 0)

	userDTO, err := userutils.GetUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, result.Fail(err.Error()))
		return
	}

	writeResult(c, h.blogService.QueryBlogOfFollow(
		c.Request.Context(),
		lastID,
		offset,
		userDTO.ID,
	))
}
