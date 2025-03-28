package router

import (
	"go-admin/controller"

	"github.com/gin-gonic/gin"
)

func initRouterSysAuth(publicRouter *gin.RouterGroup, privateGroup *gin.RouterGroup) {

	authPubRouter := publicRouter.Group("auth") //.Use(middleware.OperationRecord())
	authPriRouter := privateGroup.Group("auth")

	baseApi := controller.ControlerUserAuth

	authPubRouter.POST("login", baseApi.Login)
	authPriRouter.DELETE("logout", baseApi.Logout)
	authPriRouter.PATCH("password", baseApi.SetPassword)
	authPubRouter.POST("regedit", baseApi.Regedit)
	authPriRouter.GET("info", baseApi.Info)
	authPubRouter.GET("captcha", baseApi.Captcha)
}
