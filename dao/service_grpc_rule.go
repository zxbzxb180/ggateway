package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceGrpcRule struct {
	ID             int64  `json:"id" gorm:"primary_key"`
	ServiceID      int64  `json:"service_id" gorm:"column:service_id" description:"服务id	"`
	Port           int    `json:"port" gorm:"column:port" description:"端口	"`
	HeaderTransfor string `json:"header_transfor" gorm:"column:header_transfor" description:"header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue"`
}

func (t *ServiceGrpcRule) TableName() string {
	return "gateway_service_grpc_rule"
}

func (t *ServiceGrpcRule) Find(c *gin.Context, tx *gorm.DB, search *ServiceGrpcRule) (*ServiceGrpcRule, error) {
	out := &ServiceGrpcRule{}
	err := tx.WithContext(c).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *ServiceGrpcRule) Exist(c *gin.Context, tx *gorm.DB, search *ServiceGrpcRule) (*ServiceGrpcRule, error) {
	out := &ServiceGrpcRule{}
	err := tx.WithContext(c).Where(search).First(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *ServiceGrpcRule) Save(c *gin.Context, tx *gorm.DB) error {

	return tx.WithContext(c).Save(t).Error

}
