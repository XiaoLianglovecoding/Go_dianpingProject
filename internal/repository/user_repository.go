package repository

import (
	"context"

	"hmdp-go/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	// FindUserByPhone 根据手机号查询用户，登录时使用。
	FindUserByPhone(ctx context.Context, phone string) (*model.User, error)
	// FindUserByID 根据用户 id 查询用户，博客补作者信息时会用。
	FindUserByID(ctx context.Context, id int64) (*model.User, error)
	// FindUserInfoByID 查询用户扩展资料。
	FindUserInfoByID(ctx context.Context, id int64) (*model.UserInfo, error)
	// CreateUser 创建新用户。
	CreateUser(ctx context.Context, user *model.User) error
}

type userRepository struct {
	// db 是 GORM 数据库连接对象。
	db *gorm.DB
}

// NewUserRepository 创建用户 Repository。
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// FindUserByPhone 后面实现登录时会用：通过手机号查 tb_user。
func (r *userRepository) FindUserByPhone(ctx context.Context, phone string) (*model.User, error) {
	// TODO: Query tb_user by phone with GORM.
	return nil, nil
}

// FindUserByID 根据主键 id 查询用户。
//
// 这个方法现在已经被热门博客接口使用，用来给每篇博客补充作者昵称和头像。
func (r *userRepository) FindUserByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindUserInfoByID 查询用户扩展资料表 tb_user_info。
func (r *userRepository) FindUserInfoByID(ctx context.Context, id int64) (*model.UserInfo, error) {
	// TODO: Query tb_user_info by user_id with GORM.
	return nil, nil
}

// CreateUser 保存新用户。
func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	// TODO: Insert a new user into tb_user.
	return nil
}
