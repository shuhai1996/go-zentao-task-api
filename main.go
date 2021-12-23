package main

import (
	"go-zentao-task/core"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//服务监听
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	var env = "development"
	go core.RunServer(env, quit) //运行服务
	<-quit
	log.Println("关闭服务 ...")
	log.Println("服务已退出")
}
