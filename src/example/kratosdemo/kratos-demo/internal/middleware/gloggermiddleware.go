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
		method:    ctx.Request.Method,
		url:       ctx.Request.URL.String(),
		referer:   ctx.Request.Header.Get("Referer"),
		userAgent: ctx.Request.Header.Get("User-Agent"),
	}
	ri.ip = requestGetRemoteAddress(ctx.Request)
	// this runs handler h and captures information about
	// HTTP request
	// m := httpsnoop.CaptureMetrics()
	ri.size = ctx.Request.ContentLength
	context := initLogContext(ctx.Request, ri)
	ctx.Request.WithContext(context)
}

type HttpRequestPayload struct {
	method    string
	url       string
	ip        string
	referer   string
	userAgent string
	size      int64
	duration  int64
}

var _ HttpRequestPayload = HttpRequestPayload{}

func initLogContext(r *http.Request, info *HttpRequestPayload) context.Context {
	start := time.Now()
	ctx := r.Context()
	ctx = context.WithValue(ctx, glogger.RequestID, r.Header.Get(glogger.RequestID))
	ctx = context.WithValue(ctx, glogger.UserFlag, r.Header.Get(glogger.UserFlag))
	ctx = context.WithValue(ctx, glogger.PlatformID, r.Header.Get(glogger.PlatformID))
	ctx = context.WithValue(ctx, "referer", r.Header.Get(glogger.PlatformID))
	ctx = context.WithValue(ctx, "userAgent", r.Header.Get(glogger.PlatformID))
	ctx = context.WithValue(ctx, "size", info.size)
	ctx = context.WithValue(ctx, "duration", start)
	return ctx
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
