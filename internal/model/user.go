package model

import "time"

// User 对应数据库表 tb_user。
//
// 这个结构体主要表示登录账号的基础信息。
type User struct {
	ID       int64  `json:"id" gorm:"column:id;primaryKey"`
	Phone    string `json:"phone" gorm:"column:phone"`        // 手机号，对应 tb_user.phone。
	Password string `json:"password" gorm:"column:password"`  // 密码字段，返回给前端时通常不要使用。
	NickName string `json:"nickName" gorm:"column:nick_name"` // 昵称，数据库列名是 nick_name。
	Icon     string `json:"icon" gorm:"column:icon"`          // 头像路径。
	TimeFields
}

// TableName 告诉 GORM：User 这个结构体对应数据库里的 tb_user 表。
func (User) TableName() string {
	return "tb_user"
}

// UserInfo 对应数据库表 tb_user_info，保存用户扩展资料。
//
// 它和 tb_user 的关系通常是一对一，主键是 user_id。
type UserInfo struct {
	UserID    int64      `json:"userId" gorm:"column:user_id;primaryKey"`
	City      string     `json:"city" gorm:"column:city"`
	Introduce string     `json:"introduce" gorm:"column:introduce"`
	Fans      int        `json:"fans" gorm:"column:fans"`
	Followee  int        `json:"followee" gorm:"column:followee"`
	Gender    int        `json:"gender" gorm:"column:gender"`
	Birthday  *time.Time `json:"birthday" gorm:"column:birthday"`
	Credits   int        `json:"credits" gorm:"column:credits"`
	Level     int        `json:"level" gorm:"column:level"`
	TimeFields
}

// TableName 告诉 GORM：UserInfo 对应 tb_user_info 表。
func (UserInfo) TableName() string {
	return "tb_user_info"
}
