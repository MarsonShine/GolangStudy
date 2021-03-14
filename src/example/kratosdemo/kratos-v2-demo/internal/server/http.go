package server

import (
	"bytes"
	"encoding/gob"
	v1 "kratos-v2-demo/api/helloworld/v1"
	"kratos-v2-demo/commons/middlewares"
	"kratos-v2-demo/internal/conf"
	"kratos-v2-demo/internal/service"
	"log"
	http1 "net/http"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService) *http.Server {
	var opts = []http.ServerOption{}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	m := http.Middleware(
		middleware.Chain(
			recovery.Recovery(),
			tracing.Server(),
			logging.Server(),
			middlewares.Server(),
		),
	)
	m = http.ResponseEncoder(RegisterDataResponseEncoder)
	srv.HandlePrefix("/", v1.NewGreeterHandler(greeter, m))
	return srv
}

func RegisterDataResponseEncoder(w http1.ResponseWriter, r *http1.Request, body interface{}) error {
	dataWrapper := DataResponseWrapper{
		Success: true,
	}
	dataWrapper.Data = body

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(dataWrapper)
	// buffer := &bytes.Buffer{}
	// err := binary.Write(buffer, binary.LittleEndian, dataWrapper)
	if err != nil {
		log.Fatalf("错误，%v", err)
	}
	bytes := buffer.Bytes()
	w.Write(bytes)
	return nil
}

type DataResponseWrapper struct {
	Success bool
	Message string
	Code    int
	Data    interface{}
}
