package docs

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-zentao-task/pkg/config"
	"io/ioutil"
	"net/http"
	"os"
)

type SwaggerRequest struct {
	Module    string `form:"module" binding:"required"`
	Timestamp string `form:"timestamp" binding:"required"`
	Sign      string `form:"sign" binding:"required"`
}

func Swagger(c *gin.Context) {
	var r SwaggerRequest
	if err := c.ShouldBindQuery(&r); err != nil {
		c.Abort()
		return
	}

	if fmt.Sprintf("%x", md5.Sum([]byte(r.Timestamp+config.Get("swagger.api_secret")))) != r.Sign {
		c.Abort()
		return
	}

	file, err := os.Open("docs/swagger/" + r.Module + ".yaml")
	if err != nil {
		c.Abort()
		return
	}

	byts, err := ioutil.ReadAll(file)
	if err != nil {
		c.Abort()
		return
	}
	c.String(http.StatusOK, string(byts))
}
