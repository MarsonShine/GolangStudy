package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/MSLibs/glogger"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
)

func UseGLogger(e *bm.Engine) {
	e.Use(GLoggerMiddleware{})
}

type GLoggerMiddleware struct{}

var _ GLoggerMiddleware = GLoggerMiddleware{}

func (g GLoggerMiddleware) ServeHTTP(ctx *bm.Context) {
	ri := &HttpRequestPayload{
		method:     ctx.Request.Method,
		url:        ctx.Request.URL.String(),
		referer:    ctx.Request.Header.Get("Referer"),
		userAgent:  ctx.Request.Header.Get("User-Agent"),
		requestId:  ctx.Request.Header.Get(glogger.RequestID),
		userflag:   ctx.Request.Header.Get(glogger.UserFlag),
		platformId: ctx.Request.Header.Get(glogger.PlatformID),
	}
	ri.ip = requestGetRemoteAddress(ctx.Request)
	// this runs handler h and captures information about
	// HTTP request
	// m := httpsnoop.CaptureMetrics()
	ri.size = ctx.Request.ContentLength
	initLogContext(ctx, ri)
}

type HttpRequestPayload struct {
	method     string
	url        string
	ip         string
	referer    string
	userAgent  string
	requestId  string
	platformId string
	userflag   string
	size       int64
	duration   int64
}

var _ HttpRequestPayload = HttpRequestPayload{}

func initLogContext(ctx *bm.Context, info *HttpRequestPayload) {
	start := time.Now()
	// requestId := ctx.Request.Header.Get(glogger.RequestID)
	// userflag := ctx.Request.Header.Get(glogger.UserFlag)
	// platformId := ctx.Request.Header.Get(glogger.PlatformID)
	ctx.Context = context.WithValue(ctx.Context, glogger.RequestID, info.requestId)
	ctx.Context = context.WithValue(ctx.Context, glogger.UserFlag, info.userflag)
	ctx.Context = context.WithValue(ctx.Context, glogger.PlatformID, info.platformId)
	ctx.Context = context.WithValue(ctx.Context, "referer", info.referer)
	ctx.Context = context.WithValue(ctx.Context, "userAgent", info.userAgent)
	ctx.Context = context.WithValue(ctx.Context, "size", info.size)
	ctx.Context = context.WithValue(ctx.Context, "duration", start)
	// ctx.Set(glogger.RequestID, info.requestId)
	// ctx.Set(glogger.UserFlag, info.userflag)
	// ctx.Set(glogger.PlatformID, info.platformId)
	// ctx.Set("referer", info.referer)
	// ctx.Set("userAgent", info.userAgent)
	// ctx.Set("size", info.size)
	// ctx.Set("duration", start)
}

func requestGetRemoteAddress(r *http.Request) string {
	hdr := r.Header
	hdrRealIP := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	if hdrRealIP == "" && hdrForwardedFor == "" {
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}
	if hdrForwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		// TODO: should return first non-local address
		return parts[0]
	}
	return hdrRealIP
}

func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}
