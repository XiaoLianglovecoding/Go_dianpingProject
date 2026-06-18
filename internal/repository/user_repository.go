package repository

import (
	"context"

	"hmdp-go/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindUserByPhone(ctx context.Context, phone string) (*model.User, error)
	FindUserByID(ctx context.Context, id int64) (*model.User, error)
	FindUserInfoByID(ctx context.Context, id int64) (*model.UserInfo, error)
	CreateUser(ctx context.Context, user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindUserByPhone(ctx context.Context, phone string) (*model.User, error) {
	// TODO: Query tb_user by phone with GORM.
	return nil, nil
}

func (r *userRepository) FindUserByID(ctx context.Context, id int64) (*model.User, error) {
	// TODO: Query tb_user by id with GORM.
	return nil, nil
}

func (r *userRepository) FindUserInfoByID(ctx context.Context, id int64) (*model.UserInfo, error) {
	// TODO: Query tb_user_info by user_id with GORM.
	return nil, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	// TODO: Insert a new user into tb_user.
	return nil
}
