package service

import (
	"context"
	pb "go-zentao-task-api/grpc/proto/hello"
	"log"
)

// Server server 用于实现 HelloServiceServer 接口
type Server struct {
	pb.UnimplementedGreeterServer
}

// SayHello Hello implements helloworld.GreeterServer
func (m *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}
