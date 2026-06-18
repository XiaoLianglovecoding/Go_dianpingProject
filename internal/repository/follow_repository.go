package repository

import (
	"context"

	"hmdp-go/internal/model"

	"gorm.io/gorm"
)

type FollowRepository interface {
	// FindFollow 查询当前用户是否关注了目标用户。
	FindFollow(ctx context.Context, userID int64, followUserID int64) (*model.Follow, error)
	// SaveFollow 新增关注关系。
	SaveFollow(ctx context.Context, follow *model.Follow) error
	// DeleteFollow 删除关注关系。
	DeleteFollow(ctx context.Context, userID int64, followUserID int64) error
	// FindCommonFollows 查询共同关注。
	FindCommonFollows(ctx context.Context, userID int64, otherUserID int64) ([]model.User, error)
}

type followRepository struct {
	// db 是 GORM 数据库连接对象。
	db *gorm.DB
}

// NewFollowRepository 创建关注 Repository。
func NewFollowRepository(db *gorm.DB) FollowRepository {
	return &followRepository{db: db}
}

// FindFollow 后面会查询 tb_follow 是否存在一条关注记录。
func (r *followRepository) FindFollow(ctx context.Context, userID int64, followUserID int64) (*model.Follow, error) {
	// TODO: Query tb_follow for a single follow relation.
	return nil, nil
}

// SaveFollow 保存关注关系。
func (r *followRepository) SaveFollow(ctx context.Context, follow *model.Follow) error {
	// TODO: Insert follow relation.
	return nil
}

// DeleteFollow 取消关注。
func (r *followRepository) DeleteFollow(ctx context.Context, userID int64, followUserID int64) error {
	// TODO: Delete follow relation.
	return nil
}

// FindCommonFollows 查询两个用户都关注的人。
func (r *followRepository) FindCommonFollows(ctx context.Context, userID int64, otherUserID int64) ([]model.User, error) {
	// TODO: Query common follow users, or use Redis set intersection later.
	return nil, nil
}
