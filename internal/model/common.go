package model

import "time"

type TimeFields struct {
	CreateTime time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time"`
}
