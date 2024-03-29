package service

import (
	"context"
	"fmt"
	"math/rand"

	pb "kratos-demo/api"
	"kratos-demo/internal/dao"

	"github.com/MSLibs/glogger"
	"github.com/go-kratos/kratos/pkg/conf/paladin"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/wire"
)

var Provider = wire.NewSet(New, wire.Bind(new(pb.DemoServer), new(*Service)))

// Service service.
type Service struct {
	ac  *paladin.Map
	dao dao.Dao
	log glogger.GLogger
}

// New new a service and return.
func New(d dao.Dao) (s *Service, cf func(), err error) {
	s = &Service{
		ac:  &paladin.TOML{},
		dao: d,
		log: *dao.CreateLogger(),
	}
	cf = s.Close
	err = paladin.Watch("application.toml", s.ac)
	return
}

// SayHello grpc demo func.
func (s *Service) SayHello(ctx context.Context, req *pb.HelloReq) (reply *empty.Empty, err error) {
	reply = new(empty.Empty)
	s.log.SetContext(&ctx).Infof("hello %s", req.Name)
	fmt.Printf("hello %s", req.Name)
	return
}

// SayHelloURL bm demo func.
func (s *Service) SayHelloURL(ctx context.Context, req *pb.HelloReq) (reply *pb.HelloResp, err error) {
	reply = &pb.HelloResp{
		Content: "hello " + req.Name,
		Id:      rand.Int63(),
	}
	s.log.SetContext(&ctx).Infof("hello url %s", req.Name)
	ss, _ := s.dao.GetDemo(ctx, "demo")
	// fmt.Printf("这里获取 redis 的数据 = %s", ss)
	s.log.Infof("这里获取 redis 的数据 = %s", ss)
	return
}

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, s.dao.Ping(ctx)
}

// Close close the resource.
func (s *Service) Close() {
}
