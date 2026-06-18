package handler

import (
	"hmdp-go/internal/pkg/result"
	"hmdp-go/internal/service"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	uploadService service.UploadService
}

func NewUploadHandler(uploadService service.UploadService) *UploadHandler {
	return &UploadHandler{uploadService: uploadService}
}

func (h *UploadHandler) UploadBlog(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		writeResult(c, result.Fail("missing upload file"))
		return
	}
	writeResult(c, h.uploadService.UploadBlogImage(c.Request.Context(), file))
}

func (h *UploadHandler) DeleteBlog(c *gin.Context) {
	name := c.Query("name")
	writeResult(c, h.uploadService.DeleteBlogImage(c.Request.Context(), name))
}
