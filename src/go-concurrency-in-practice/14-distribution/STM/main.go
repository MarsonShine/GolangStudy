package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

/*
STM: 分布式事务
etcd 的事务基于CAS方式实现的，融合了Get，Put，Delete
执行过程: Txn().If(cond1. cond2, ...).Then(op1, op2, ...).Else(op1, op2, ...)


demo 案例：
下面这个例子创建了 5 个银行账号，然后随机选择一些账号两两转账。在转账的时候，要
把源账号一半的钱要转给目标账号。这个例子启动了 10 个 goroutine 去执行这些事务，
每个 goroutine 要完成 100 个事务。
*/

var (
	addr = flag.String("addr", "http://127.0.0.1:2379", "etcd addresses")
)

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

	// 设置5个账号，每个账户都有100元，总共有500元
	totalAccounts := 5
	for i := 0; i < totalAccounts; i++ {
		k := fmt.Sprintf("accts/%d", i)
		if _, err := cli.Put(context.TODO(), k, "100"); err != nil {
			log.Fatal(err)
		}
	}
	// STM的应用函数，主要的事务逻辑
	exchange := func(stm concurrency.STM) error {
		// 随机得到两个账户
		from, to := rand.Intn(totalAccounts), rand.Intn(totalAccounts)
		if from == to {
			return nil
		}
		// 读取账户的值
		fromK, toK := fmt.Sprintf("accts/%d", from), fmt.Sprintf("accts/%d", to)
		fromV, toV := stm.Get(fromK), stm.Get(toK)
		fromInt, toInt := 0, 0
		fmt.Sscanf(fromV, "%d", &fromInt)
		fmt.Sscanf(toV, "%d", &toInt)

		// 把源账户一半的钱转给目标账户
		xfer := fromInt / 2
		fromInt, toInt = fromInt-xfer, toInt+xfer
		// 写回交易后的值
		stm.Put(fromK, fmt.Sprintf("%d", fromInt))
		stm.Put(toK, fmt.Sprintf("%d", toInt))
		return nil
	}

	// 启动10个worker
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				if _, serr := concurrency.NewSTM(cli, exchange); serr != nil {
					log.Fatal(serr)
				}
			}
		}()
	}
	wg.Wait()

	// 检查数目是否一致
	sum := 0
	accts, err := cli.Get(context.TODO(), "accts/", clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}
	for _, kv := range accts.Kvs {
		v := 0
		fmt.Sscanf(string(kv.Value), "%d", &v)
		sum += v
		log.Printf("account %s: %d", kv.Key, v)
	}
	log.Println("account sum is", sum) // 总数
}

// fmt.Sscanf 用法：https://www.geeksforgeeks.org/fmt-sscanf-function-in-golang-with-examples/
