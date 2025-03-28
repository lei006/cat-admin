package router

import (
	"go-admin/controller"

	"github.com/gin-gonic/gin"
)

func initRouterSysSetup(publicRouter *gin.RouterGroup, privateGroup *gin.RouterGroup) {

	//optionRouter := privateGroup

	controller := controller.SysSetupControl{}

	{
		// add
		//optionRouter.POST("/setup", controller.Create) // 新建报告的数据结构
		// del
		//authPriRouter.DELETE("/user/:id", controller.DeleteOne) // 删除报告的数据结构
		//optionRouter.DELETE("/setup", controller.DeleteMany) // 批量删除报告的数据结构
		// change
		//optionRouter.PUT("/setup/:id", controller.PutOne)     // 更新报告的数据结构
		publicRouter.PATCH("/setup/:id", controller.PatchOne) // 更新报告的数据结构
		// get
		publicRouter.GET("/setup/:id", controller.GetOneById)          // 根据ID获取报告的数据结构
		publicRouter.GET("/setup_name/:name", controller.GetOneByName) // 根据ID获取报告的数据结构
		publicRouter.PUT("/setup_name/:name", controller.PutOneByName) // 根据ID获取报告的数据结构
		publicRouter.GET("/setup", controller.GetList)                 // 获取报告的数据结构列表
		publicRouter.PUT("/setup/:id", controller.PutOne)              // 获取报告的数据结构列表
	}

}
