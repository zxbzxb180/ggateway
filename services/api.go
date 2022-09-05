package services

import (
	"ggateway/dao"
	"ggateway/dto"
	"ggateway/middleware"
	"github.com/gin-gonic/gin"
	"github.com/zxbzxb180/golang_common/lib"
)

type ApiService struct {
}

func (as *ApiService) AddUser(c *gin.Context, input *dto.AddUserInput) error {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return err
	}
	user := &dao.User{
		Name:  input.Name,
		Sex:   input.Sex,
		Age:   input.Age,
		Birth: input.Birth,
		Addr:  input.Addr,
	}
	if err := user.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2002, err)
		return err
	}
	return nil
}
