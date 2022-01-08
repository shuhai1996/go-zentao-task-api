package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
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
	gin.SetMode(gin.ReleaseMode)
	// 创建一个监听器
	fmt.Println("开始创建一个监听器")
	l, err := net.Listen("tcp4", "127.0.0.1:8899")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("创建成功")
	// 创建一个端口多路复用连接器
	m := cmux.New(l)
	// 创建grpc服务
	grpcL := m.Match(cmux.Any())
	grpcS := grpc.NewServer()
	fmt.Println(grpcS.GetServiceInfo())
	//注册服务
	pb.RegisterGreeterServer(grpcS, &service.Server{})
	pbm.RegisterMathGreeterServer(grpcS,&service.MathServer{})

	// 启动grpc服务
	go grpcS.Serve(grpcL)

	// 创建gin应用, 先匹配http1 请求
	httpL := m.Match(cmux.HTTP1Fast())
	r:= router.Register(env)
	go func() {
		r.RunListener(httpL)//监听http请求
	}()


	// 开始端口复用服务，将在同一个端口上提供了grpc和http的服务
	if err := m.Serve(); !strings.Contains(err.Error(), "use of closed network connection") {
		panic(err)
	}
}
