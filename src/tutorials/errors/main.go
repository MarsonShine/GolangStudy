package errors

// https://go.dev/blog/error-handling-and-go
import (
	"errors"
	"net/http"
)

func init() {
	http.HandleFunc("/index", viewRecord)
	// 第二版
	http.Handle("/index", appHandler(viewRecord2))
	// 版本三，可以返回自定义错误类型
}

func viewRecord(w http.ResponseWriter, r *http.Request) {
	key := dbGetKey()
	if err := dbGet(key); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// ...
	// 如果有更多其它操作，那么返回的错误重复率会越来越高
}

// 为了减少代码的重复，我们可以定义自己的 HTTP 处理程序，并返回 error

func dbGetKey() string {
	return "key from database"
}

func dbGet(key string) error {
	if key == "" {
		return errors.New("key is null")
	}
	return nil
}
