package main

import (
	"go-admin/config"
	"go-admin/logger"
	"go-admin/model"
	"go-admin/router"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lei006/go-daemon/daemontool"
	"github.com/sohaha/zlsgo/zlog"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

// 目录结构：
// controller  控制器
// model       模型
// router      路由
// entity      结构体
// schedule    定时任务
// middleware  中间插件层
// extend      扩展--
// units       小组件

func main() {

	///////////////////////////////////////////
	// 如果 参数 alone 则不是服务模式
	if len(os.Args) == 2 {
		if os.Args[1] == "alone" {
			config.IsDaemon = false
		}
	}

	err := config.OnInit()
	if err != nil {
		zlog.Error(err)
		return
	}

	if config.RunAtVscode {
		config.IsDaemon = false
	}

	///////////////////////////////////////////
	// 如果 开始程序

	if config.IsDaemon {
		daemonTool := daemontool.DefDaemonTool
		daemonTool.Run(config.AppName, config.AppDesc, RunApp)
	} else {
		RunApp()
	}

}

func RunApp() {
	zlog.Debug("app RunApp enter")

	// 加载配置-日志系统
	err := logger.OnInit()
	if err != nil {
		return
	}
	err = model.OnInit()
	if err != nil {
		return
	}
	if !config.Config.Api.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()

	// 加载-路由
	err = router.LoadRouter(engine)
	if err != nil {
		return
	}

	addr := config.Config.Api.Addr
	zlog.Info("WebServer Listen :" + addr)
	err = engine.Run(addr)
	if err != nil {
		zlog.Debug("WebServer error :" + err.Error())
	}

}
