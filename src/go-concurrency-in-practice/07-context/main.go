package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// ctx, cancel := context.WithCancel(context.Background())

	// go func() {
	// 	defer func() {
	// 		fmt.Println("goroutine exit")
	// 	}()

	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			return
	// 		default:
	// 			time.Sleep(time.Second)
	// 		}
	// 	}
	// }()

	// time.Sleep(time.Second)
	// cancel()
	// time.Sleep(2 * time.Second)

	cancellation()
}

// 父子context，父cancel，子是否也会被cancel
type newName string

func cancellation() {
	p := context.Background()
	ctx, cancel := context.WithCancel(p)

	child := context.WithValue(ctx, newName("name"), "marsonshine")
	go func() {
		for {
			select {
			case <-child.Done():
				fmt.Println("子context结束")
				return
			default:
				name := child.Value(newName("name"))
				fmt.Println("name:", name)
				time.Sleep(1 * time.Second)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-p.Done():
				fmt.Println("background context done!") // ctx 取消并不影响 background context 运行，因为它不是可取消的
				return
			default:
				fmt.Println("background context work well")
				time.Sleep(500 * time.Microsecond)
			}
		}
	}()
	// 父context结束
	time.Sleep(1 * time.Second)
	fmt.Println("父context取消")
	cancel()
	time.Sleep(2 * time.Second)

}
