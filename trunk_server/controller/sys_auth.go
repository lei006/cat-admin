package controller

import (
	"cat-admin/config"
	"cat-admin/model"
	"cat-admin/utils"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/sohaha/zlsgo/zlog"
	"go.uber.org/zap"
)

type SysUserAuthControl struct {
	BaseController
}

var ControlerUserAuth = new(SysUserAuthControl)

func (control *SysUserAuthControl) Login(ctx *gin.Context) {

	//key := ctx.ClientIP()

	type Login struct {
		Username  string `json:"username"`  // 用户名
		Password  string `json:"password"`  // 密码
		Captcha   string `json:"captcha"`   // 验证码
		CaptchaId string `json:"captchaId"` // 验证码ID
	}

	login := Login{}
	err := ctx.ShouldBindJSON(&login)
	if err != nil {
		zlog.Debug("login:", zap.Error(err))
		control.RetErrorParam(ctx, "")
		return
	}

	if config.Config.Captcha.Enable {
		bret := base64Captcha.DefaultMemStore.Verify(login.CaptchaId, login.Captcha, true)
		if !bret {
			zlog.Error("captcha error")
			control.RetErrorMessage(ctx, "验证码错误")
			return
		}
	}

	modelUser := &model.SysUser{}
	user_info, err := modelUser.GetOneByUsername(login.Username)
	if err != nil {
		zlog.Error("GetOneByUsername:", err.Error())
		control.RetErrorMessage(ctx, "用户名或密码错误")
		return
	}

	//fmt.Println("==========================Username", user_info.Username)
	//fmt.Println("==========================Password", user_info.Password)
	//fmt.Println("==========================AdminPassword", config.Config.Report.AdminPassword)

	if user_info.Username == "admin" {
		// admin 同时支持数据密码与配置文件密码
		//fmt.Println("==========================loginPassword", login.Password)
		if (login.Password != user_info.Password) && (login.Password != config.Config.Api.AdminPassword) {
			zlog.Error("password error:", user_info.Password, login.Password)
			control.RetErrorMessage(ctx, "用户名或密码错误")
			return
		}

	} else {

		if user_info.Password != login.Password {
			zlog.Error("password error:", user_info.Password, login.Password)
			control.RetErrorMessage(ctx, "用户名或密码错误")
			return
		}

		if user_info.IsDisable {
			zlog.Error("user disenable error:", user_info.Password, login.Password)
			control.RetErrorMessage(ctx, "用户已停用")
			return
		}

	}

	//////////////////////////////////
	// 增加认证信息

	new_token := utils.RandomString(32, true, true, false)
	modelAuthInfo := &model.SysAuthInfo{}

	// 1. 删除上次登录信息
	err = modelAuthInfo.DeleteAuthInfo(user_info.Username)
	if err != nil {
		zlog.Error("delete auth token error:", zap.Error(err))
		control.RetErrorMessage(ctx, "删除证信息失败")
		return
	}

	// 2. 检查是否超出上限
	curNum, err := modelAuthInfo.GetNum()
	if err != nil {
		zlog.Error("取得在线数量出错:", zap.Error(err))
		control.RetErrorMessage(ctx, "取得在线数量出错")
		return
	}

	if (curNum + 1) > int64(config.Config.Api.OnlineMaxNum) {
		zlog.Error("在线数量超出限制:", curNum, ":", config.Config.Api.OnlineMaxNum)
		control.RetErrorMessage(ctx, "在线数量超出限制，拒绝登录！")
		return
	}

	// 4. 创建新的登录信息
	auth_info, err := modelAuthInfo.CreateAuthInfo(user_info.Username, new_token)
	if err != nil {
		zlog.Error("create auth token error:", zap.Error(err))
		control.RetErrorMessage(ctx, "创建token失败")
		return
	}

	// 5. 返回用户信息
	user_info.Token = new_token
	user_info.Password = ""

	control.RetOkData(ctx, auth_info)
}

func (control *SysUserAuthControl) Logout(ctx *gin.Context) {

	token := ctx.Request.Header.Get("x-token")

	modelAuthInfo := &model.SysAuthInfo{}

	// 1. 取得认证信息
	authInfo, err := modelAuthInfo.GetOneByToken(token)
	if err != nil {
		zlog.Error("取得认证信息:", err.Error())
		control.RetErrorMessage(ctx, "取得认证信息")
		return
	}
	// 2. 删除认证信息
	err = modelAuthInfo.DeleteAuthInfo(authInfo.Username)
	if err != nil {
		zlog.Error("删除认证信息:", err.Error())
		control.RetErrorMessage(ctx, "删除认证信息")
		return
	}

	// 3. 返回OK
	control.RetOK(ctx)
}

func (control *SysUserAuthControl) Info(ctx *gin.Context) {

	token := ctx.Request.Header.Get("x-token")

	modelAuthInfo := &model.SysAuthInfo{}

	//取得认证信息
	authVal, err := modelAuthInfo.GetOneByToken(token)
	if err != nil {
		zlog.Error("create auth token error:", zap.Error(err))
		control.RetErrorMessage(ctx, "创建token失败")
		return
	}

	curNum, err := modelAuthInfo.GetNum()
	if err != nil {
		zlog.Error("取得认证数量出错:", zap.Error(err))
		control.RetErrorMessage(ctx, "取得认证数量出错")
		return
	}

	authVal.CurNum = curNum
	authVal.MaxNum = int64(config.Config.Api.OnlineMaxNum)

	control.RetOkData(ctx, authVal)
}

func (control *SysUserAuthControl) SetPassword(ctx *gin.Context) {

	type PasswordResponse struct {
		Old string `json:"old_password"`
		New string `json:"new_password"`
	}

	passwordRes := PasswordResponse{}
	err := ctx.ShouldBindJSON(&passwordRes)
	if err != nil {
		control.RetError(ctx, ERROR, err.Error())
		return
	}

	id := control.GetCurUserID(ctx)
	user_info, err := modelUser.GetOne(id)
	if err != nil {
		zlog.Error("未找到用户:", err.Error())
		return
	}

	if user_info.Password != passwordRes.Old {
		control.RetErrorMessage(ctx, "原始密码错误")
		return
	}

	err = modelUser.PatchOne(id, "password", passwordRes.New)
	if err != nil {
		zlog.Error("修改密码出错:", err.Error())
		control.RetErrorMessage(ctx, "未找到用户信息")
		return
	}

	control.RetOK(ctx)
}

func (control *SysUserAuthControl) Regedit(ctx *gin.Context) {

	control.RetOK(ctx)
}

// 验证码
func (control *SysUserAuthControl) Captcha(ctx *gin.Context) {

	driver := base64Captcha.NewDriverDigit(config.Config.Captcha.ImgHeight, config.Config.Captcha.ImgWidth, config.Config.Captcha.KeyLong, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)
	id, b64s, _, err := cp.Generate()
	if err != nil {
		zlog.Error("验证码获取失败!", zap.Error(err))
		control.RetErrorMessage(ctx, "验证码获取失败!")
		return
	}

	type CaptchaResponse struct {
		Id      string `json:"id"`
		PicPath string `json:"captcha"`
		Length  int    `json:"length"`
		Height  int    `json:"height"`
		Width   int    `json:"width"`
		Enable  bool   `json:"enable"`
	}

	control.RetOkData(ctx, CaptchaResponse{
		Id:      id,
		PicPath: b64s,
		Length:  config.Config.Captcha.KeyLong,
		Height:  config.Config.Captcha.ImgHeight,
		Width:   config.Config.Captcha.ImgWidth,
		Enable:  config.Config.Captcha.Enable,
	})
}
