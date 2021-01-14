package main

import (
	"net/http"
	"time"
)

type F interface {
	f()
}

type S1 struct{}

func (s S1) f() {}

type S2 struct{}

func (s S2) f() {}

func main() {
	// var f1 F = S1{}
	// var f2 F = &S2{}
	poll_bad(10)                               // 是几秒钟还是几毫秒，不清楚
	poll_good(time.Duration(10) * time.Second) // 能清晰的知道是秒
}

type Handler struct {
}

func (h *Handler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {

}

var _ http.Handler = (*Handler)(nil) // 会及时报错

// 枚举从 1 开始
type Operation int

const (
	Add Operation = iota + 1
	Subtract
	Multipy
)

// 时间的操作都要用 time 包里的，以表示准确
func isActive_Bad(now, start, stop int) bool {
	return start <= now && now < stop
}

func isActive_Good(now, start, stop time.Time) bool {
	return (start.Before(now) || start.Equal(now) && now.Before(stop))
}

// 使用 time.Duration 表示时间段
func poll_bad(delay int) {
	for {
		time.Sleep(time.Duration(delay) * time.Microsecond)
	}
}

func poll_good(delay time.Duration) {
	for {
		time.Sleep(delay)
	}
}
