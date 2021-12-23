package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func SetV() gin.HandlerFunc {
	return func(c *gin.Context) {
		pathInfo := strings.Split(strings.Trim(c.Request.URL.Path, "/"), "/")
		if len(pathInfo) < 4 {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.Set("v1", pathInfo[1])
		c.Set("v2", pathInfo[2])
		c.Set("v3", pathInfo[3])
		c.Next()
	}
}
