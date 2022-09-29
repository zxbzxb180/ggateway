package controller

import (
	"fmt"
	"ggateway/dao"
	"ggateway/dto"
	"ggateway/middleware"
	"ggateway/public"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/zxbzxb180/golang_common/lib"
	"gorm.io/gorm"
	"strings"
	"time"
)

type ServiceController struct{}

func ServiceRegister(group *gin.RouterGroup) {
	service := &ServiceController{}
	group.GET("/service_list", service.ServiceList)
	group.GET("/service_delete", service.ServiceDelete)
	group.GET("/service_detail", service.ServiceDetail)
	group.GET("/service_statistics", service.ServiceStatistics)
	group.POST("/service_add_http", service.ServiceAddHttp)
	group.POST("/service_update_http", service.ServiceUpdateHTTP)
	group.POST("/service_add_tcp", service.ServiceAddTcp)
	group.POST("/service_update_tcp", service.ServiceUpdateTcp)
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
func (serviceController *ServiceController) ServiceList(c *gin.Context) {
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
func (serviceController *ServiceController) ServiceDelete(c *gin.Context) {
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
func (serviceController *ServiceController) ServiceDetail(c *gin.Context) {
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

// ServiceStatistics godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags 服务管理
// @ID /service/service_statistics
// @Accept  json
// @Produce  json
// @Param service_id query string true "服务ID"
// @Success 200 {object} middleware.Response{data=dto.ServiceStatisticsOutput} "success"
// @Router /service/service_statistics [get]
func (serviceController *ServiceController) ServiceStatistics(c *gin.Context) {
	params := &dto.ServiceStatisticsInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 4001, err)
		return
	}

	// 连接orm
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 4002, err)
		return
	}

	// 通过查询服务详情
	serviceInfo := &dao.ServiceInfo{ID: params.ServiceId}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	//serviceDetail, err := serviceInfo.GetServiceDetail(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 4003, err)
		return
	}
	counter, err := public.FlowCounterHandler.GetCounter(public.FlowServicePrefix + serviceInfo.ServiceName)
	if err != nil {
		middleware.ResponseError(c, 4004, err)
		return
	}
	todayList := []int64{}
	currentTime := time.Now()
	for i := 0; i <= currentTime.Hour(); i++ {
		dateTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		hourData, _ := counter.GetHourData(dateTime)
		todayList = append(todayList, hourData)
	}

	yesterdayList := []int64{}
	yesterTime := currentTime.Add(-1 * time.Duration(time.Hour*24))
	for i := 0; i <= 23; i++ {
		dateTime := time.Date(yesterTime.Year(), yesterTime.Month(), yesterTime.Day(), i, 0, 0, 0, lib.TimeLocation)
		hourData, _ := counter.GetHourData(dateTime)
		yesterdayList = append(yesterdayList, hourData)
	}
	middleware.ResponseSuccess(c, &dto.ServiceStatisticsOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	})

}

// ServiceAddHttp godoc
// @Summary 添加HTTP服务
// @Description 添加HTTP服务
// @Tags 服务管理
// @ID /service/service_add_http
// @Accept  json
// @Produce  json
// @Param data body dto.ServiceAddHttpInput true "http detail"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_add_http [post]
func (serviceController *ServiceController) ServiceAddHttp(c *gin.Context) {
	params := &dto.ServiceAddHttpInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 5001, err)
		return
	}
	// 判断IP列表长度与权重列表长度是否一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 5002, errors.New("IP列表与权重列表数量不一致"))
		return
	}

	// 连接orm
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 5003, err)
		return
	}
	// 开始事务
	tx = tx.Begin()

	// 判断服务是否已存在
	serviceInfo := &dao.ServiceInfo{ServiceName: params.ServiceName}
	if _, err = serviceInfo.Exist(c, tx, serviceInfo); err == nil {
		tx.Rollback()
		middleware.ResponseError(c, 5004, errors.New("服务名被占用，请重新输入"))
		return
	}

	// 判断服务前缀或域名是否存在
	httpUrl := &dao.ServiceHttpRule{RuleType: params.RuleType, Rule: params.Rule}
	if _, err := httpUrl.Exist(c, tx, httpUrl); err == nil {
		tx.Rollback()
		middleware.ResponseError(c, 5005, errors.New("服务接入前缀或域名已存在"))
		return
	}

	// 新增一条ServiceInfo数据
	serviceModel := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
		LoadType:    public.LoadTypeHTTP,
	}
	if err := serviceModel.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 5006, err)
		return
	}
	// 根据前面新增的ServiceInfo获取ID，新增一条httpRule数据
	httpRule := &dao.ServiceHttpRule{
		ServiceID:      serviceModel.ID,
		RuleType:       params.RuleType,
		Rule:           params.Rule,
		NeedHttps:      params.NeedHttps,
		NeedStripUri:   params.NeedStripUri,
		NeedWebsocket:  params.NeedWebsocket,
		UrlRewrite:     params.UrlRewrite,
		HeaderTransfor: params.HeaderTransfor,
	}
	if err := httpRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 5007, err)
		return
	}

	// 根据前面新增的ServiceInfo获取ID，新增一条ServiceAccessControl数据
	accessControl := &dao.ServiceAccessControl{
		ServiceID:         serviceModel.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		ClientIPFlowLimit: params.ClientipFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 5008, err)
		return
	}

	// 根据前面新增的ServiceInfo获取ID，新增一条ServiceLoadBalance数据
	loadBalance := &dao.ServiceLoadBalance{
		ServiceID:              serviceModel.ID,
		RoundType:              params.RoundType,
		IpList:                 params.IpList,
		WeightList:             params.WeightList,
		UpstreamConnectTimeout: params.UpstreamConnectTimeout,
		UpstreamHeaderTimeout:  params.UpstreamHeaderTimeout,
		UpstreamIdleTimeout:    params.UpstreamIdleTimeout,
		UpstreamMaxIdle:        params.UpstreamMaxIdle,
	}
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 5009, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "success")
}

// ServiceUpdateHTTP godoc
// @Summary 修改HTTP服务
// @Description 修改HTTP服务
// @Tags 服务管理
// @ID /service/service_update_http
// @Accept  json
// @Produce  json
// @Param data body dto.ServiceAddHttpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_update_http [post]
func (serviceController *ServiceController) ServiceUpdateHTTP(c *gin.Context) {
	params := &dto.ServiceAddHttpInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 6001, err)
		return
	}

	// 判断IP列表长度与权重列表长度是否一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 6002, errors.New("IP列表与权重列表数量不一致"))
		return
	}

	// 连接orm
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 6003, err)
		return
	}
	// 开始事务
	tx = tx.Begin()

	// 判断服务基本信息数据是否已存在
	serviceInfo := &dao.ServiceInfo{ServiceName: params.ServiceName}
	if serviceInfo, err = serviceInfo.Exist(c, tx, serviceInfo); err != nil && err == gorm.ErrRecordNotFound {
		tx.Rollback()
		middleware.ResponseError(c, 6004, errors.New("服务不存在"))
		return
	}

	// 判断服务详细信息数据是否已存在
	serviceDetail, err := serviceInfo.GetServiceDetail(c, tx, serviceInfo)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 6005, errors.New("服务不存在"))
		return
	}

	// 修改服务详情
	info := serviceDetail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 6006, err)
		return
	}

	// 修改http规则
	httpRule := serviceDetail.HttpRule
	httpRule.Rule = params.Rule
	httpRule.NeedHttps = params.NeedHttps
	httpRule.NeedStripUri = params.NeedStripUri
	httpRule.NeedWebsocket = params.NeedWebsocket
	httpRule.UrlRewrite = params.UrlRewrite
	httpRule.HeaderTransfor = params.HeaderTransfor
	if err := httpRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 6007, err)
		return
	}

	// 修改访问控制
	accessControl := serviceDetail.AccessControl
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.ClientIPFlowLimit = params.ClientipFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 6008, err)
		return
	}

	loadbalance := serviceDetail.LoadBalance
	loadbalance.RoundType = params.RoundType
	loadbalance.IpList = params.IpList
	loadbalance.WeightList = params.WeightList
	loadbalance.UpstreamConnectTimeout = params.UpstreamConnectTimeout
	loadbalance.UpstreamHeaderTimeout = params.UpstreamHeaderTimeout
	loadbalance.UpstreamIdleTimeout = params.UpstreamIdleTimeout
	loadbalance.UpstreamMaxIdle = params.UpstreamMaxIdle
	if err := loadbalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 6009, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "success")
}

// ServiceAddTcp godoc
// @Summary tcp服务添加
// @Description tcp服务添加
// @Tags 服务管理
// @ID /service/service_add_tcp
// @Accept  json
// @Produce  json
// @Param data body dto.ServiceAddTcpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_add_tcp [post]
func (serviceController *ServiceController) ServiceAddTcp(c *gin.Context) {
	params := &dto.ServiceAddTcpInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 7001, err)
		return
	}
	//判断ip与权重数量是否一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 7002, errors.New("ip列表与权重设置不匹配"))
		return
	}
	// 连接orm
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 7003, err)
		return
	}

	info := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
	}

	if _, err := info.Exist(c, tx, info); err == nil {
		middleware.ResponseError(c, 7004, errors.New("服务名被占用，请重新输入"))
		return
	}
	//验证端口是否被占用?
	tcpRuleSearch := &dao.ServiceTcpRule{
		Port: params.Port,
	}
	if _, err := tcpRuleSearch.Exist(c, tx, tcpRuleSearch); err == nil {
		middleware.ResponseError(c, 7005, errors.New("服务端口被占用，请重新输入"))
		return
	}

	//开始事务
	tx = tx.Begin()
	newInfo := &dao.ServiceInfo{
		LoadType:    public.LoadTypeTCP,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := newInfo.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}
	loadBalance := &dao.ServiceLoadBalance{
		ServiceID:  newInfo.ID,
		RoundType:  params.RoundType,
		IpList:     params.IpList,
		WeightList: params.WeightList,
		ForbidList: params.ForbidList,
	}
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}

	httpRule := &dao.ServiceTcpRule{
		ServiceID: newInfo.ID,
		Port:      params.Port,
	}
	if err := httpRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2008, err)
		return
	}

	accessControl := &dao.ServiceAccessControl{
		ServiceID:         newInfo.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		WhiteHostName:     params.WhiteHostName,
		ClientIPFlowLimit: params.ClientIPFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2009, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "success")
	return
}

// ServiceUpdateTcp godoc
// @Summary tcp服务更新
// @Description tcp服务更新
// @Tags 服务管理
// @ID /service/service_update_tcp
// @Accept  json
// @Produce  json
// @Param data body dto.ServiceUpdateTcpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_update_tcp [post]
func (serviceController ServiceController) ServiceUpdateTcp(c *gin.Context) {
	params := &dto.ServiceUpdateTcpInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 8001, err)
		return
	}
	//判断ip与权重数量是否一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(c, 8002, errors.New("ip列表与权重设置不匹配"))
		return
	}
	// 连接orm
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 8003, err)
		return
	}
	service := &dao.ServiceInfo{
		ID: params.ID,
	}
	detail, err := service.GetServiceDetail(c, tx, service)
	if err != nil {
		middleware.ResponseError(c, 8004, err)
		return
	}

	// 开始事务
	tx = tx.Begin()
	info := detail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 8005, err)
		return
	}
	loadBalance := &dao.ServiceLoadBalance{}
	if detail.LoadBalance != nil {
		loadBalance = detail.LoadBalance
	}
	loadBalance.ServiceID = info.ID
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IpList
	loadBalance.WeightList = params.WeightList
	loadBalance.ForbidList = params.ForbidList
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 8006, err)
		return
	}

	tcpRule := &dao.ServiceTcpRule{}
	if detail.TcpRule != nil {
		tcpRule = detail.TcpRule
	}
	tcpRule.ServiceID = info.ID
	tcpRule.Port = params.Port
	if err := tcpRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 8007, err)
		return
	}

	accessControl := &dao.ServiceAccessControl{}
	if detail.AccessControl != nil {
		accessControl = detail.AccessControl
	}
	accessControl.ServiceID = info.ID
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.WhiteHostName = params.WhiteHostName
	accessControl.ClientIPFlowLimit = params.ClientIPFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 8008, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(c, "success")
	return
}
