package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kratos-demo/internal/di"

	"github.com/MSLibs/glogger"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/log"
)

func main() {
	flag.Parse()
	// log.Init(nil) // debug flag: log.dir={path}
	// defer log.Close()
	// log.Info("kratos-demo start")
	log := glogger.CreateLog(glogger.GLoggerConfig{})
	paladin.Init()
	_, closeFunc, err := di.InitApp()
	if err != nil {
		panic(err)
	}

	// // 服务发现
	// bulder, err := mynaming.NewConsulDiscovery(mynaming.Config{Zone: "zone01", Env: "dev", Region: "region01"})
	// if err != nil {
	// 	panic(err)
	// }
	// resolver.Register(bulder)

	// // 服务注册
	// ip := "127.0.0.1" // NOTE: 必须拿到您实例节点的真实IP，
	// port := "9000"    // NOTE: 必须拿到您实例grpc监听的真实端口，warden默认监听9000
	// hn, _ := os.Hostname()
	// ins := &naming.Instance{
	// 	Zone:     "zone01",
	// 	Env:      env.DeployEnv,
	// 	AppID:    AppID,
	// 	Hostname: hn,
	// 	Addrs: []string{
	// 		"grpc://" + ip + ":" + port,
	// 	},
	// }

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Infof("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeFunc()
			log.Info("kratos-demo exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func logInit() {
	log.Init(&log.Config{
		Dir: "log.log",
		V:   1,
	})
}
