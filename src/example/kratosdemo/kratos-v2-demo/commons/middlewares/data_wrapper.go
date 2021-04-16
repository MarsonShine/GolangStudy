package middlewares

import (
	"context"
	"reflect"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/runtime/protoimpl"

	"github.com/go-kratos/kratos/v2/middleware"
)

type Option func(*options)

type options struct {
	logger log.Logger
}

func Server(opts ...Option) middleware.Middleware {
	options := options{
		logger: log.DefaultLogger,
	}
	for _, o := range opts {
		o(&options)
	}
	log := log.NewHelper("middleware/datawrapper", options.logger)

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			reply, err := handler(ctx, req)
			res := DataResponseWrapper{}
			if err != nil {
				res.message = err.Error()
				res.success = false
				res.code = 500
			} else {
				res.success = true
				// 提取指针对象
				v := reflect.ValueOf(reply)
				if v.Kind() == reflect.Ptr {
					v = v.Elem()
				}
				var realV interface{}
				for i := 0; i < v.NumField(); i++ {
					var field = v.Type().Field(i)
					if field.Type == reflect.TypeOf(protoimpl.MessageState{}) ||
						field.Name == "unknownFields" ||
						field.Name == "sizeCache" {
						continue
					}
					realV = v.Field(i).Interface()
				}
				res.data = realV
			}
			log.Debugf("data wrapper: origin = %v, now = %v", reply, res)
			return &res, nil
		}
	}
}

type DataResponseWrapper struct {
	success bool
	message string
	code    int
	data    interface{}
}
