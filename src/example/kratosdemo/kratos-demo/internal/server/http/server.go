package http

import (
	"net/http"

	pb "kratos-demo/api"
	"kratos-demo/internal/middleware"
	"kratos-demo/internal/middleware/cors"
	"kratos-demo/internal/model"

	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/log"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
)

var svc pb.DemoServer

// New new a bm server.
func New(s pb.DemoServer) (engine *bm.Engine, err error) {
	var (
		cfg bm.ServerConfig
		ct  paladin.TOML
	)
	if err = paladin.Get("http.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Server").UnmarshalTOML(&cfg); err != nil {
		return
	}
	svc = s
	engine = bm.DefaultServer(&cfg)
	pb.RegisterDemoBMServer(engine, s)
	engine.Use(middleware.NewRecordRequestElapsedTime())
	// 跨域
	cors.NewCors().UseCros(engine)
	initRouter(engine)
	// 限流
	middleware.UseRateLimiter(engine)
	err = engine.Start()
	return
}

func initRouter(e *bm.Engine) {
	e.Ping(ping)
	g := e.Group("/kratos-demo")
	{
		g.GET("/start", howToStart)
		g.GET("/user/:id", getUserHandler)
		g.GET("/bigint", getInt64)
		g.GET("/bigint2", getInt64FromProtobuf)
	}
}

func ping(ctx *bm.Context) {
	if _, err := svc.Ping(ctx, nil); err != nil {
		log.Error("ping error(%v)", err)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

// example for http request handler.
func howToStart(c *bm.Context) {
	k := &model.Kratos{
		Hello: "Golang 大法好 !!!",
	}
	c.JSON(k, nil)
}
