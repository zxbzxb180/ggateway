package controller

import (
	"fmt"
	"ggateway/dao"
	"ggateway/dto"
	"ggateway/middleware"
	"ggateway/public"
	"github.com/gin-gonic/gin"
	"github.com/zxbzxb180/golang_common/lib"
)

type ServiceController struct{}

func ServiceRegister(group *gin.RouterGroup) {
	service := &ServiceController{}
	group.GET("/service_list", service.ServiceList)
	group.GET("/service_delete", service.ServiceDelete)
	group.GET("/service_detail", service.ServiceDetail)
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
func (service *ServiceController) ServiceList(c *gin.Context) {
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

	// 分页读取
	serviceInfo := &dao.ServiceInfo{}
	list, total, err := serviceInfo.GetPageList(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 1003, err)
		return
	}

	// 格式化输出
	outputList := []dto.ServiceListOutputItem{}
	for _, listItem := range list {
		// 1. HTTP后缀接入：clusterIP + clusterPort + path
		// 2. HTTP域名接入：domain
		// 3. tcp、grpc接入：clusterIP + servicePort
		serviceDetail, err := listItem.GetServiceDetail(c, tx, &listItem)
		if err != nil {
			middleware.ResponseError(c, 1004, err)
			return
		}
		serviceAddr := "unknown"
		clusterIP := lib.GetStringConf("base.cluster.cluster_ip")
		clusterPort := lib.GetStringConf("base.cluster.cluster_port")
		clusterSSLPort := lib.GetStringConf("base.cluster.cluster_ssl_port")

		// HTTP 后缀接入
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP && serviceDetail.HttpRule.RuleType == public.HTTPRulePrefixURL && serviceDetail.HttpRule.NeedHttps == 0 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterPort, serviceDetail.HttpRule.Rule)
		}
		// HTTPS 后缀接入
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP && serviceDetail.HttpRule.RuleType == public.HTTPRulePrefixURL && serviceDetail.HttpRule.NeedHttps == 1 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterSSLPort, serviceDetail.HttpRule.Rule)
		}
		// HTTP 域名接入
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP && serviceDetail.HttpRule.RuleType == public.HTTPRuleTypeDomain {
			serviceAddr = serviceDetail.HttpRule.Rule
		}
		// TCP
		if serviceDetail.Info.LoadType == public.LoadTypeTCP {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, serviceDetail.TcpRule.Port)
		}
		// GRPC
		if serviceDetail.Info.LoadType == public.LoadTypeGRPC {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, serviceDetail.GrpcRule.Port)
		}

		outItem := dto.ServiceListOutputItem{
			ID:          listItem.ID,
			ServiceName: listItem.ServiceName,
			ServiceDesc: listItem.ServiceDesc,
			LoadType:    listItem.LoadType,
			ServiceAddr: serviceAddr,
			Qps:         0,
			Qpd:         0,
			TotalNode:   len(serviceDetail.LoadBalance.GetIpListByModel()),
		}
		outputList = append(outputList, outItem)
	}

	out := &dto.ServiceListOutput{
		Total: total,
		List:  outputList,
	}
	middleware.ResponseSuccess(c, out)
}

// ServiceDelete godoc
// @Summary 服务删除
// @Description 服务删除
// @Tags 服务管理
// @ID /service/service_delete
// @Accept json
// @Produce	json
// @Param service_id query int true "服务id"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_delete [get]
func (service *ServiceController) ServiceDelete(c *gin.Context) {
	params := &dto.ServiceDeleteInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	// 查找服务
	serviceInfo := &dao.ServiceInfo{ID: params.ServiceId}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	serviceInfo.IsValid = 0
	if err = serviceInfo.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2004, err)
		return
	}

	middleware.ResponseSuccess(c, "")
}

// ServiceDetail godoc
// @Summary 服务详情
// @Description 服务详情
// @Tags 服务管理
// @ID /service/service_detail
// @Accept  json
// @Produce  json
// @Param service_id query int true "服务ID"
// @Success 200 {object} middleware.Response{data=dao.ServiceDetail} "success"
// @Router /service/service_detail [get]
func (service *ServiceController) ServiceDetail(c *gin.Context) {
	params := &dto.ServiceDetailInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 3001, err)
		return
	}

	// 连接orm
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 3002, err)
		return
	}
	// 读取基本信息
	serviceInfo := &dao.ServiceInfo{ID: params.ServiceId}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 3003, err)
		return
	}
	// 读取详情
	serviceDetail, err := serviceInfo.GetServiceDetail(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 3004, err)
		return
	}
	middleware.ResponseSuccess(c, serviceDetail)
}
