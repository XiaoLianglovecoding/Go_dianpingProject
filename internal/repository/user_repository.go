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
	FindUsersByIDs(ctx context.Context, ids []int64) ([]model.User, error)
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
	var user model.User
	res := r.db.WithContext(ctx).
		Where("phone = ?", phone).
		Limit(1).Find(&user) // 优化：使用 Find 替代 First，阻止 GORM 自动打印错误日志

	if res.Error != nil {
		return nil, res.Error // 真正的数据库异常（如断网）
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound // 巧妙伪装：手动返回没查到，保持与 Service 层的兼容
	}
	return &user, nil
}

// FindUserByID 根据主键 id 查询用户。
//
// 这个方法现在已经被热门博客接口使用，用来给每篇博客补充作者昵称和头像。
func (r *userRepository) FindUserByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	res := r.db.WithContext(ctx).
		Where("id = ?", id).
		Limit(1).Find(&user) // 同理优化

	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &user, nil
}

// FindUserInfoByID 查询用户扩展资料表 tb_user_info。
func (r *userRepository) FindUserInfoByID(ctx context.Context, id int64) (*model.UserInfo, error) {
	var userinfo model.UserInfo
	res := r.db.WithContext(ctx).
		Where("user_id = ?", id).
		Limit(1).Find(&userinfo) // 同理优化

	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &userinfo, nil
}

// CreateUser 保存新用户。
// 相当于：INSERT INTO `tb_user`
// (`phone`, `nick_name`, `password`, `icon`, `create_time`, `update_time`)
// VALUES
// ('13800138000', 'user_550e8400', ”, ”, '2026-06-20 16:44:46', '2026-06-20 16:44:46');
func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

//	新增：FindUsersByIDs 批量查询用户列表
//
// 相当于 SQL: SELECT * FROM tb_user WHERE id IN (?, ?, ?);
func (r *userRepository) FindUsersByIDs(ctx context.Context, ids []int64) ([]model.User, error) {
	if len(ids) == 0 {
		return []model.User{}, nil
	}
	var users []model.User
	// GORM 的精髓：传一个切片给 ?，它会自动帮你转换成 IN () 的语法
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
