package middleware

import (
	"net/http"
	"strconv"

	"github.com/sohaha/zlsgo/zlog"

	"go-admin/config"
	"go-admin/controller"
	"go-admin/model"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {

	return func(c *gin.Context) {

		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localStorage中
		// 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := c.Request.Header.Get("x-token")
		if token == "" {
			controller.BaseControllerApp.RetError(c, controller.ERROR_Unauthorized, "未登录或非法访问")
			c.Abort()
			return
		}

		modelAuthInfo := &model.SysAuthInfo{}

		//删除所以过期的token
		err := modelAuthInfo.DeleteAuthInfoByUpdatedTime(config.Config.Api.TokenExpired)
		if err != nil {
			zlog.Error("删除过期认证信息出错:", err.Error())
			controller.BaseControllerApp.RetError(c, controller.ERROR_Unauthorized, "删除过期认证信息出错")
			c.Abort()
			return
		}

		// 更新token
		modelAuthInfoVal, err := modelAuthInfo.UpdateOneByToken(token)
		if err != nil {
			zlog.Error("jwt 未找到登录信息:", err.Error())
			controller.BaseControllerApp.RetError(c, controller.ERROR_Unauthorized, "jwt 未找到登录信息")
			c.Abort()
			return
		}

		//取得认证信息
		modelUser := model.SysUser{}
		user_info, err := modelUser.GetOneByUsername(modelAuthInfoVal.Username)
		if err != nil {
			if model.ErrRecordNotFound(err) {
				zlog.Error("未找到token:", err.Error())
				controller.BaseControllerApp.RetError(c, http.StatusUnauthorized, "token非法")
			} else {
				zlog.Error("其它错误:", err.Error())
				controller.BaseControllerApp.RetError(c, http.StatusInternalServerError, "内部错误")
			}
			c.Abort()
			return
		}

		//检查用户是否被冻结
		if user_info.Username != "admin" && user_info.IsDisable {
			controller.BaseControllerApp.RetError(c, controller.ERROR, "用户已被停用")
			c.Abort()
			return
		}

		//把user_info.ID 转为字符串
		c.Request.Header.Set("user_id", strconv.FormatUint(uint64(user_info.ID), 10))

		if user_info.IsAdmin {
			c.Request.Header.Set("IsAdmin", "true")
		} else {
			c.Request.Header.Set("IsAdmin", "false")
		}

		c.Next()
	}
}
