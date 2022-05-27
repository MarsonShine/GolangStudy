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

var (
	addr      = flag.String("addr", "http://127.0.0.1:2379", "etcd addresses")
	queueName = flag.String("name", "my-test-queue", "queue name")
)

// 利用etcd实现分布式队列
func main() {
	flag.Parse()

	// 解析etcd地址
	endpoints := strings.Split(*addr, ",")

	// 创建etcd的client
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// 创建、获取队列
	q := recipe.NewQueue(cli, *queueName)

	// 从命令行读取命令
	consoleScanner := bufio.NewScanner(os.Stdin)
	for consoleScanner.Scan() {
		action := consoleScanner.Text()
		items := strings.Split(action, " ")
		switch items[0] {
		case "push": // 加入队列
			if len(items) != 2 {
				fmt.Println("must set value to push")
				continue
			}
			q.Enqueue(items[1])
		case "pop":
			v, err := q.Dequeue()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(v)
		case "quit", "exit":
			return
		default:
			fmt.Println("unknown action")
		}
	}
}
