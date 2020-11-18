package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func main() {
	// doc, err := html.Parse(os.Stdin)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
	// 	os.Exit(1)
	// }
	// for _, link := range visit(nil, doc) {
	// 	fmt.Println(link)
	// }

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: wait url\n")
		os.Exit(1)
	}
	url := os.Args[1]
	//!+main
	// (In function main.)
	if err := waitForServer(url); err != nil {
		fmt.Fprintf(os.Stderr, "Site is down: %v\n", err)
		os.Exit(1) // 捕捉到错误，主动结束进程
	}
	//!-main
	//函数值
	functionValue()
	// 文件读取时，读到最后没有内容读取会报 EOF 错，这个时候不要当作程序错误了，要特殊处理
	in := bufio.NewReader(os.Stdin)
	for {
		r, _, err := in.ReadRune()
		if err == io.EOF {
			break // 表示读取整个内容成功
		}
		if err != nil {
			_ = fmt.Errorf("read failed:%v", err)
		}
		fmt.Print(r)
	}
}

// 递归
func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}

// 错误处理机制，如果错误是偶然的，那么最好选择重试，限制重试次数，时间间隔
func waitForServer(url string) error {
	const timeout = 1 * time.Minute
	deadline := time.Now().Add(timeout)
	// time.now() 的时间如果在 deadline 之前（早），返回 true，否则返回 false
	for tries := 0; time.Now().Before(deadline); tries++ {
		_, err := http.Head(url)
		if err == nil {
			return nil
		}
		log.Printf("server not responding (%s);retrying…", err)
		time.Sleep(time.Second << uint(tries))
	}
	return fmt.Errorf("server %s failed to respond after %s", url, timeout)
}

func functionValue() {
	f := square       // 相当于委托，将函数传递给一个变量
	fmt.Println(f(3)) // "9"
	f = negative
	fmt.Println(f(3))     // "-3"
	fmt.Printf("%T\n", f) // "func(int) int"
}

func square(n int) int     { return n * n }
func negative(n int) int   { return -n }
func product(m, n int) int { return m * n }

// 匿名函数
// 就是在关键字 func 后面没有函数名，这就是匿名函数
func anonyFunction() {
	strings.Map(func(r rune) rune { return 1 + 1 }, "HAL-9000")
	f := squares()
	fmt.Println(f()) // "1"
	fmt.Println(f()) // "4"
	fmt.Println(f()) // "9"
	fmt.Println(f()) // "16"
}

// 闭包
func squares() func() int {
	var x int
	return func() int {
		x++
		return x * x
	}
}
