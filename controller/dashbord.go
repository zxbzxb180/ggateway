package controller

import "github.com/gin-gonic/gin"

type DashboardController struct{}

func DashboardRegister(group *gin.RouterGroup) {
	dashboard := &DashboardController{}
	//group.GET("/panel_group_data", dashboard.)
}
