package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	recipe "go.etcd.io/etcd/client/v3/experimental/recipes"
)

/*
Enter: 会被阻塞住，直到一共有 count（初始化这个栅栏的时候设定的值）个节点调用了 Enter，这 count 个被阻塞的节点才能继续执行
Leave: 节点调用 Leave 方法的时候，会被阻塞，直到有 count 个节点，都调用了 Leave 方法，这些节点才能继续执行
*/
var (
	addr        = flag.String("addr", "http://127.0.0.1:2379", "etcd addresses")
	barrierName = flag.String("name", "my-test-barrier", "barrier name")
	count       = flag.Int("c", 2, "")
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

	// 创建session
	s1, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s1.Close()

	// 获取、创建berrier
	b := recipe.NewDoubleBarrier(s1, *barrierName, *count)
	// 从命令行读取指令
	consoleScanner := bufio.NewScanner(os.Stdin)
	for consoleScanner.Scan() {
		action := consoleScanner.Text()
		items := strings.Split(action, " ")
		switch items[0] {
		case "enter": // 持有barrier
			b.Enter()
			fmt.Println("enter")
		case "leave":
			b.Leave()
			fmt.Println("leave")
		case "quit", "exit":
			return
		default:
			fmt.Println("unknown actions")
		}
	}
}
