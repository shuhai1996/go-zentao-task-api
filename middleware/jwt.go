package middleware

import (
	"go-zentao-task-api/core"
	"go-zentao-task-api/service/errcode"
	"go-zentao-task-api/service/zentao"
)

var service = zentao.InitializeService()

type AuthHeader struct {
	Token string `header:"Authorization-Token" binding:"required"`
}

// JwtAuth 后台鉴权中间件
func JwtAuth(c *core.Context) {
	var h AuthHeader
	if err := c.ShouldBindHeader(&h); err != nil {
		c.FailWithErrCode(errcode.ErrAdminLoginExpired, nil)
		c.Abort()
	}
	res, err := service.GetUserInfoByToken(h.Token)
	if err != nil || res.Code != 200 {
		c.FailWithErrCode(errcode.ErrAdminLoginExpired, nil)
		c.Abort()
	}
	userInfo := res.Data
	c.Set("user_info", userInfo)
	c.Next()
}
