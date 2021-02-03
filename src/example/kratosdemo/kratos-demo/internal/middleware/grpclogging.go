package middleware

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func GrpcClientLogging() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ri := &HttpRequestPayload{}

		fillRequestHeaders(ctx, ri)

		// call server handler

		startTime := time.Now()
		var remoteIP string
		if peerInfo, ok := peer.FromContext(ctx); ok {
			remoteIP = peerInfo.Addr.String()
		}
		ri.ip = remoteIP
		ri.startTime = startTime
		ri.url = method
		ri.method = method
		return nil
	}
}

func GrpcServerLogging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ri := &HttpRequestPayload{}
		fillRequestHeaders(ctx, ri)
		// call server handler

		startTime := time.Now()
		var remoteIP string
		if peerInfo, ok := peer.FromContext(ctx); ok {
			remoteIP = peerInfo.Addr.String()
		}
		ri.ip = remoteIP
		ri.startTime = startTime
		ri.url = info.FullMethod
		ri.method = info.FullMethod
		return handler(ctx, req)
	}
}

func fillRequestHeaders(ctx context.Context, payload *HttpRequestPayload) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return
	}
	requestId := md["requestId"]
	platformId := md["platformId"]
	userflag := md["userflag"]
	if len(requestId) != 0 {
		payload.requestId = requestId[0]
	}
	if len(platformId) != 0 {
		payload.platformId = platformId[0]
	}
	if len(userflag) != 0 {
		payload.userflag = userflag[0]
	}
}
