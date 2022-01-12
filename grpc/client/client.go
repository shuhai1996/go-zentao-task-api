package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	pb "go-zentao-task-api/grpc/proto/hello"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net/http"
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
	c:= pb.NewGreeterClient(conn)
	//m := pbm.NewMathGreeterClient(conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
	//rm, err := m.Calculate(ctx, &pbm.Num{Name1: *num1, Name2: *num2})
	//if err != nil {
	//	log.Fatalf("could not greet: %v", err)
	//}
	//log.Printf("Greeting: %v", rm.GetMessage())

	// 构造POST请求
	uri := "http://127.0.0.1:8899/rpc/math/calculate"
	params := []byte(`{
    "name1": 4,
    "name2": 6
	}`)
	re, _:= http.NewRequest("POST", uri, bytes.NewReader(params))
	re.Header.Add("Content-Type", "application/grpc+json")
	client := &http.Client{}
	resp, err := client.Do(re)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()

	statuscode := resp.StatusCode
	hea := resp.Header
	bodyw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(hea)
	fmt.Println(statuscode)
	fmt.Println(string(bodyw))
}