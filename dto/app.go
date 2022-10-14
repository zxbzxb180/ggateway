package dto

import (
	"ggateway/public"
	"github.com/gin-gonic/gin"
)

type APPListInput struct {
	Info     string `json:"info" form:"info" comment:"查找信息" validate:""`
	PageSize int    `json:"page_size" form:"page_size" comment:"页数" validate:"required,min=1,max=999"`
	PageNo   int    `json:"page_no" form:"page_no" comment:"页码" validate:"required,min=1,max=999"`
}

func (params *APPListInput) GetValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, params)
}
