package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "example.com/learn-grpc-01/ecommerce"
	"example.com/learn-grpc-01/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func main() {
	creds, _ := credentials.NewClientTLSFromFile(
		"/root/workspace/learn-grpc-01/key/test.pem",
		"*.heliu.site",
	)

	var opts []grpc.DialOption
	// 不带TLS这里是grpc.WithTransportCredentials(insecure.NewCredentials())
	opts = append(opts, grpc.WithTransportCredentials(creds))
	opts = append(opts, grpc.WithPerRPCCredentials(auth.BasicAuth{
		Username: "root",
		Password: "123",
	}))

	// 连接server端，使用ssl加密通信
	conn, err := grpc.NewClient("127.0.0.1:9090", opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewOrderManagementClient(conn)

	fmt.Printf("now-Time: %s\n", time.Now().Format(time.DateTime))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 执行rpc调用(这个方法在服务器端来实现并返回结构)
	resp, err := client.GetOrder(ctx, &wrapperspb.StringValue{Value: "2"})

	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	// id:"2"  items:"6"  items:"5"  items:"4"  items:"3"  items:"2"  items:"1"  destination:"102"
	fmt.Println(resp.String())
}
