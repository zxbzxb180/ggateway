package dao

import (
	"ggateway/dto"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type ServiceInfo struct {
	ID          int64     `json:"id" gorm:"column:id" form:"id"`                               // id
	ServiceName string    `json:"service_name" gorm:"column:service_name" form:"service_name"` // 服务名称
	ServiceDesc string    `json:"service_desc" gorm:"column:service_desc" form:"service_desc"` // 服务描述
	LoadType    int       `json:"load_type" gorm:"column:load_type" form:"load_type"`          // 负载类型
	CreatedAt   time.Time `json:"create_at" gorm:"column:create_at" form:"create_at"`          // 创建时间
	UpdatedAt   time.Time `json:"update_at" gorm:"column:update_at" form:"update_at"`          // 更新时间
	IsValid     int8      `json:"is_valid"  gorm:"column:is_valid" form:"is_valid"`            // 是否有效
}

func (t *ServiceInfo) TableName() string {
	return "gateway_service_info"
}

func (t *ServiceInfo) GetServiceDetail(c *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceDetail, error) {
	httpRule := &ServiceHttpRule{ServiceID: search.ID}
	httpRule, err := httpRule.Find(c, tx, httpRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	tcpRule := &ServiceTcpRule{ServiceID: search.ID}
	tcpRule, err = tcpRule.Find(c, tx, tcpRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	grpcRule := &ServiceGrpcRule{ServiceID: search.ID}
	grpcRule, err = grpcRule.Find(c, tx, grpcRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	loadBalance := &ServiceLoadBalance{ServiceID: search.ID}
	loadBalance, err = loadBalance.Find(c, tx, loadBalance)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	accessControl := &ServiceAccessControl{ServiceID: search.ID}
	accessControl, err = accessControl.Find(c, tx, accessControl)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	detail := &ServiceDetail{
		Info:          search,
		HttpRule:      httpRule,
		TcpRule:       tcpRule,
		GrpcRule:      grpcRule,
		LoadBalance:   loadBalance,
		AccessControl: accessControl,
	}
	return detail, nil
}

func (t *ServiceInfo) GetPageList(c *gin.Context, tx *gorm.DB, param *dto.ServiceListInput) ([]ServiceInfo, int64, error) {
	total := int64(0)
	list := []ServiceInfo{}
	offset := (param.PageNumber - 1) * param.PageSize
	query := tx.WithContext(c).Table(t.TableName()).Where("is_valid=1")
	if param.Info != "" {
		query = query.Where("(service_name like ? or service_desc like ?)", "%"+param.Info+"%", "%"+param.Info+"%")
	}
	if err := query.Limit(param.PageSize).Offset(offset).Order("id").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	query.Limit(param.PageSize).Offset(offset).Count(&total)
	return list, total, nil
}

func (t *ServiceInfo) Find(c *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceInfo, error) {
	out := &ServiceInfo{}
	err := tx.WithContext(c).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *ServiceInfo) Save(c *gin.Context, tx *gorm.DB) error {

	return tx.WithContext(c).Save(t).Error

}
