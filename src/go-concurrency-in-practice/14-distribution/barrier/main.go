package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
)

/*
Hold 方法是创建一个 Barrier。如果 Barrier 已经创建好了，有节点调用它的 Wait 方法，就会被阻塞。
Release 方法是释放这个 Barrier，也就是打开栅栏。如果使用了这个方法，所有被阻塞的节点都会被放行，继续执行。
Wait 方法会阻塞当前的调用者，直到这个 Barrier 被 release。如果这个栅栏不存在，调用者不会被阻塞，而是会继续执行。
*/

var (
	addr        = flag.String("addr", "http://127.0.0.1:2379", "etcd addresses")
	barrierName = flag.String("name", "my-test-barrier", "barrier name")
)

func main() {
	flag.Parse()
	// 解析etcd地址
	endpoints := strings.Split(*addr, ",")
	// 创建etcd client
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// 获取、创建berrier
	b := recipe.NewBarrier(cli, *barrierName)
	// 从命令行读取指令
	consoleScanner := bufio.NewScanner(os.Stdin)
	for consoleScanner.Scan() {
		action := consoleScanner.Text()
		items := strings.Split(action, " ")
		switch items[0] {
		case "hold": // 持有barrier
			b.Hold()
			fmt.Println("hold")
		case "release":
			b.Release()
			fmt.Println("released")
		case "wait":
			fmt.Println("waiting")
			b.Wait()
			fmt.Println("after wait")
		case "quit", "exit":
			return
		default:
			fmt.Println("unknown actions")
		}
	}
}
