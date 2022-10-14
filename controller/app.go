package controller

import (
	"ggateway/dto"
	"ggateway/middleware"
	"github.com/gin-gonic/gin"
)

type AppController struct {
}

func (app *AppController) AppList(c *gin.Context) {
	params := &dto.APPListInput{}
	if err := params.GetValidParams(c); err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

}
