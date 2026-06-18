package model

type Follow struct {
	ID           int64 `json:"id" gorm:"column:id;primaryKey"`
	UserID       int64 `json:"userId" gorm:"column:user_id"`
	FollowUserID int64 `json:"followUserId" gorm:"column:follow_user_id"`
	TimeFields
}

func (Follow) TableName() string {
	return "tb_follow"
}
