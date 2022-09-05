package controller

import (
	"encoding/json"
	"fmt"
	"ggateway/dao"
	"ggateway/dto"
	"ggateway/middleware"
	"ggateway/public"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/zxbzxb180/golang_common/lib"
)

type AdminController struct{}

func AdminRegister(group *gin.RouterGroup) {
	adminInfo := &AdminController{}
	group.GET("/admin_info", adminInfo.AdminInfo)
	group.POST("/change_pwd", adminInfo.ChangePwd)
}

// AdminInfo godoc
// @Summary 管理员信息
// @Description 管理员信息
// @Tags 管理员接口
// @ID /admin/admin_info
// @Accept json
// @Produce	json
// @Success 200 {object} middleware.Response{data=dto.AdminInfoOutput} "success"
// @Router /admin/admin_info [get]
func (admininfo *AdminController) AdminInfo(c *gin.Context) {
	// 1 读取session key对应json，转换为结构体
	// 2 取出数据封装输出

	session := sessions.Default(c)
	sessionInfo := session.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessionInfo)), adminSessionInfo); err != nil {
		middleware.ResponseError(c, 1001, err)
		return
	}
	out := &dto.AdminInfoOutput{
		ID:           adminSessionInfo.ID,
		UserName:     adminSessionInfo.UserName,
		LoginTime:    adminSessionInfo.LoginTime,
		Avatar:       "",
		Introduction: "I am a super administrator",
		Roles:        []string{"admin"},
	}
	middleware.ResponseSuccess(c, out)
}

// ChangePwd godoc
// @Summary 修改密码
// @Description 修改密码
// @Tags 管理员接口
// @ID /admin/change_pwd
// @Accept json
// @Produce	json
// @Param body body dto.ChangePwdInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin/change_pwd [post]
func (admininfo *AdminController) ChangePwd(c *gin.Context) {

	params := &dto.ChangePwdInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	// 从数据库读取adminInfo
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	// 1 读取session到结构体 sessionInfo
	session := sessions.Default(c)
	sessionInfo := session.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessionInfo)), adminSessionInfo); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	// 2 利用ID查询用户信息
	adminInfo := &dao.Admin{}
	adminInfo, err = adminInfo.Find(c, tx, &dao.Admin{UserName: adminSessionInfo.UserName, ID: adminSessionInfo.ID})
	if err != nil {
		middleware.ResponseError(c, 2004, err)
		return
	}
	// 3 params.password + salt  sha256 =  saltPassword
	saltPassword := public.GenSaltPassword(params.Password, adminInfo.Salt)
	adminInfo.Password = saltPassword

	// 4 save
	if err = adminInfo.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2005, err)
		return
	}

	middleware.ResponseSuccess(c, "")
}
