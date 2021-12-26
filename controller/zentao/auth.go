package zentao

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-zentao-task-api/core"
	"go-zentao-task-api/service/zentao"
)

type Auth struct {
}

var service = zentao.InitializeService()

type AuthLoginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 登录
func (*Auth) Login(c *core.Context) {
	var r AuthLoginRequest
	fmt.Println(c.Request.Body)
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Fail(40100, "缺少参数，请重试", nil)
		return
	}
	res, err := service.UserLogin(r.Account, r.Password)
	if err != nil {
		c.Fail(400, err.Error(), res)
		return
	}

	if res.Code != 200 {
		c.Fail(400, res.Msg, res)
		return
	}

	c.Success(res.Data)
}

func (*Auth) GetAuthInfo(c *core.Context) {
	userInfo := c.GetStringMap("user_info")
	c.Success(gin.H{
		"user_info": gin.H{
			"account": userInfo["account"],
			"type":    userInfo["type"],
		},
	})
}

func (*Auth) Logout(c *core.Context) {
	//清除用户token数据
	_, err := service.Logout(c.GetHeader("Authorization-Token"))
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	c.Success(nil)
}
