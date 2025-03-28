package middleware

import (
	"go-admin/config"

	"github.com/gin-gonic/gin"
)

//var casbinService = service.CasbinServiceApp

// CasbinHandler 拦截器
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.Config.System.Env != "develop" {
			// 解码: jwt
			//waitUse, _ := controller.BaseControllerApp.GetClaims(c)

			/*
				//获取请求的PATH
				path := c.Request.URL.Path
				obj := strings.TrimPrefix(path, config.Config.System.RouterPrefix)
				// 获取请求方法
				act := c.Request.Method
				// 获取用户的角色
				sub := strconv.Itoa(int(waitUse.AuthorityId))
				e := casbinService.Casbin() // 判断策略中是否存在
				success, _ := e.Enforce(sub, obj, act)
				if !success {
					controller.BaseControllerApp.FailWithDetailed(gin.H{}, "权限不足", c)
					c.Abort()
					return
				}
			*/
		}
		c.Next()
	}
}
