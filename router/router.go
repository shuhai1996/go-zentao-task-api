package router

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"go-zentao-task-api/controller/zentao"
	"go-zentao-task-api/core"
	"go-zentao-task-api/middleware"
	"net/http"
)

func Register(env string) *gin.Engine {

	r := gin.New()
	r.Use(middleware.Recovery())
	r.Use(static.Serve("/g/static", static.LocalFile("static", false)))

	r.LoadHTMLGlob("template/**/*")
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "error/404", nil)
	})
	r.Any("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "ok",
		})
	})

	r.Use(middleware.SetV())
	r.Use(core.Handle(middleware.Logging))
	r.Use(core.Handle(middleware.I18n))
	r.Use(core.Handle(middleware.Session))
	r.Use(middleware.Secure())

	g := r.Group("/zentao")
	g.POST("/user/auth/login", core.Handle(new(zentao.Auth).Login))
	zt := g.Group(
		"/user",
		core.Handle(middleware.JwtAuth),
	)
	{
		//用户
		zt.GET("/auth/info", core.Handle(new(zentao.Auth).GetAuthInfo))
		zt.POST("/auth/logout", core.Handle(new(zentao.Auth).Logout))

		//任务
		zt.GET("task/list", core.Handle(new(zentao.Task).Index))
	}

	g.GET("/v1/v2/i18n", core.Handle(func(c *core.Context) {
		c.String(http.StatusOK, c.Tr("author.info", 18)+" "+c.Tr("section.language"))
	}))

	g.GET("/v1/v2/session", core.Handle(func(c *core.Context) {
		var count int
		v := c.GetSession("count")
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count++
		}
		c.SetSession("count", count)
		c.SaveSession() //nolint
		c.JSON(http.StatusOK, gin.H{"count": count})
	}))

	return r
}
