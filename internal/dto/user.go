package dto

// LoginFormDTO 对应前端登录请求体。
//
// DTO 的意思是 Data Transfer Object：只负责网络传输，不一定和数据库表完全一样。
type LoginFormDTO struct {
	Phone    string `json:"phone" form:"phone"` // 手机号。
	Code     string `json:"code" form:"code"`   // 验证码登录时使用。
	Password string `json:"password" form:"password"`
}

// UserDTO 是返回给前端看的用户信息。
//
// 注意这里没有 Password、Phone 等敏感字段，避免把不该暴露的数据返回出去。
type UserDTO struct {
	ID       int64  `json:"id"`
	NickName string `json:"nickName"`
	Icon     string `json:"icon"`
}

// UserInfoDTO 返回给前端的用户扩展信息
type UserInfoDTO struct {
	City      string `json:"city"`
	Introduce string `json:"introduce"`
	Fans      int    `json:"fans"`
	Followee  int    `json:"followee"`
	Gender    int    `json:"gender"`
	Birthday  string `json:"birthday"` // 建议日期转成格式化字符串给前端，比如 "1999-01-01"
	Credits   int    `json:"credits"`
	Level     int    `json:"level"`
}
