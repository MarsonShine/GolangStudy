package grpc

import (
	pb "kratos-demo/api"
	"kratos-demo/internal/middleware"

	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/net/rpc/warden"
)

// New new a grpc server.
func New(svc pb.DemoServer) (ws *warden.Server, err error) {
	var (
		cfg warden.ServerConfig
		ct  paladin.TOML
	)
	if err = paladin.Get("grpc.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Server").UnmarshalTOML(&cfg); err != nil {
		return
	}
	ws = warden.NewServer(&cfg)
	ws.Use(middleware.GrpcServerLogging())
	// grpcServer := ws.Server()
	// grpc.WithUnaryInterceptor(logInterceptor)
	pb.RegisterDemoServer(ws.Server(), svc)
	ws, err = ws.Start()
	return
}
