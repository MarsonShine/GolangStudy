package main

import (
	"fmt"
	"io"
	"net"
	"runtime"
	"syscall"
)

func main() {
	// 通信域：IPv4、IPv6、Unix域
	// Socket类型：Sock_DGRAM（UDP）、Sock_STREAM（TCP）、SOCK_SEQPACKET、SOCK_RAW
	// Socket使用的协议
	fd1, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		fmt.Printf("socket error: %v\n", err)
		return
	}
	defer syscall.Close(fd1)
	fmt.Printf("The file descriptor of socket：%d\n", fd1)

	// 网络编程 demo1
	network := "tcp"
	host := "google.cn"
	reqStrTpl := `HEAD / HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Host: %s
User-Agent: Dialer/%s
`
	// demo1
	network1 := network + "4"
	address1 := host + ":80"
	fmt.Printf("Dial %q with network %q ...\n", address1, network1)
	conn1, err := net.Dial(network1, address1)
	if err != nil {
		fmt.Printf("dial error: %v\n", err)
		return
	}
	defer conn1.Close()

	reqStr1 := fmt.Sprintf(reqStrTpl, host, runtime.Version())
	fmt.Printf("The request:\n%s\n", reqStr1)
	_, err = io.WriteString(conn1, reqStr1)
	if err != nil {
		fmt.Printf("write error: %v\n", err)
		return
	}
	fmt.Println()
}
