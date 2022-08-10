package dto

type AdminLoginInput struct {
	UserName string `json:"username" form:"username" comment:"用户名" example:"admin" validate:"required"`
	Password string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"`
}
