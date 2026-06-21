package model

// Blog 对应数据库表 tb_blog，表示一篇探店笔记。
type Blog struct {
	ID       int64   `json:"id" gorm:"column:id;primaryKey"`
	ShopID   int64   `json:"shopId" gorm:"column:shop_id"`    // 关联店铺 id。
	UserID   int64   `json:"userId" gorm:"column:user_id"`    // 发布博客的用户 id。
	Title    string  `json:"title" gorm:"column:title"`       // 博客标题。
	Images   string  `json:"images" gorm:"column:images"`     // 图片路径，多个图片用逗号分隔。
	Content  string  `json:"content" gorm:"column:content"`   // 正文内容。
	Liked    int     `json:"liked" gorm:"column:liked"`       // 点赞数。
	Comments int     `json:"comments" gorm:"column:comments"` // 评论数。
	Name     string  `json:"name,omitempty" gorm:"-"`         // 作者昵称，不在 tb_blog 表里，查询后手动补充。
	Icon     string  `json:"icon,omitempty" gorm:"-"`         // 作者头像，不在 tb_blog 表里，查询后手动补充。
	IsLike   bool    `json:"isLike" gorm:"-"`                 // 当前登录用户是否点赞，不在 tb_blog 表里。
	Distance float64 `json:"distance,omitempty" gorm:"-"`     // 距离字段，后续做附近商户/博客时可能使用。
	TimeFields
}

// TableName 告诉 GORM：Blog 对应 tb_blog 表。
func (Blog) TableName() string {
	return "tb_blog"
}

// BlogComments 对应数据库表 tb_blog_comments，表示博客评论。
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

// TableName 告诉 GORM：BlogComments 对应 tb_blog_comments 表。
func (BlogComments) TableName() string {
	return "tb_blog_comments"
}
