package dto

type LoginFormDTO struct {
	Phone    string `json:"phone" form:"phone"`
	Code     string `json:"code" form:"code"`
	Password string `json:"password" form:"password"`
}

type UserDTO struct {
	ID       int64  `json:"id"`
	NickName string `json:"nickName"`
	Icon     string `json:"icon"`
}
