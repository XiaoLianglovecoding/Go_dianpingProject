package model

import "time"

type User struct {
	ID       int64  `json:"id" gorm:"column:id;primaryKey"`
	Phone    string `json:"phone" gorm:"column:phone"`
	Password string `json:"password" gorm:"column:password"`
	NickName string `json:"nickName" gorm:"column:nick_name"`
	Icon     string `json:"icon" gorm:"column:icon"`
	TimeFields
}

func (User) TableName() string {
	return "tb_user"
}

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

func (UserInfo) TableName() string {
	return "tb_user_info"
}
