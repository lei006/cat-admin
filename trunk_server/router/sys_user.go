package router

import (
	"cat-admin/controller"

	"github.com/gin-gonic/gin"
)

func initRouterSysUser(publicRouter *gin.RouterGroup, privateGroup *gin.RouterGroup) {

	//authPubRouter := publicRouter.Group("user") //.Use(middleware.OperationRecord())
	authPriRouter := privateGroup

	controller := controller.SysUserControl{}

	{
		// add
		authPriRouter.POST("/user", controller.Create) // 新建报告的数据结构
		// del
		authPriRouter.DELETE("/user/:id", controller.DeleteOne) // 删除报告的数据结构
		authPriRouter.DELETE("/user", controller.DeleteMany)    // 批量删除报告的数据结构
		// change
		authPriRouter.PUT("/user/:id", controller.PutOne)     // 更新报告的数据结构
		authPriRouter.PATCH("/user/:id", controller.PatchOne) // 更新报告的数据结构
		// get
		authPriRouter.GET("/user/:id", controller.GetOne) // 根据ID获取报告的数据结构
		authPriRouter.GET("/user", controller.GetList)    // 获取报告的数据结构列表
	}

}
