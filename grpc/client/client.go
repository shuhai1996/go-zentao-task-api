package main

import (
	"context"
	"flag"
	"fmt"
	pb "go-zentao-task-api/grpc/proto/hello"
	pbm "go-zentao-task-api/grpc/proto/math"
	"google.golang.org/grpc"
	"log"
	"time"
)


var (
	addr = flag.String("addr", "127.0.0.1:8899", "the address to connect to")
	name = flag.String("name", "world", "Name to greet")
	num1 = flag.Int64("num1", 1, "num1 to greet")
	num2 = flag.Int64("num2", 2, "num2 to greet")
)


func main() {
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())

	if err != nil {
		log.Fatal("dialing"+ err.Error())
	}
	fmt.Println(conn.GetState().String())
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	m := pbm.NewMathGreeterClient(conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
	rm, err := m.Calculate(ctx, &pbm.Num{Name1: *num1, Name2: *num2})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %v", rm.GetMessage())
}