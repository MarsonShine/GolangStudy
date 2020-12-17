package main

// TODO 自定义字段相关知识：https://github.com/uber-go/zap/issues/693
import (
	"context"
	"fmt"
	"golog/src/fzlog"
	"net/http"
	"strings"
	"time"
)

func initLogServer() {
	fzlog.CreateLog()
}

func makeHttpServer() *http.Server {
	mux := &http.ServeMux{}
	mux.HandleFunc("/", handleIndex)
	// ... potentially add more handlers

	var handler http.Handler = mux
	// wrap mux with our logger. this will
	handler = logRequestHandler(handler)
	// ... potentially add more middleware handlers

	srv := &http.Server{
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  120 * time.Second, // introduced in Go 1.8
		Handler:      handler,
	}
	return srv
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	msg := fmt.Sprintf("You've called url %s", r.URL.String())
	fzlog.WithContext(&ctx).Info("这是一个测试demo")
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK) // 200
	w.Write([]byte(msg))
}

func logRequestHandler(next http.Handler) http.Handler {
	initLogServer()
	fn := func(w http.ResponseWriter, r *http.Request) {
		ri := &HTTPReqInfo{
			method:    r.Method,
			url:       r.URL.String(),
			referer:   r.Header.Get("Referer"),
			userAgent: r.Header.Get("User-Agent"),
		}

		ri.ipaddr = requestGetRemoteAddress(r)

		// this runs handler h and captures information about
		// HTTP request
		// m := httpsnoop.CaptureMetrics()
		ri.size = r.ContentLength

		// ri.code = m.Code
		// ri.size = m.Written
		// ri.duration = m.Duration
		ctx := initLogContext(r, ri)
		// next
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func initLogContext(r *http.Request, info *HTTPReqInfo) context.Context {
	start := time.Now()
	ctx := r.Context()
	ctx = context.WithValue(ctx, fzlog.RequestID, r.Header.Get(fzlog.RequestID))
	ctx = context.WithValue(ctx, fzlog.UserFlag, r.Header.Get(fzlog.UserFlag))
	ctx = context.WithValue(ctx, fzlog.PlatformID, r.Header.Get(fzlog.PlatformID))
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

func skipHTTPRequestLogging(ri *HTTPReqInfo) bool {
	// we always want to know about failures and other
	// non-200 responses
	if ri.code != 200 {
		return false
	}

	// we want to know about slow requests.
	// 100 ms threshold is somewhat arbitrary
	if ri.duration > 100*time.Millisecond {
		return false
	}

	// this is linked from every page
	if ri.url == "/favicon.png" {
		return true
	}

	if ri.url == "/favicon.ico" {
		return true
	}

	if strings.HasSuffix(ri.url, ".css") {
		return true
	}
	return false
}
