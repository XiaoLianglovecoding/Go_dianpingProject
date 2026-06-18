package model

type Blog struct {
	ID       int64   `json:"id" gorm:"column:id;primaryKey"`
	ShopID   int64   `json:"shopId" gorm:"column:shop_id"`
	UserID   int64   `json:"userId" gorm:"column:user_id"`
	Title    string  `json:"title" gorm:"column:title"`
	Images   string  `json:"images" gorm:"column:images"`
	Content  string  `json:"content" gorm:"column:content"`
	Liked    int     `json:"liked" gorm:"column:liked"`
	Comments int     `json:"comments" gorm:"column:comments"`
	Name     string  `json:"name,omitempty" gorm:"-"`
	Icon     string  `json:"icon,omitempty" gorm:"-"`
	IsLike   bool    `json:"isLike,omitempty" gorm:"-"`
	Distance float64 `json:"distance,omitempty" gorm:"-"`
	TimeFields
}

func (Blog) TableName() string {
	return "tb_blog"
}

type BlogComments struct {
	ID       int64  `json:"id" gorm:"column:id;primaryKey"`
	UserID   int64  `json:"userId" gorm:"column:user_id"`
	BlogID   int64  `json:"blogId" gorm:"column:blog_id"`
	ParentID int64  `json:"parentId" gorm:"column:parent_id"`
	AnswerID int64  `json:"answerId" gorm:"column:answer_id"`
	Content  string `json:"content" gorm:"column:content"`
	Liked    int    `json:"liked" gorm:"column:liked"`
	Status   int    `json:"status" gorm:"column:status"`
	TimeFields
}

func (BlogComments) TableName() string {
	return "tb_blog_comments"
}
