package middleware

import (
	"go-zentao-task/core"
	"go-zentao-task/pkg/logging"
	"go-zentao-task/pkg/rbac"
	"go-zentao-task/service/errcode"
)

func RBAC(c *core.Context) {
	userInfo := c.GetStringMap("user_info")
	if userInfo == nil {
		c.FailWithErrCode(errcode.ErrAdminLoginExpired, nil)
		return
	}

	if userInfo["role_type"].(string) == "super" {
		c.Next()
		return
	}

	sub := userInfo["account"].(string)
	obj := c.Request.URL.Path
	act := c.Request.Method

	if err := rbac.Enforcer.LoadPolicy(); err != nil {
		logging.Fatal("go", "middleware", "rbac", "LoadPolicy失败", "", err.Error())
		c.FailWithErrCode(errcode.ErrAdminNetworkBusy, nil)
		return
	}

	ok, err := rbac.Enforcer.Enforce(sub, obj, act)
	if err != nil {
		logging.Fatal("go", "middleware", "rbac", "Enforce失败", "", err.Error())
		c.FailWithErrCode(errcode.ErrAdminNetworkBusy, nil)
		return
	}
	if !ok {
		c.FailWithErrCode(errcode.ErrAdminResourceForbidden, nil)
		return
	}
	c.Next()
}
