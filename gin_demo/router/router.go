package router

import (
	"gin_demo/conf"
	"gin_demo/controller"
	"gin_demo/library/cors"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

func Init() {
	f, _ := os.Create(conf.Conf.Log.Dir +"/run.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	// 用于开发测试
	gin.SetMode(gin.DebugMode)
	//gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.Cors())

	notice := new(controller.NoticeController)
	n := r.Group("/notice")
	{
		n.POST("/",notice.Add)
		n.DELETE("/:id",notice.Del)
		n.PUT("/:id",notice.Update)
		n.GET("/:id",notice.Get)
		n.GET("/",notice.GetAll)
	}

	if err := r.Run(conf.Conf.HttpServer.Addr); err != nil {
		panic(err)
	}
}


