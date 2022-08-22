package dto

import (
	"ggateway/public"
	"github.com/gin-gonic/gin"
)

type ServiceListInput struct {
	Info       string `json:"info" form:"info" comment:"关键词" example:"" validate:""`                       // 关键词
	PageNumber int    `json:"page_number" form:"page_number" comment:"页码" example:"1" validate:"required"` // 页码
	PageSize   int    `json:"page_size" form:"page_size" comment:"每页数量" example:"20" validate:"required"`  // 每页数量
}

type ServiceListOutputItem struct {
	ID          int64  `json:"id" form:"id"`                     // id
	ServiceName string `json:"service_name" form:"service_name"` // 服务名称
	ServiceDesc string `json:"service_desc" form:"service_desc"` // 服务描述
	LoadType    int    `json:"load_type" form:"load_type"`       // 负载类型
	ServiceAddr string `json:"service_addr" form:"service_addr"` // 服务地址
	Qps         int64  `json:"qps" form:"qps"`                   // 每秒请求数量
	Qpd         int64  `json:"qpd" form:"qpd"`                   // 每日请求数量
	TotalNode   int    `json:"total_node" form:"total_node"`     // 节点总数
}

type ServiceListOutput struct {
	Total int64                   `json:"total" form:"total" comment:"总数" example:"" validate:""`         // 总数
	List  []ServiceListOutputItem `json:"list" form:"list" comment:"服务列表" example:"" validate:"required"` // 服务列表
}

func (param *ServiceListInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}
