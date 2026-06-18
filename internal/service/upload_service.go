package service

import (
	"context"
	"mime/multipart"

	"hmdp-go/internal/pkg/result"
)

type UploadService interface {
	// UploadBlogImage 上传博客图片。
	UploadBlogImage(ctx context.Context, file *multipart.FileHeader) result.Result
	// DeleteBlogImage 删除博客图片。
	DeleteBlogImage(ctx context.Context, name string) result.Result
}

type uploadService struct{}

// NewUploadService 创建上传 Service。
func NewUploadService() UploadService {
	return &uploadService{}
}

// UploadBlogImage 后面会把图片保存到静态资源目录或专门的上传目录。
func (s *uploadService) UploadBlogImage(ctx context.Context, file *multipart.FileHeader) result.Result {
	// TODO: Store uploaded blog image under the frontend/static image directory or a Go-owned upload directory.
	return result.Fail("TODO: upload blog image")
}

// DeleteBlogImage 后面会删除指定图片。
//
// 实现时要注意路径安全，不能让用户传 ../ 删除任意文件。
func (s *uploadService) DeleteBlogImage(ctx context.Context, name string) result.Result {
	// TODO: Delete uploaded blog image by relative path after validating it is inside the upload directory.
	return result.Fail("TODO: delete blog image")
}
