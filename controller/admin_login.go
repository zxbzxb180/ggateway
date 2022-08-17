package controller

import (
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/zxbzxb180/ggateway/dao"
	"github.com/zxbzxb180/ggateway/dto"
	"github.com/zxbzxb180/ggateway/middleware"
)

type AdminLoginController struct{}

func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("/login", adminLogin.AdminLogin)
}

// AdminLogin godoc
// @Summary 管理员登录
// @Description 管理员登录
// @Tags 管理员接口
// @ID /admin_login/login
// @Accept json
// @Produce	json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminLoginOutput} "success"
// @Router /admin_login/login [post]
func (adminlogin *AdminLoginController) AdminLogin(c *gin.Context) {
	params := &dto.AdminLoginInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 1001, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 1002, err)
	}

	admin := &dao.Admin{}
	admin, err = admin.LoginCheck(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 1003, err)
	}

	out := &dto.AdminLoginOutput{Token: params.UserName}
	middleware.ResponseSuccess(c, out)
}
