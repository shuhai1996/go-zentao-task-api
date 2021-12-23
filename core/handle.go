package core

import (
	"github.com/gin-gonic/gin"
)

const Ctx = "__context__"

type HandlerFunc func(c *Context)

func Handle(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := getContext(c)
		h(ctx)
	}
}

func getContext(c *gin.Context) *Context {
	ctx, ok := c.Get(Ctx)
	if ok {
		return ctx.(*Context)
	}

	context := &Context{
		Context: c,
		V1:      c.GetString("v1"),
		V2:      c.GetString("v2"),
		V3:      c.GetString("v3"),
	}
	c.Set(Ctx, context)
	return context
}
