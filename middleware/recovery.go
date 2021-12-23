package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go-zentao-task/pkg/logging"
	"net/http"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var errMsg string
				switch err := err.(type) {
				case error:
					errMsg = fmt.Sprintf("%+v", errors.WithStack(err.(error)))
				default:
					errMsg = fmt.Sprintf("%v", err)
				}
				logging.Fatal("go", "middleware", "recovery", "服务器内部错误", "", errMsg)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
