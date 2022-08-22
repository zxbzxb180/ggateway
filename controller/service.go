package controller

import (
	"ggateway/dao"
	"ggateway/dto"
	"ggateway/middleware"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
)

type ServiceController struct{}

func ServiceRegister(group *gin.RouterGroup) {
	service := &ServiceController{}
	group.GET("/service_list", service.ServiceList)
}

// ServiceList godoc
// @Summary 服务列表
// @Description 服务列表
// @Tags 服务管理
// @ID /service/service_list
// @Accept json
// @Produce	json
// @Param info query string false "关键词"
// @Param page_number query int true "页码"
// @Param page_size query int true "每页数量"
// @Success 200 {object} middleware.Response{data=dto.ServiceListOutput} "success"
// @Router /service/service_list [get]
func (service ServiceController) ServiceList(c *gin.Context) {
	params := &dto.ServiceListInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 1001, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 1002, err)
		return
	}

	serviceInfo := &dao.ServiceInfo{}
	list, total, err := serviceInfo.PageList(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 1003, err)
		return
	}

	outputList := []dto.ServiceListOutputItem{}
	for _, listItem := range list {
		outItem := dto.ServiceListOutputItem{
			ID:          listItem.ID,
			ServiceName: listItem.ServiceName,
			ServiceDesc: listItem.ServiceDesc,
			LoadType:    listItem.LoadType,
		}
		outputList = append(outputList, outItem)
	}

	out := &dto.ServiceListOutput{
		Total: total,
		List:  outputList,
	}
	middleware.ResponseSuccess(c, out)
}
