package router

import (
	"cat-admin/config"
	"cat-admin/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sohaha/zlsgo/zlog"
)

func LoadRouter(engine *gin.Engine) error {

	engine.Use(middleware.CorsByRules()) // 按照配置的规则放行跨域请求

	{
		// 应映页面
		staticPath := config.WorkPath + "/static"
		engine.Static("/static", staticPath)
		zlog.Info("map : /static -- >", staticPath)
	}

	{
		viewPath := config.WorkPath + "/view"
		engine.StaticFile("/index.html", viewPath+"/index.html")
		engine.StaticFile("/favicon.ico", viewPath+"/favicon.ico")
		engine.StaticFile("/", viewPath+"/index.html")
		engine.StaticFS("/assets", http.Dir(viewPath+"/assets"))
		zlog.Info("map : / -- >", viewPath)

	}

	{

		// 检查健康测试
		engine.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
	}

	publicGroup := engine.Group(config.Config.System.RouterPrefix)  // 公有路由 （无权限检查）
	privateGroup := engine.Group(config.Config.System.RouterPrefix) // 私有路由 （权限检查）
	{
		privateGroup.Use(middleware.JWTAuth()) //支持JWT
		//privateGroup.Use(middleware.CasbinHandler()) //权限管理
	}

	{

		//
		initRouterSysSetup(publicGroup, privateGroup)  // 路由系统设置
		initRouterSysAuth(publicGroup, privateGroup)   // 路由授权
		initRouterSysUser(publicGroup, privateGroup)   // 路由用户管理
		initRouterSysOption(publicGroup, privateGroup) // 路由操作记录
	}

	return nil
}
