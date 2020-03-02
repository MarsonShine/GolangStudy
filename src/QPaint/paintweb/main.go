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
var testdata *struct {
	a *[7]int
}

func main() {
	// go paintserver.Main()
	// http.Handle("/api/", http.StripPrefix("/api/", apiReverseProxy))
	// http.HandleFunc("/", handleDefault)
	// http.ListenAndServe(":8888", nil)

	// for i := 0; i <= 3; i++ {
	// 	defer fmt.Println(i)
	// }
	// fmt.Println("==================")
	// p := f()
	// fmt.Println(p)

	var pp = P{
		"id": 11,
	}
	if val, ok := pp["ids"]; ok {
		fmt.Println(val)
	} else {
		fmt.Println("val 不存在")
	}

	// for 1 < 2 {
	// 	fmt.Println("1 < 2") // 死循环
	// }
	fmt.Println("==================for range ===================")
	for i, _ := range *testdata.a {
		fmt.Print(i)
		fmt.Print(" ")
		fmt.Print(len(*testdata.a))
		fmt.Print(" ")
	}

	var a [10]string
	for i, s := range a {
		fmt.Print(i)
		fmt.Print(" ")
		fmt.Print(s)
		fmt.Print(" ")
	}

	var key string
	var val interface{} // element type of m is assignable to val
	m := map[string]int{"mon": 0, "tue": 1, "wed": 2, "thu": 3, "fri": 4, "sat": 5, "sun": 6}
	for key, val = range m {
		fmt.Print("key=" + key)
		fmt.Print(" ")
		fmt.Print(val) // interface{}
		fmt.Print(" ")
	}

	// key == last map key encountered in iteration
	// val == map[key]
}

func f() (result int) {
	// defer func() {
	// 	result *= 7
	// }()
	return 6
}

type P map[string]interface{}
