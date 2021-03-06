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
	r.Use(Cors())
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
		zt.GET("auth/info", core.Handle(new(zentao.Auth).GetAuthInfo))
		zt.POST("auth/logout", core.Handle(new(zentao.Auth).Logout))

		//任务
		zt.GET("task/list", core.Handle(new(zentao.Task).Index))

		//行为
		zt.GET("action/list", core.Handle(new(zentao.Actions).Index))
		zt.GET("action/:id", core.Handle(new(zentao.Actions).View))
		zt.DELETE("action/:id", core.Handle(new(zentao.Actions).Delete))
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

// 跨域中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
