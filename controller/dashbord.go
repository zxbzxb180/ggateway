package controller

import (
	"ggateway/dao"
	"ggateway/dto"
	"ggateway/middleware"
	"github.com/gin-gonic/gin"
	"github.com/zxbzxb180/golang_common/lib"
)

type DashboardController struct{}

func DashboardRegister(group *gin.RouterGroup) {
	dashboard := &DashboardController{}
	group.GET("/panel_group_data", dashboard.PannelGroupData)
}

// PanelGroupData godoc
// @Summary 指标统计
// @Description 指标统计
// @Tags 首页大盘
// @ID /dashboard/panel_group_data
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.PanelGroupDataOutput} "success"
// @Router /dashboard/panel_group_data [get]
func (dashboard *DashboardController) PannelGroupData(c *gin.Context) {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	serviceInfo := &dao.ServiceInfo{}
	_, serviceNum, err := serviceInfo.GetPageList(c, tx, &dto.ServiceListInput{PageNumber: 1, PageSize: 1})
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	app := &dao.App{}
}
