package app
import (
	"github.com/gin-gonic/gin"
	"github.com/soheilhy/cmux"
	"go-zentao-task-api/pkg/config"
	"go-zentao-task-api/pkg/db"
	"go-zentao-task-api/pkg/elasticsearch"
	"go-zentao-task-api/pkg/gredis"
	"go-zentao-task-api/pkg/logging"
	"go-zentao-task-api/router"
	"google.golang.org/grpc"
	"log"
	"net"
)

func setup(env string) {
	config.Setup(env)
	logging.Setup(env, logging.Stdout)
	db.Setup()     //初始化数据库
	gredis.Setup() //初始化缓存
	elasticsearch.Setup() //初始化es连接
}

func RunServer(env string) { //服务运行
	setup(env)
	gin.SetMode(gin.ReleaseMode)
	// Create the main listener.
	l, err := net.Listen("tcp", ":8899")
	if err != nil {
		log.Fatal(err)
	}

	// Create a cmux.
	m := cmux.New(l)

	// grpc
	grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	grpcS := grpc.NewServer()

	go grpcS.Serve(grpcL)

	// gin
	httpL := m.Match(cmux.HTTP1Fast())
	r:= router.Register(env)
	go func() {
		r.RunListener(httpL)//监听
	}()

	// Start serving!
	m.Serve()
}
