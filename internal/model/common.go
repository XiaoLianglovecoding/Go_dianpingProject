package model

import "time"

// TimeFields 抽出所有表里常见的创建时间和更新时间字段。
//
// gorm:"column:create_time" 表示这个 Go 字段对应数据库列 create_time。
// json:"createTime" 表示返回给前端时字段名叫 createTime。
type TimeFields struct {
	// 加上 autoCreateTime，插入数据时 GORM 自动填入当前时间
	CreateTime time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime"`

	// 加上 autoUpdateTime，更新/插入数据时 GORM 自动填入当前时间
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time;autoUpdateTime"`
}
