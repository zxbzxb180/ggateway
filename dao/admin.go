package dao

import (
	"github.com/gin-gonic/gin"
	"github.com/zxbzxb180/ggateway/dto"
	"gorm.io/gorm"
	"time"
)

type Admin struct {
	Id        int       `json:"id" gorm:"primary_key" description:"自增主键"`
	UserName  string    `json:"user_name" gorm:"column:user_name" description:"管理员名称"`
	Salt      int       `json:"salt" gorm:"column:salt" description:"管理员加盐值"`
	Password  int64     `json:"password" gorm:"column:password" description:"密码"`
	UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsValid   int       `json:"is_valid" gorm:"column:is_valid" description:"是否有效"`
}

func (t *Admin) TableName() string {
	return "gateway_admin"
}

func (t *Admin) LoginCheck(c *gin.Context, tx *gorm.DB, param *dto.AdminLoginInput) (*Admin, error) {
	t.Find(c, tx, &Admin{UserName: param.UserName, IsValid: 1})
	return nil, nil
}

func (t *Admin) Find(c *gin.Context, tx *gorm.DB, search *Admin) (*Admin, error) {
	admin := &Admin{}
	err := tx.WithContext(c).Where(search).Find(admin).Error
	if err != nil {
		return nil, err
	}
	return admin, nil
}
