package core

import (
	"github.com/gin-gonic/gin"
	"go-zentao-task/pkg/i18n"
	"go-zentao-task/pkg/logging"
	"go-zentao-task/pkg/session"
	"go-zentao-task/service/errcode"
	"net/http"
)

type Context struct {
	*gin.Context
	*i18n.Locale
	*session.Session
	*logging.Logging

	V1 string
	V2 string
	V3 string
}

func (c *Context) Success(data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": data,
	})
}

func (c *Context) Fail(code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func (c *Context) FailWithErrCode(code errcode.ErrCode, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  code.String(),
		"data": data,
	})
}

func (c *Context) ErrorPage(msg string) {
	c.HTML(http.StatusOK, "error/error", gin.H{
		"msg": msg,
	})
}
