package app

import (
	"github.com/gin-gonic/gin"
	"go-zentao-task-api/pkg/config"
	"go-zentao-task-api/pkg/db"
	"go-zentao-task-api/pkg/gredis"
	"go-zentao-task-api/pkg/logging"
	"go-zentao-task-api/router"
	"log"
	"net/http"
)

func setup(env string) {
	config.Setup(env)
	logging.Setup(env, logging.Stdout)
	db.Setup()     //初始化数据库
	gredis.Setup() //初始化缓存
	//rbac.Setup()
}

func RunServer(env string) { //服务运行
	setup(env)
	gin.SetMode(gin.ReleaseMode)
	r := router.Register(env)
	srv := &http.Server{
		Addr:    ":8899",
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}
