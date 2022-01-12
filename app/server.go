package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/soheilhy/cmux"
	pb "go-zentao-task-api/grpc/proto/hello"
	pbm "go-zentao-task-api/grpc/proto/math"
	"go-zentao-task-api/grpc/service"
	"go-zentao-task-api/pkg/config"
	"go-zentao-task-api/pkg/db"
	"go-zentao-task-api/pkg/elasticsearch"
	"go-zentao-task-api/pkg/gredis"
	"go-zentao-task-api/pkg/logging"
	"go-zentao-task-api/router"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"strings"
)

func setup(env string) {
	config.Setup(env)//加载配置文件
	logging.Setup(env, logging.Stdout)
	db.Setup()     //初始化数据库
	gredis.Setup() //初始化缓存
	elasticsearch.Setup() //初始化es连接
}

func RunServer(env string) { //服务运行
	setup(env)
	gin.SetMode(gin.DebugMode)
	// 创建一个监听器
	fmt.Println("开始创建一个监听器")
	l, err := net.Listen("tcp4", "127.0.0.1:8899")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("创建成功")
	// 创建一个端口多路复用连接器
	m := cmux.New(l)
	//匹配http 的grpc请求
	grpWs := m.Match(cmux.HTTP1HeaderField("content-type", "application/grpc+json"))
	//匹配http1 请求
	httpL := m.Match(cmux.HTTP1Fast())
	//匹配任意类型请求
	grpcL := m.Match(cmux.Any())

	grpcS := grpc.NewServer()
	fmt.Println(grpcS.GetServiceInfo())
	//注册服务
	pb.RegisterGreeterServer(grpcS, &service.Server{})

	//基于grpc getaway 的restful api
	mux := runtime.NewServeMux()
	pbm.RegisterMathGreeterHandlerServer(context.Background(), mux, &service.MathServer{})
	gwS := &http.Server{
		Handler: mux,
	}

	//创建gin应用, 先匹配http1 请求
	r:= router.Register(env)
	r.SetTrustedProxies([]string{"192.168.1.2"})
	go r.RunListener(httpL)//监听http请求
	// 启动grpc服务
	go gwS.Serve(grpWs)
	go grpcS.Serve(grpcL)

	// 开始端口复用服务，将在同一个端口上提供了grpc和http的服务
	if err := m.Serve(); !strings.Contains(err.Error(), "use of closed network connection") {
		panic(err)
	}
}

