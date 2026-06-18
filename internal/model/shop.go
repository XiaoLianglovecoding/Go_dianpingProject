package model

// Shop 对应数据库表 tb_shop，表示一个店铺。
type Shop struct {
	ID        int64   `json:"id" gorm:"column:id;primaryKey"`
	Name      string  `json:"name" gorm:"column:name"`            // 店铺名称。
	TypeID    int64   `json:"typeId" gorm:"column:type_id"`       // 店铺类型 id，对应 tb_shop_type.id。
	Images    string  `json:"images" gorm:"column:images"`        // 店铺图片，Java 版里多张图片用逗号拼接。
	Area      string  `json:"area" gorm:"column:area"`            // 商圈，例如“大关”。
	Address   string  `json:"address" gorm:"column:address"`      // 详细地址。
	X         float64 `json:"x" gorm:"column:x"`                  // 经度。
	Y         float64 `json:"y" gorm:"column:y"`                  // 纬度。
	AvgPrice  int64   `json:"avgPrice" gorm:"column:avg_price"`   // 人均价格，单位通常是分。
	Sold      int     `json:"sold" gorm:"column:sold"`            // 销量。
	Comments  int     `json:"comments" gorm:"column:comments"`    // 评论数。
	Score     int     `json:"score" gorm:"column:score"`          // 评分，Java 版常用整数保存。
	OpenHours string  `json:"openHours" gorm:"column:open_hours"` // 营业时间。
	Distance  float64 `json:"distance,omitempty" gorm:"-"`        // 距离字段不在数据库中，后面做 GEO 查询时返回给前端。
	TimeFields
}

// TableName 告诉 GORM：Shop 对应 tb_shop 表。
func (Shop) TableName() string {
	return "tb_shop"
}

// ShopType 对应数据库表 tb_shop_type，表示首页顶部的分类图标。
type ShopType struct {
	ID   int64  `json:"id" gorm:"column:id;primaryKey"` // 分类 id。
	Name string `json:"name" gorm:"column:name"`        // 分类名称，例如“美食”。
	Icon string `json:"icon" gorm:"column:icon"`        // 分类图标路径，例如 /imgs/types/ms.png。
	Sort int    `json:"sort" gorm:"column:sort"`        // 排序值，越小越靠前。
	TimeFields
}

// TableName 告诉 GORM：ShopType 对应 tb_shop_type 表。
func (ShopType) TableName() string {
	return "tb_shop_type"
}
