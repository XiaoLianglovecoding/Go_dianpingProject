package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"hmdp-go/internal/pkg/result"

	"hmdp-go/internal/config"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
)

type UploadService interface {
	// UploadBlogImage 上传博客图片。
	UploadBlogImage(ctx context.Context, file *multipart.FileHeader) result.Result
	// DeleteBlogImage 删除博客图片。
	DeleteBlogImage(ctx context.Context, name string) result.Result
}

type uploadService struct {
	ossCfg config.OSSConfig
}

// NewUploadService 创建上传 Service。
func NewUploadService(ossCfg config.OSSConfig) UploadService {
	return &uploadService{ossCfg: ossCfg}
}

// UploadBlogImage 后面会把图片保存到静态资源目录或专门的上传目录。
func (s *uploadService) UploadBlogImage(ctx context.Context, file *multipart.FileHeader) result.Result {
	_ = ctx

	if file == nil {
		return result.Fail("missing upload file")
	}
	if file.Size > 5*1024*1024 {
		return result.Fail("图片不能超过 5MB")
	}

	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return result.Fail("只能上传图片")
	}

	src, err := file.Open()
	if err != nil {
		return result.Fail("读取上传文件失败")
	}
	defer src.Close()

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" {
		ext = ".jpg"
	}

	objectName := buildOSSObjectName(s.ossCfg.Dir, ext)

	client, err := oss.New(
		s.ossCfg.Endpoint,
		s.ossCfg.AccessKeyID,
		s.ossCfg.AccessKeySecret,
	)
	if err != nil {
		return result.Fail("连接 OSS 失败")
	}

	bucket, err := client.Bucket(s.ossCfg.Bucket)
	if err != nil {
		return result.Fail("获取 OSS Bucket 失败")
	}

	err = bucket.PutObject(
		objectName,
		src,
		oss.ContentType(contentType),
	)
	if err != nil {
		return result.Fail("上传图片到 OSS 失败")
	}

	// 前端会自动拼接 "/imgs" 前缀，所以这里必须模仿 Java 版返回 "/blogs/xxx.jpg"。
	return result.OKWithData("/" + objectName)
}

func buildOSSObjectName(dir string, ext string) string {
	prefix := strings.Trim(dir, "/")
	datePath := time.Now().Format("2006/01/02")
	fileName := uuid.New().String() + ext

	if prefix == "" {
		return fmt.Sprintf("%s/%s", datePath, fileName)
	}
	return fmt.Sprintf("%s/%s/%s", prefix, datePath, fileName)
}

// DeleteBlogImage 后面会删除指定图片。
//
// 实现时要注意路径安全，不能让用户传 ../ 删除任意文件。
func (s *uploadService) DeleteBlogImage(ctx context.Context, name string) result.Result {
	_ = ctx

	objectName := normalizeBlogObjectName(name, s.ossCfg)
	if objectName == "" {
		return result.Fail("非法图片路径")
	}

	client, err := oss.New(
		s.ossCfg.Endpoint,
		s.ossCfg.AccessKeyID,
		s.ossCfg.AccessKeySecret,
	)
	if err != nil {
		return result.Fail("连接 OSS 失败")
	}

	bucket, err := client.Bucket(s.ossCfg.Bucket)
	if err != nil {
		return result.Fail("获取 OSS Bucket 失败")
	}

	if err := bucket.DeleteObject(objectName); err != nil {
		return result.Fail("删除 OSS 图片失败")
	}

	return result.OK()
}

func normalizeBlogObjectName(name string, cfg config.OSSConfig) string {
	objectName := strings.TrimSpace(name)
	if objectName == "" {
		return ""
	}

	publicPrefix := strings.TrimRight(cfg.PublicHost, "/") + "/"
	objectName = strings.TrimPrefix(objectName, publicPrefix)
	objectName = strings.TrimPrefix(objectName, "/imgs/")
	objectName = strings.TrimPrefix(objectName, "imgs/")
	objectName = strings.TrimPrefix(objectName, "/")

	if strings.Contains(objectName, "..") || strings.Contains(objectName, "\\") {
		return ""
	}

	dir := strings.Trim(cfg.Dir, "/")
	if dir != "" && objectName != dir && !strings.HasPrefix(objectName, dir+"/") {
		return ""
	}

	return objectName
}
