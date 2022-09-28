package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceHttpRule struct {
	ID             int64  `json:"id" gorm:"primary_key"`
	ServiceID      int64  `json:"service_id" gorm:"column:service_id" description:"服务id"`
	RuleType       int    `json:"rule_type" gorm:"column:rule_type" description:"匹配类型 domain=域名, url_prefix=url前缀"`
	Rule           string `json:"rule" gorm:"column:rule" description:"type=domain表示域名，type=url_prefix时表示url前缀"`
	NeedHttps      int    `json:"need_https" gorm:"column:need_https" description:"type=支持https 1=支持"`
	NeedWebsocket  int    `json:"need_websocket" gorm:"column:need_websocket" description:"启用websocket 1=启用"`
	NeedStripUri   int    `json:"need_strip_uri" gorm:"column:need_strip_uri" description:"启用strip_uri 1=启用"`
	UrlRewrite     string `json:"url_rewrite" gorm:"column:url_rewrite" description:"url重写功能，每行一个	"`
	HeaderTransfor string `json:"header_transfor" gorm:"column:header_transfor" description:"header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue	"`
}

func (t *ServiceHttpRule) TableName() string {
	return "gateway_service_http_rule"
}

func (t *ServiceHttpRule) Find(c *gin.Context, tx *gorm.DB, search *ServiceHttpRule) (*ServiceHttpRule, error) {
	out := &ServiceHttpRule{}
	err := tx.WithContext(c).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *ServiceHttpRule) Exist(c *gin.Context, tx *gorm.DB, search *ServiceHttpRule) (*ServiceHttpRule, error) {
	out := &ServiceHttpRule{}
	err := tx.WithContext(c).Where(search).First(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *ServiceHttpRule) Save(c *gin.Context, tx *gorm.DB) error {

	return tx.WithContext(c).Save(t).Error

}
