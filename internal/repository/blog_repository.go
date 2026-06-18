package repository

import (
	"context"

	"hmdp-go/internal/model"

	"gorm.io/gorm"
)

type BlogRepository interface {
	SaveBlog(ctx context.Context, blog *model.Blog) error
	FindBlogByID(ctx context.Context, id int64) (*model.Blog, error)
	FindBlogsByHot(ctx context.Context, current int) ([]model.Blog, error)
	FindBlogsByUserID(ctx context.Context, userID int64, current int) ([]model.Blog, error)
	UpdateBlogLiked(ctx context.Context, id int64, delta int) error
}

type blogRepository struct {
	db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) BlogRepository {
	return &blogRepository{db: db}
}

func (r *blogRepository) SaveBlog(ctx context.Context, blog *model.Blog) error {
	// TODO: Insert a new blog into tb_blog.
	return nil
}

func (r *blogRepository) FindBlogByID(ctx context.Context, id int64) (*model.Blog, error) {
	// TODO: Query tb_blog by id.
	return nil, nil
}

func (r *blogRepository) FindBlogsByHot(ctx context.Context, current int) ([]model.Blog, error) {
	// TODO: Query hot blogs ordered by liked desc.
	return nil, nil
}

func (r *blogRepository) FindBlogsByUserID(ctx context.Context, userID int64, current int) ([]model.Blog, error) {
	// TODO: Query blogs by user_id.
	return nil, nil
}

func (r *blogRepository) UpdateBlogLiked(ctx context.Context, id int64, delta int) error {
	// TODO: Increase or decrease liked count atomically.
	return nil
}
