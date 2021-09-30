package main

import (
	"fmt"
	"net"
	"net/url"
)

func main() {
	// 解析 url
	// url 包含 schema、认证信息、主机名、端口、路径、查询参数以及查询片段
	s := "postgres://user:pass@host.com:5432/path?k=v#f"

	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(u.Scheme)
	fmt.Println(u.User)
	fmt.Println(u.User.Username())
	fmt.Println(u.User.Password())

	fmt.Println(u.Host)
	host, port, _ := net.SplitHostPort(u.Host)
	fmt.Println(host)
	fmt.Println(port)

	fmt.Println("RawQeury:", u.Path)
	fmt.Println("RawQeury:", u.Fragment)

	fmt.Println("RawQeury:", u.RawQuery)
	m, _ := url.ParseQuery(u.RawQuery)
	fmt.Println(m)
	fmt.Println(m["k"][0])
}
