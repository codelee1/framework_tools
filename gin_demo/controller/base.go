package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseController struct {
}

func (b *BaseController) Success(g *gin.Context)  {
	g.JSON(http.StatusOK,gin.H{
		"code":http.StatusOK,
		"data":"",
		"msg":"success",
	})
}

func (b *BaseController) SuccessWithData(g *gin.Context, data interface{})  {
	g.JSON(http.StatusOK,gin.H{
		"code":http.StatusOK,
		"data":data,
		"msg":"success",
	})
}

func (b *BaseController) FailureForClient(g *gin.Context,msg string)  {
	g.JSON(http.StatusOK,gin.H{
		"code":http.StatusNotFound,
		"data":"",
		"msg":msg,
	})
}

func (b *BaseController) FailureForServer(g *gin.Context,msg string)  {
	g.JSON(http.StatusOK,gin.H{
		"code":http.StatusInternalServerError,
		"data":"",
		"msg":msg,
	})
}
