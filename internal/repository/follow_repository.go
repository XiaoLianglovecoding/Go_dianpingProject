package repository

import (
	"context"

	"hmdp-go/internal/model"

	"gorm.io/gorm"
)

type FollowRepository interface {
	FindFollow(ctx context.Context, userID int64, followUserID int64) (*model.Follow, error)
	SaveFollow(ctx context.Context, follow *model.Follow) error
	DeleteFollow(ctx context.Context, userID int64, followUserID int64) error
	FindCommonFollows(ctx context.Context, userID int64, otherUserID int64) ([]model.User, error)
}

type followRepository struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) FollowRepository {
	return &followRepository{db: db}
}

func (r *followRepository) FindFollow(ctx context.Context, userID int64, followUserID int64) (*model.Follow, error) {
	// TODO: Query tb_follow for a single follow relation.
	return nil, nil
}

func (r *followRepository) SaveFollow(ctx context.Context, follow *model.Follow) error {
	// TODO: Insert follow relation.
	return nil
}

func (r *followRepository) DeleteFollow(ctx context.Context, userID int64, followUserID int64) error {
	// TODO: Delete follow relation.
	return nil
}

func (r *followRepository) FindCommonFollows(ctx context.Context, userID int64, otherUserID int64) ([]model.User, error) {
	// TODO: Query common follow users, or use Redis set intersection later.
	return nil, nil
}
