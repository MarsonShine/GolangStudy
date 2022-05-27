package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vardius/gollback"
	"golang.org/x/sync/errgroup"
)

/*
分组并发原语：分组执行一批相同的或类似的任务
使用场景：需要将一个通用的父任务拆成几个小任务并发执行的场景

将一个大的任务拆成几个小任务并发执行，可以有效地提高程序的并发度
*/

func main() {
	returnFirstError()
}

// 启动三个任务，其中任务2会返回失败，其它两个任务执行成功
// Wait会等到这三个任务一起完成（无论是否失败）
func returnFirstError() {
	var g errgroup.Group

	// 启动第一个任务
	g.Go(func() error {
		time.Sleep(5 * time.Second)
		fmt.Println("exec #1")
		return nil
	})
	// 启动第二个任务
	g.Go(func() error {
		time.Sleep(10 * time.Second)
		fmt.Println("exec #2")
		return errors.New("failed to exec #2")
	})
	// 启动第三个任务
	g.Go(func() error {
		time.Sleep(15 * time.Second)
		fmt.Println("exec #4")
		return nil
	})
	// 等待三个任务全部完成
	if err := g.Wait(); err == nil {
		fmt.Println("Successfully exec all")
	} else {
		fmt.Println("failed:", err)
	}
}

func returnAllError() {
	var g errgroup.Group
	var result = make([]error, 3)

	// 启动第一个任务
	g.Go(func() error {
		time.Sleep(5 * time.Second)
		fmt.Println("exec #1")
		result[0] = nil
		return nil
	})

	// 启动第二个任务
	g.Go(func() error {
		time.Sleep(10 * time.Second)
		fmt.Println("exec #2")
		result[1] = errors.New("failed to exec #2")
		return result[1]
	})

	// 启动第三个任务
	g.Go(func() error {
		time.Sleep(15 * time.Second)
		fmt.Println("exec #3")
		result[2] = nil
		return nil
	})
	// 等待三个任务全部完成
	if err := g.Wait(); err == nil {
		fmt.Println("Successfully exec all")
	} else {
		fmt.Println("failed:", err)
	}
}

// 同时处理多个任务，并返回多个错误信息（如果有的话）
func handleMultipleTask() {
	rs, errs := gollback.All(
		context.Background(),
		func(ctx context.Context) (interface{}, error) {
			time.Sleep(3 * time.Second)
			return 1, nil
		},
		func(ctx context.Context) (interface{}, error) {
			return nil, errors.New("failed")
		},
		func(ctx context.Context) (interface{}, error) {
			return 3, nil
		},
	)
	fmt.Println(rs...)
	fmt.Println(errs)

	// 只要有任意一个函数有错误，就立马返回
	// gollback.Race(...)
	// 重试，再执行一个任务时，如果执行失败会尝试一定的次数
	// gollback.Retry()
}
