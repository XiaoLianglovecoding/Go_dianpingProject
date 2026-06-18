package handler

import (
	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/service"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	// uploadService 负责文件上传/删除业务逻辑。
	uploadService service.UploadService
}

// NewUploadHandler 创建上传 Handler。
func NewUploadHandler(uploadService service.UploadService) *UploadHandler {
	return &UploadHandler{uploadService: uploadService}
}

// UploadBlog 处理 POST /upload/blog，上传博客图片。
func (h *UploadHandler) UploadBlog(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		writeResult(c, result.Fail("missing upload file"))
		return
	}
	writeResult(c, h.uploadService.UploadBlogImage(c.Request.Context(), file))
}

// DeleteBlog 处理 GET /upload/blog/delete?name=xxx，删除博客图片。
func (h *UploadHandler) DeleteBlog(c *gin.Context) {
	name := c.Query("name")
	writeResult(c, h.uploadService.DeleteBlogImage(c.Request.Context(), name))
}
