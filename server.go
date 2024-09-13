package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	pb "example.com/learn-grpc-01/ecommerce"
	"example.com/learn-grpc-01/pkg/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func main() {
	// 方法一
	// creds, err1 := credentials.NewServerTLSFromFile(
	//		"/root/workspace/learn-grpc/key/test.pem",
	//		"/root/workspace/learn-grpc/key/test.key",
	//	)
	//
	//	if err1 != nil {
	//		fmt.Printf("证书错误：%v", err1)
	//		return
	//	}

	// 方法二
	cert, err := tls.LoadX509KeyPair(
		"/root/workspace/learn-grpc-01/key/test.pem",
		"/root/workspace/learn-grpc-01/key/test.key")
	if err != nil {
		fmt.Printf("私钥错误：%v", err)
		return
	}
	creds := credentials.NewServerTLSFromCert(&cert)

	listen, _ := net.Listen("tcp", ":9090")
	grpcServer := grpc.NewServer(
		grpc.Creds(creds), // TLS
		grpc.UnaryInterceptor(interceptor.UnaryServerAuthInterceptor()), // 拦截器
	)
	pb.RegisterOrderManagementServer(grpcServer, &server{})

	// 启动服务
	err = grpcServer.Serve(listen)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 编译器静态检查，*server是否实现了pb.OrderManagementServer接口
var _ pb.OrderManagementServer = (*server)(nil)

var orders = make(map[string]pb.Order, 8)

func init() {
	// 测试数据
	orders["1"] = pb.Order{Id: "1", Items: []string{"1", "2", "3", "4", "5", "6"}, Destination: "101"}
	orders["2"] = pb.Order{Id: "2", Items: []string{"6", "5", "4", "3", "2", "1"}, Destination: "102"}
}

type server struct {
	pb.UnimplementedOrderManagementServer
}

// GetOrder 获取订单信息
func (s *server) GetOrder(ctx context.Context, orderId *wrapperspb.StringValue) (*pb.Order, error) {
	ord, exists := orders[orderId.Value]
	if exists {
		return &ord, status.New(codes.OK, "").Err()
	}

	return nil, status.Errorf(codes.NotFound, "Order does not exist. : ", orderId)
}
