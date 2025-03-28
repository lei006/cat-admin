package controller

import (
	"cat-admin/model"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

type SysSetupControl struct {
	BaseController
}

var modelSetup model.SysSetup

func (control *SysSetupControl) PatchOne(ctx *gin.Context) {
	id := ctx.Param("id")
	// 把id 转为 uint

	req := PatchReq{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		control.RetError(ctx, ERROR, err.Error())
		return
	}

	err = modelSetup.PatchOne(id, req.Field, req.Data)
	if err != nil {
		control.RetError(ctx, ERROR, err.Error())
		return
	}

	item, err := modelSetup.GetOne(id)
	if err != nil {
		control.RetError(ctx, ERROR, err.Error())
		return
	}

	control.RetOkData(ctx, item)
}

func (control *SysSetupControl) GetOneById(ctx *gin.Context) {
	id := ctx.Param("id")

	item, err := modelSetup.GetOne(id)
	if err != nil {
		control.RetError(ctx, ERROR, err.Error())
		return
	}

	control.RetOkData(ctx, item)
}

func (control *SysSetupControl) PutOne(ctx *gin.Context) {

	id := ctx.Param("id")

	body, _ := io.ReadAll(ctx.Request.Body)

	fmt.Println("SysSetupControl SetOne id:", id, string(body))

	//检查是否存在，如果不存在则创建，如果存在则更新
	_, err := modelSetup.GetOne(id)
	if err != nil {
		val, err := modelSetup.CreateOne(id, string(body))
		if err != nil {
			control.RetError(ctx, ERROR, err.Error())
			return
		}
		control.RetOkData(ctx, val)
	} else {
		//更新数据
		err = modelSetup.PatchOne(id, "data", string(body))
		if err != nil {
			control.RetError(ctx, ERROR, err.Error())
			return
		}
		control.RetOK(ctx)
	}
}

func (control *SysSetupControl) GetOneByName(ctx *gin.Context) {
	name := ctx.Param("name")
	item, err := modelSetup.GetOne(name)
	if err != nil {
		new_val := &model.SysSetup{Name: name}
		err := modelSetup.Create(new_val)
		if err != nil {
			control.RetError(ctx, ERROR, err.Error())
		} else {
			control.RetOkData(ctx, new_val)
		}
	} else {
		control.RetOkData(ctx, item)
	}

}

func (control *SysSetupControl) PutOneByName(ctx *gin.Context) {
	name := ctx.Param("name")
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		control.RetError(ctx, ERROR, err.Error())
		return
	}

	val, err := modelSetup.GetOne(name)
	if err != nil {
		val, err := modelSetup.CreateOne(name, string(bodyBytes))
		if err != nil {
			control.RetError(ctx, ERROR, err.Error())
			return
		}
		control.RetOkData(ctx, val)
	} else {
		//更新数据
		err = modelSetup.PatchOne(fmt.Sprintf("%d", val.ID), "data", string(bodyBytes))
		if err != nil {
			control.RetError(ctx, ERROR, err.Error())
			return
		}
		control.RetOK(ctx)
	}

}

func (control *SysSetupControl) GetList(ctx *gin.Context) {

	items, _, err := modelSetup.GetList()
	if err != nil {
		control.RetError(ctx, ERROR, err.Error())
		return
	}

	control.RetOkList(ctx, items)
}
