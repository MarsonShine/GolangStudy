package main

import "time"

// HTTPReqInfo describes info about HTTP request
type HTTPReqInfo struct {
	// GET POST
	method  string
	url     string
	referer string
	ipaddr  string
	// 200, 404
	code int
	// 请求包大小
	size int64
	// 执行时间
	duration  time.Duration
	userAgent string
}
