package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go-zentao-task/core"
	"go-zentao-task/pkg/logging"
	"io/ioutil"
)

type BodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *BodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Logging(c *core.Context) {
	byts, _ := c.GetRawData()
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(byts))

	w := &BodyLogWriter{ResponseWriter: c.Writer, body: bytes.NewBufferString("")}
	c.Writer = w

	v1 := c.GetString("v1")
	v2 := c.GetString("v2")
	v3 := c.GetString("v3")

	c.Logging = &logging.Logging{
		V1: v1,
		V2: v2,
		V3: v3,
	}

	c.Next()
	logging.Info(v1, v2, v3, "接口请求与响应", string(byts), w.body.String())
}
