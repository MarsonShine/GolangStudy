package main

import (
	"flag"

	pb "kratos-v2-demo/api/helloworld/v1"
	"kratos-v2-demo/internal/biz"
	"kratos-v2-demo/internal/conf"
	"kratos-v2-demo/internal/data"
	"kratos-v2-demo/internal/service"

	"github.com/devfeel/mapper"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"gopkg.in/yaml.v2"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
	createMapper()
}

func createMapper() {
	mapper.Register(&biz.UserDto{})
	mapper.Register(&data.UserEntity{})
}

func newApp(logger log.Logger, hs *http.Server, gs *grpc.Server, greeter *service.GreeterService) *kratos.App {
	pb.RegisterGreeterServer(gs, greeter)
	pb.RegisterGreeterHTTPServer(hs, greeter)
	return kratos.New(
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			gs,
		),
	)
}

func main() {
	flag.Parse()
	logger := log.NewStdLogger()
	// log := log.NewHelper("main", logger)
	config := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
		config.WithDecoder(func(kv *config.KeyValue, v map[string]interface{}) error {
			return yaml.Unmarshal(kv.Value, v)
		}),
	)

	// val := config.Value("server")
	// log.Infof("config:server=%v", val)
	if err := config.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := config.Scan(&bc); err != nil {
		panic(err)
	}

	app, err := initApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
