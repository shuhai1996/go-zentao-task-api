package middleware

import (
	"go-zentao-task/core"
	"net/http"
)

func CORS(c *core.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Length, Content-Type, X-CSRF-Token")
	c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT")

	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
}
