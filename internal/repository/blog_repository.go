package repository

import (
	"context"

	"hmdp-go/internal/model"

	"gorm.io/gorm"
)

type BlogRepository interface {
	// SaveBlog 保存一篇新博客。
	SaveBlog(ctx context.Context, blog *model.Blog) error
	// FindBlogByID 根据博客 id 查询一篇博客。
	FindBlogByID(ctx context.Context, id int64) (*model.Blog, error)
	// FindBlogsByHot 查询热门博客列表。
	FindBlogsByHot(ctx context.Context, current int) ([]model.Blog, error)
	// FindBlogsByUserID 查询某个用户发布的博客。
	FindBlogsByUserID(ctx context.Context, userID int64, current int) ([]model.Blog, error)
	// UpdateBlogLiked 修改博客点赞数，delta 可以是 +1 或 -1。
	UpdateBlogLiked(ctx context.Context, id int64, delta int) error
}

type blogRepository struct {
	// db 是 GORM 数据库连接对象。
	db *gorm.DB
}

// NewBlogRepository 创建博客 Repository。
func NewBlogRepository(db *gorm.DB) BlogRepository {
	return &blogRepository{db: db}
}

// SaveBlog 负责向 tb_blog 插入一条博客记录。
func (r *blogRepository) SaveBlog(ctx context.Context, blog *model.Blog) error {
	// TODO: Insert a new blog into tb_blog.
	return nil
}

// FindBlogByID 负责按主键查询博客详情。
func (r *blogRepository) FindBlogByID(ctx context.Context, id int64) (*model.Blog, error) {
	var blog model.Blog
	err := r.db.WithContext(ctx).
		First(&blog, id).Error
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

// FindBlogsByHot 查询热门博客。
//
// 当前实现先按 liked 点赞数倒序排序，每页 10 条。
// 对应 SQL 大概是:
// SELECT * FROM tb_blog ORDER BY liked DESC LIMIT 10 OFFSET ?;
func (r *blogRepository) FindBlogsByHot(ctx context.Context, current int) ([]model.Blog, error) {
	var blogs []model.Blog

	// current 是页码，前端没传时默认是 1。这里兜底避免非法页码。
	if current < 1 {
		current = 1
	}

	// pageSize 是每页数量；offset 是从第几条开始查。
	const pageSize = 10
	offset := (current - 1) * pageSize

	// Order("liked DESC") 表示点赞数从高到低；
	// Offset + Limit 是 MySQL 分页；
	// Find(&blogs) 查询多条博客记录。
	err := r.db.WithContext(ctx).Order("liked DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&blogs).Error

	if err != nil {
		return nil, err
	}
	return blogs, nil
}

// FindBlogsByUserID 后面用于个人主页，查询某个用户的博客列表。
func (r *blogRepository) FindBlogsByUserID(ctx context.Context, userID int64, current int) ([]model.Blog, error) {
	var blogs []model.Blog
	if current < 1 {
		current = 1
	}
	const pageSize = 10
	offset := (current - 1) * pageSize
	err := r.db.WithContext(ctx).Order("create_time DESC").
		Where("user_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Find(&blogs).Error
	if err != nil {
		return nil, err
	}
	return blogs, nil
}

// UpdateBlogLiked 后面用于点赞/取消点赞时更新 tb_blog.liked。
func (r *blogRepository) UpdateBlogLiked(ctx context.Context, id int64, delta int) error {
	// TODO: Increase or decrease liked count atomically.
	return nil
}
