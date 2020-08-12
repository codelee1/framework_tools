package main

import (
	"flag"
	"go.uber.org/zap"
	"learn_tools/gin_demo/conf"
	"learn_tools/gin_demo/controller"
	"learn_tools/gin_demo/global"
	"learn_tools/gin_demo/router"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	flag.Parse()

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	log := logger.Sugar()

	// 读取配置
	if err := conf.Init(); err != nil {
		log.Errorf("conf.Init() error(%v)", err)
		panic(err)
	}

	controller.InitLogger(conf.Conf.Log.Dir + "/controller.log")
	defer controller.Sync()

	// 初始化全局变量
	global.Init()
	defer global.Close()

	// 初始化路由
	router.Init()

	// 优雅退出
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-ch
		log.Infof("canal get a signal %s", s.String())
		switch s {
		case os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			//Close()
			log.Info("app exit")
			time.Sleep(time.Second)
			os.Exit(0)
			return
		case syscall.SIGHUP:
		default:
			return
		}

	}
}
