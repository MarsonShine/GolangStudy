package middleware

import (
	"context"
	"kratos-demo/utils"
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
	serverip   string
	referer    string
	userAgent  string
	requestId  string
	platformId string
	userflag   string
	size       int64
	duration   int64
	startTime  time.Time
}

var _ HttpRequestPayload = HttpRequestPayload{}

func initLogContext(bctx *bm.Context, info *HttpRequestPayload) {
	start := time.Now()
	ctx := context.WithValue(bctx.Context, glogger.RequestID, info.requestId)
	ctx = context.WithValue(ctx, "userflag", info.userflag)
	ctx = context.WithValue(ctx, glogger.PlatformID, info.platformId)
	ctx = context.WithValue(ctx, "referer", info.referer)
	ctx = context.WithValue(ctx, "userAgent", info.userAgent)
	ctx = context.WithValue(ctx, "size", info.size)
	ctx = context.WithValue(ctx, "duration", start)
	ctx = context.WithValue(ctx, "url", bctx.Request.URL.String())
	ctx = context.WithValue(ctx, "sourceip", requestGetRemoteAddress(bctx.Request))
	if serverip, err := utils.ExternalIP(); err == nil {
		ctx = context.WithValue(ctx, "serverip", serverip)
	}
	bctx.Context = ctx
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
