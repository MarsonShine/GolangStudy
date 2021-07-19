package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
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
	// poll_bad(10)                               // 是几秒钟还是几毫秒，不清楚
	// poll_good(time.Duration(10) * time.Second) // 能清晰的知道是秒
	reflectExample()
	// editReflectValue()
	editAddressReflectValue()
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

// 错误处理最佳实践
// 1 对于只是简单的错误信息，则用 error.New 即可
// 2 对于要在客户端处理错误时，则需要用自定义错误类型，实现 Error() 方法
// 3 要传递给下游函数的错误，则需要通过错误包裹处理
// 4 其它用 fmt.Errorf() 即可
func Open() error {
	return errors.New("could not open")
}

func error_bad() {
	if err := Open(); err != nil {
		if err.Error() == "could not open" {
			// handle
		} else {
			panic("unkonw error")
		}
	}
}

// 自定义错误
var ErrCouldNotOpen = errors.New("could not open")

func open_good() error {
	return ErrCouldNotOpen
}
func error_good() {
	if err := open_good(); err != nil {
		if err == ErrCouldNotOpen {
			// handle
		} else {
			panic("unkonw error")
		}
	}
}

// 2 如果您有可能需要客户端检测的错误，并且想向其中添加更多信息（例如，它不是静态字符串），则应使用自定义类型。
func open_bad(file string) error {
	return fmt.Errorf("file %q not found", file)
}
func error_bad2() {
	if err := open_bad("testfile.txt"); err != nil {
		if strings.Contains(err.Error(), "not found") {
			// handle
		} else {
			panic("unknown error")
		}
	}
}

// 自定义类型
type errNotFound struct {
	file string
}

func (e errNotFound) Error() string {
	return fmt.Sprintf("file %q not found", e.file)
}

// 一般还暴露一个api检查该错误
func IsNotFoundError(err error) bool {
	_, ok := err.(errNotFound)
	return ok
}

func open_good2(file string) error {
	return errNotFound{file: file}
}
func error_good2() {
	if err := open_good2("testfile.txt"); err != nil {
		if IsNotFoundError(err) {
			// handle
		}
		if _, ok := err.(errNotFound); ok {
			// handle
		} else {
			panic("unknown error")
		}
	}
}

// 4 错误包装，一般用于传递错误信息给下层函数
// 4.1 如果没有要添加的其他上下文，并且您想要维护原始错误类型，则返回原始错误。
// 4.2 添加上下文，使用 "pkg/errors".Wrap 以便错误消息提供更多上下文 ,"pkg/errors".Cause 可用于提取原始错误。
// 4.3 如果调用者不需要检测或处理的特定错误情况，使用 fmt.Errorf。

// 在程序中不要抛出 panic
func panic_bad() {
	run := func(args []string) {
		if len(args) == 0 {
			panic("an argument is required")
		}
		// ...
	}

	run(os.Args[1:])
}

// 要返回错误
func panic_never() error {
	run := func(args []string) error {
		if len(args) == 0 {
			return errors.New("an argument is required")
		}
		// ...
		return nil
	}
	var err error
	if err = run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return err
}
