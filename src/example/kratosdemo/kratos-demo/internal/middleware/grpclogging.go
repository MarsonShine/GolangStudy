package middleware

import (
	"context"
	"strconv"
	"time"

	"github.com/MSLibs/glogger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GrpcClientLogging() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		nctx := writeMetadataToContext(ctx)
		err := invoker(nctx, method, req, reply, cc, opts...)
		return err
	}
}

func GrpcServerLogging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		readFromMetadataToContext := func(cctx context.Context) context.Context {
			md, ok := metadata.FromIncomingContext(ctx)
			if ok {
				cctx = context.WithValue(cctx, "ip", md["sourceip"][0])
				cctx = context.WithValue(cctx, "platformId", md["platformid"][0])
				cctx = context.WithValue(cctx, "requestId", md["requestid"][0])
				cctx = context.WithValue(cctx, "userflag", md["userflag"][0])
				cctx = context.WithValue(cctx, "url", md["url"][0])
				cctx = context.WithValue(cctx, "serverip", md["serverip"][0])
				cctx = context.WithValue(cctx, "userAgent", md["useragent"][0])
				cctx = context.WithValue(cctx, "referer", md["referer"][0])
				cctx = context.WithValue(cctx, "serverip", md["serverip"][0])
				size, ok := strconv.ParseInt(md["size"][0], 10, 64)
				if ok == nil {
					cctx = context.WithValue(cctx, "size", size)
				}
				duration, ok := time.Parse(time.UnixDate, md["duration"][0])
				if ok == nil {
					cctx = context.WithValue(cctx, "duration", duration)
				}
			}
			return cctx
		}
		return handler(readFromMetadataToContext(ctx), req)
	}
}

func writeMetadataToContext(ctx context.Context) context.Context {
	md := metadata.Pairs("userflag", ctx.Value("userflag").(string),
		glogger.RequestID, ctx.Value(glogger.RequestID).(string),
		glogger.PlatformID, ctx.Value(glogger.PlatformID).(string),
		glogger.Duration, ctx.Value(glogger.Duration).(time.Time).Format(time.UnixDate),
		"referer", ctx.Value("referer").(string),
		"userAgent", ctx.Value("userAgent").(string),
		"size", strconv.FormatInt(ctx.Value("size").(int64), 10),
		"url", ctx.Value("url").(string),
		"sourceip", ctx.Value("sourceip").(string),
		"serverip", ctx.Value("serverip").(string))
	return metadata.NewOutgoingContext(ctx, md)
}
