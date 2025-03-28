package controller

import (
	"go-admin/model"
	"go-admin/utils"

	"github.com/gin-gonic/gin"
	"github.com/sohaha/zlsgo/zlog"
	"go.uber.org/zap"
)

var modelUser model.SysUser

type SysUserControl struct {
	BaseController
}

func (control *SysUserControl) Create(ctx *gin.Context) {

	user_info := model.SysUser{}
	err := ctx.ShouldBindJSON(&user_info)
	if err != nil {
		zlog.Debug("Create:", zap.Error(err))
		control.RetErrorParam(ctx, "")
		return
	}
	zlog.Debugf("user_info: %+v \n", user_info)
	err = modelUser.Create(&user_info)
	if err != nil {
		zlog.Debug("create new user:", err)
		control.RetErrorParam(ctx, "")
		return
	}

	control.RetOkData(ctx, user_info)
}

func (control *SysUserControl) DeleteOne(ctx *gin.Context) {

	id := ctx.Param("id")

	err := modelUser.DeleteOne(id)
	if err != nil {
		zlog.Debug("Delete user error:", err)
		control.RetErrorParam(ctx, "")
		return
	}

	control.RetOK(ctx)
}

func (control *SysUserControl) DeleteMany(ctx *gin.Context) {
	var ids []uint
	err := ctx.ShouldBindJSON(&ids)
	if err != nil {
		control.RetError(ctx, ERROR, err.Error())
		return
	}
	err = modelUser.DeleteMany(ids)
	if err != nil {
		control.RetError(ctx, ERROR, err.Error())
		return
	}

	control.RetOK(ctx)
}

func (control *SysUserControl) PutOne(ctx *gin.Context) {

	var reportItem model.SysUser
	err := ctx.ShouldBindJSON(&reportItem)
	if err != nil {
		control.RetError(ctx, ERROR, err.Error())
		return
	}
	verify := utils.Rules{
		"ID": {utils.NotEmpty()},
	}
	if err := utils.Verify(reportItem, verify); err != nil {
		control.RetError(ctx, ERROR, err.Error())
		return
	}

	if err := modelUser.UpdateOne(&reportItem); err != nil {
		zlog.Error("更新失败!", zap.Error(err))
		control.RetError(ctx, ERROR, err.Error())
		return
	}
	control.RetOK(ctx)
}

func (control *SysUserControl) PatchOne(ctx *gin.Context) {
	id := ctx.Param("id")
	// 把id 转为 uint

	req := PatchReq{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		control.RetError(ctx, ERROR, err.Error())
		return
	}

	// 1. 管理员可以修改一切
	// 2. 非管理员只能修改自己
	// 3. 非管理员不能用此接口修改密码
	if !control.IsAdmin(ctx) {
		// 如果不是管理员
		if control.GetCurUserID(ctx) != id {
			control.RetError(ctx, ERROR, "非管理员不能修改其他用户信息")
			return
		}
		if req.Field == "password" {
			control.RetError(ctx, ERROR, "非管理员不能用此接口修改密码")
			return
		}
	}

	err = modelUser.PatchOne(id, req.Field, req.Data)
	if err != nil {
		control.RetError(ctx, ERROR, err.Error())
		return
	}
	data, err := modelUser.GetOne(id)
	if err != nil {
		control.RetError(ctx, ERROR, err.Error())
		return
	}
	zlog.Debug("patch:", id, req, data)

	control.RetOkData(ctx, data)
}

func (control *SysUserControl) GetOne(ctx *gin.Context) {
	id := ctx.Param("id")
	item, err := modelUser.GetOne(id)
	if err != nil {

		item_val, err := modelUser.GetOneByUsername(id)
		if err != nil {
			control.RetError(ctx, ERROR, err.Error())
			return
		}
		control.RetOkData(ctx, item_val)
		return
	}

	control.RetOkData(ctx, item)
}

func (control *SysUserControl) GetList(ctx *gin.Context) {

	if control.IsAdmin(ctx) {
		user_list, _, err := modelUser.GetList()
		if err != nil {
			control.RetError(ctx, ERROR, err.Error())
			return
		}

		control.RetOkList(ctx, user_list)
	} else {
		id := control.GetCurUserID(ctx)
		user_list, _, err := modelUser.GetListOne(id)
		if err != nil {
			control.RetError(ctx, ERROR, err.Error())
			return
		}
		control.RetOkList(ctx, user_list)
	}
}
