package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceTcpRule struct {
	ID        int64 `json:"id" gorm:"primary_key"`
	ServiceID int64 `json:"service_id" gorm:"column:service_id" description:"服务id	"`
	Port      int   `json:"port" gorm:"column:port" description:"端口	"`
}

func (t *ServiceTcpRule) TableName() string {
	return "gateway_service_tcp_rule"
}

func (t *ServiceTcpRule) Find(c *gin.Context, tx *gorm.DB, search *ServiceTcpRule) (*ServiceTcpRule, error) {
	out := &ServiceTcpRule{}
	err := tx.WithContext(c).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *ServiceTcpRule) Save(c *gin.Context, tx *gorm.DB) error {

	return tx.WithContext(c).Save(t).Error

}
