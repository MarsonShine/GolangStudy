package api

import (
	"context"
	"fmt"
	"kratos-demo/internal/middleware"

	"github.com/go-kratos/kratos/pkg/net/rpc/warden"

	"google.golang.org/grpc"
)

// AppID .
// const AppID = "TODO: ADD APP ID"
const target = "192.168.3.67:9000"

// NewClient new grpc client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (DemoClient, error) {
	client := warden.NewClient(cfg, opts...).Use(middleware.GrpcClientLogging())
	// cc, err := client.Dial(context.Background(), fmt.Sprintf("discovery://default/%s", AppID))
	cc, err := client.Dial(context.Background(), fmt.Sprintf("direct://default/%s", target))
	if err != nil {
		return nil, err
	}
	return NewDemoClient(cc), nil
}

// 直接传递 http.context 会报错，应该附加context值
func NewClientFromHttp(ctx context.Context, cfg *warden.ClientConfig, opts ...grpc.DialOption) (DemoClient, error) {
	client := warden.NewClient(cfg, opts...).Use(middleware.GrpcClientLogging())
	// cc, err := client.Dial(context.Background(), fmt.Sprintf("discovery://default/%s", AppID))
	cc, err := client.Dial(ctx, fmt.Sprintf("direct://default/%s", target))
	if err != nil {
		return nil, err
	}
	return NewDemoClient(cc), nil
}

// 生成 gRPC 代码
//go:generate kratos tool protoc --grpc api.proto
