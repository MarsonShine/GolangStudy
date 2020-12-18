package httpmiddleware

import (
	"fmt"
	"net/http"
)

func main() {
	mux := &http.ServeMux{}
	mux.HandleFunc("/", indexHandler)
	// 注册中间件
	var handler http.Handler = mux
	handler = RegisterHandler(handler)
	srv := &http.Server{
		Handler: handler,
	}
	srv.Addr = ":5000"
	srv.ListenAndServe()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("hello world")
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK) // 200
	w.Write([]byte(msg))
}

func RegisterHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		println("每次请求都会执行这个中间件")
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
