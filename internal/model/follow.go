package model

import "time"

// Follow 对应数据库表 tb_follow，表示“谁关注了谁”。
type Follow struct {
	ID           int64     `json:"id" gorm:"column:id;primaryKey"`            // 主键。
	UserID       int64     `json:"userId" gorm:"column:user_id"`              // 发起关注的人。
	FollowUserID int64     `json:"followUserId" gorm:"column:follow_user_id"` // 被关注的人。
	CreateTime   time.Time `json:"createTime" gorm:"column:create_time;autoCreateTime"`
}

// TableName 告诉 GORM：Follow 对应 tb_follow 表。
func (Follow) TableName() string {
	return "tb_follow"
}
