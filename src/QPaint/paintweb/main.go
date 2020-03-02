package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func newReverseProxy(baseURL string) *httputil.ReverseProxy {
	rpURL, _ := url.Parse(baseURL)
	return httputil.NewSingleHostReverseProxy(rpURL)
}

func handleDefault(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		http.ServeFile(w, req, "www/index.htm")
		return
	}
	req.URL.RawQuery = "" //跳过“?param”

}

var (
	apiReverseProxy = newReverseProxy("http://localhost:9999")
	wwwServer       = http.FileServer(http.Dir("www"))
)

func main() {
	// go paintserver.Main()
	// http.Handle("/api/", http.StripPrefix("/api/", apiReverseProxy))
	// http.HandleFunc("/", handleDefault)
	// http.ListenAndServe(":8888", nil)

	for i := 0; i <= 3; i++ {
		defer fmt.Println(i)
	}
	fmt.Println("==================")
	p := f()
	fmt.Println(p)
}

func f() (result int) {
	// defer func() {
	// 	result *= 7
	// }()
	return 6
}
