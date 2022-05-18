package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

var host string = "google.cn"
var domains = []string{
	"google.com",
	"google.com.hk",
	"google.cn",
	"golang.org",
	"golang.google.cn",
}

func main() {
	get1()
	get2()
}

func get1() {

	url1 := "http://" + host
	fmt.Printf("Send request to %q with method GET ...\n", url1)
	resp1, err := http.Get(url1)
	if err != nil {
		fmt.Printf("request sending error: %v\n", err)
	}
	defer resp1.Body.Close()
	line1 := resp1.Proto + " " + resp1.Status
	fmt.Printf("The first line of response:\n%s\n", line1)
}

func get2() {
	url2 := "http://golang." + host
	fmt.Printf("Send request to %q with method GET ...\n", url2)
	var httpClient http.Client
	resp2, err := httpClient.Get(url2)
	if err != nil {
		fmt.Printf("request sending error: %v\n", err)
	}
	defer resp2.Body.Close()
	line2 := resp2.Proto + " " + resp2.Status
	fmt.Printf("The first line of response:\n%s\n", line2)
}

func transport() {
	myTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:       15 * time.Second,
			KeepAlive:     15 * time.Second,
			FallbackDelay: 0,
		}).DialContext,
		MaxConnsPerHost:       2,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   2,
		IdleConnTimeout:       30 * time.Second,
		ResponseHeaderTimeout: 0,
		ExpectContinueTimeout: 1 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
	}
	// 你可以改变myClient中的各个字段的值，
	// 并观察后面的输出会有什么不同。
	myClient := http.Client{
		Transport: myTransport,
		Timeout:   20 * time.Second,
	}

	var wg sync.WaitGroup
	wg.Add(len(domains))
	for _, domain := range domains {
		go func(domain string) {
			var logBuf strings.Builder
			var diff time.Duration
			defer func() {
				logBuf.WriteString(
					fmt.Sprintf("(elapsed time: %s)\n", diff))
				fmt.Println(logBuf.String())
				wg.Done()
			}()
			url := "https://" + domain
			logBuf.WriteString(
				fmt.Sprintf("Send request to %q with method GET ...\n", url))
			t1 := time.Now()
			resp, err := myClient.Get(url)
			diff = time.Now().Sub(t1)
			if err != nil {
				logBuf.WriteString(
					fmt.Sprintf("request sending error: %v\n", err))
				return
			}
			defer resp.Body.Close()
			line2 := resp.Proto + " " + resp.Status
			logBuf.WriteString(
				fmt.Sprintf("The first line of response:\n%s\n", line2))
		}(domain)
	}
	wg.Wait()
}

func server() {
	var wg sync.WaitGroup
	wg.Add(2)

	// 示例1。
	go startServer1(&wg)

	// 示例2。
	go startServer2(&wg)

	wg.Wait()
}

func startServer1(wg *sync.WaitGroup) {
	defer wg.Done()
	var httpServer1 http.Server
	httpServer1.Addr = "127.0.0.1:8080"
	// 由于我们没有定制handler，所以这个网络服务对任何请求都只会响应404。
	if err := httpServer1.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("HTTP server 1 closed.")
		} else {
			log.Printf("HTTP server 1 error: %v\n", err)
		}
	}
}

func startServer2(wg *sync.WaitGroup) {
	defer wg.Done()
	mux1 := http.NewServeMux()
	mux1.HandleFunc("/hi", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/hi" {
			http.NotFound(w, req)
			return
		}
		name := req.FormValue("name")
		if name == "" {
			fmt.Fprint(w, "Welcome!")
		} else {
			fmt.Fprintf(w, "Welcome, %s!", name)
		}
	})
	httpServer2 := http.Server{
		Addr:    "127.0.0.1:8081",
		Handler: mux1,
	}
	if err := httpServer2.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("HTTP server 2 closed.")
		} else {
			log.Printf("HTTP server 2 error: %v\n", err)
		}
	}
}
