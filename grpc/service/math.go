package service

import (
	"context"
	pbm "go-zentao-task-api/grpc/proto/math"
	"log"
)

// MathServer 用于实现 cla  接口
type MathServer struct {
	pbm.UnimplementedMathGreeterServer
}

// Calculate 实现计算方法
func (m *MathServer) Calculate(ctx context.Context, in *pbm.Num) (*pbm.CalReply, error) {
	log.Printf("Recived: %v,%v", in.GetName1(), in.GetName2())
	n:= in.GetName1() + in.GetName2()
	log.Printf("Calculate: %v", n)
	return &pbm.CalReply{Message:n}, nil
}
