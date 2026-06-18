package service

import (
	"context"
	"mime/multipart"

	"hmdp-go/internal/pkg/result"
)

type UploadService interface {
	UploadBlogImage(ctx context.Context, file *multipart.FileHeader) result.Result
	DeleteBlogImage(ctx context.Context, name string) result.Result
}

type uploadService struct{}

func NewUploadService() UploadService {
	return &uploadService{}
}

func (s *uploadService) UploadBlogImage(ctx context.Context, file *multipart.FileHeader) result.Result {
	// TODO: Store uploaded blog image under the frontend/static image directory or a Go-owned upload directory.
	return result.Fail("TODO: upload blog image")
}

func (s *uploadService) DeleteBlogImage(ctx context.Context, name string) result.Result {
	// TODO: Delete uploaded blog image by relative path after validating it is inside the upload directory.
	return result.Fail("TODO: delete blog image")
}
