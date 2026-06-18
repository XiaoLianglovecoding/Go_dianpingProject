package model

type Shop struct {
	ID        int64   `json:"id" gorm:"column:id;primaryKey"`
	Name      string  `json:"name" gorm:"column:name"`
	TypeID    int64   `json:"typeId" gorm:"column:type_id"`
	Images    string  `json:"images" gorm:"column:images"`
	Area      string  `json:"area" gorm:"column:area"`
	Address   string  `json:"address" gorm:"column:address"`
	X         float64 `json:"x" gorm:"column:x"`
	Y         float64 `json:"y" gorm:"column:y"`
	AvgPrice  int64   `json:"avgPrice" gorm:"column:avg_price"`
	Sold      int     `json:"sold" gorm:"column:sold"`
	Comments  int     `json:"comments" gorm:"column:comments"`
	Score     int     `json:"score" gorm:"column:score"`
	OpenHours string  `json:"openHours" gorm:"column:open_hours"`
	Distance  float64 `json:"distance,omitempty" gorm:"-"`
	TimeFields
}

func (Shop) TableName() string {
	return "tb_shop"
}

type ShopType struct {
	ID   int64  `json:"id" gorm:"column:id;primaryKey"`
	Name string `json:"name" gorm:"column:name"`
	Icon string `json:"icon" gorm:"column:icon"`
	Sort int    `json:"sort" gorm:"column:sort"`
	TimeFields
}

func (ShopType) TableName() string {
	return "tb_shop_type"
}
