package interceptor

import (
	"context"
	"encoding/base64"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid credentials")
)

// UnaryServerAuthInterceptor 服务端拦截器
func UnaryServerAuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errMissingMetadata
		}

		// 2024/09/13 16:04:08 md: map[:authority:[*.heliu.site] authorization:[Basic cm9vdDoxMjM=] content-type:[application/grpc] user-agent:[grpc-go/1.66.2]]
		log.Printf("md: %v", md)

		// md.Get("authorization")
		authorization, ok := md["authorization"]
		if !ok || len(authorization) < 1 {
			return nil, errMissingMetadata
		}

		// 2024/09/13 16:04:08 authorization: [Basic cm9vdDoxMjM=]
		log.Printf("authorization: %v", authorization)

		token := strings.TrimPrefix(authorization[0], "Basic ")
		if token != base64.StdEncoding.EncodeToString([]byte("root:123")) {
			return nil, errInvalidToken
		}

		return handler(ctx, req)
	}
}
