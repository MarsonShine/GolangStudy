package main

import (
	"bytes"
	"fmt"
	"sync"
)

func main() {
	// orchannelExample()
	// pipelineExample()
	// pipelineExample2()
	// pipelineExample3()
	// fanoutinExample()
	teeExample()
}

func ppExample() {
	// 访问范围约束
	chanOwner := func() <-chan int {
		results := make(chan int, 5) //1 在 chanOwner 函数下定义通道，约束范围，只能在下面的匿名函数能访问通道并发送值
		go func() {
			defer close(results)
			for i := 0; i <= 5; i++ {
				results <- i
			}
		}()
		return results
	}
	consumer := func(results <-chan int) { //3
		for result := range results {
			fmt.Printf("Received: %d\n", result)
		}
		fmt.Println("Done receiving!")
	}
	results := chanOwner() //2
	consumer(results)
}

func ppExample2() {
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()
		var buff bytes.Buffer
		for _, b := range data {
			fmt.Fprintf(&buff, "%c", b)
		}
		fmt.Println(buff.String())
	}
	var wg sync.WaitGroup
	wg.Add(2)
	data := []byte("golang")
	go printData(&wg, data[:3]) // 1
	go printData(&wg, data[3:]) // 2
	wg.Wait()
}

// for {
// 	select {
// 	case <-done:
// 	  return
// 	  default:
// 	}
// 	// 执行非抢占任务
// }
func ppExample3() {

}
