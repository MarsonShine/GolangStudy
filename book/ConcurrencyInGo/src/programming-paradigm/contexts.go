package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// 在并发程序中，可能会存在 goroutine 出现问题如超时，取消，出错等
// 虽然之前我们可以通过 done 来通知其它任务取消，但也只是仅仅通知而已
// 如果要给出取消的原因给其它任务，目前是做不到的
// go 1.7 引入 context 来解决这个问题
// context 主要有两个目的
// 1. 提供取消操作
// 2. 提供用于通过调用传输请求附加数据的数据包

func contextExample() {
	var wg sync.WaitGroup
	done := make(chan interface{})
	defer close(done)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printGreeting(done); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewall(done); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}()
	wg.Wait()
}

// 通过 context.Context 来实现取消操作
func contextExample2() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background()) // 1 建立一个新的 context，传递给下级函数
	defer cancel()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printGreetingWithContext(ctx); err != nil {
			fmt.Printf("cannot print greeting: %v\n", err)
			cancel() // 2 发生错误，取消操作
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewellWithContext(ctx); err != nil {
			fmt.Printf("cannot print farewell: %v\n", err)
		}
	}()
	wg.Wait()
}

// localeWithContext 函数已知需要1分钟的时间，我们可以使用deadline
func contextExample3() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printGreetingWithContext(ctx); err != nil {
			fmt.Printf("cannot print greeting: %v\n", err)
			cancel()
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewellWithContext(ctx); err != nil {
			fmt.Printf("cannot print farewell: %v\n", err)
		}
	}()
	wg.Wait()
}

// 可以给上下文请求添加数据
func contextExample4() {
	ProcessRequest("username", "marsonshine")
}

func ProcessRequest(key string, value string) {
	ctx := context.WithValue(context.Background(), key, value)
	ctx = context.WithValue(ctx, "other", "summerzhu")
	HandleResponse(ctx)
}

func HandleResponse(ctx context.Context) {
	fmt.Printf("handling response for %v (%v)",
		ctx.Value("username"),
		ctx.Value("other"))
}

func printGreetingWithContext(ctx context.Context) error {
	greeting, err := genGreetingWithContext(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", greeting)
	return nil
}
func printFarewellWithContext(ctx context.Context) error {
	farewell, err := genFarewellWithContext(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", farewell)
	return nil
}

func genGreetingWithContext(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second) //3 context.WithTimeout 包装新的 context，并1秒后自动返回从而取消操作被调用
	defer cancel()
	switch locale, err := localeWithContext(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func genFarewellWithContext(ctx context.Context) (string, error) {
	switch locale, err := localeWithContext(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}
	return "", fmt.Errorf("unsupported locale")
}
func localeWithContext(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err() // 4: 返回context错误的原因，会一直冒泡到函数的调用处，这会导致 2 注释的取消操作被调用
	case <-time.After(1 * time.Minute):
	}
	return "EN/US", nil
}

func localeWithContextDeadline(ctx context.Context) (string, error) {
	if deadline, ok := ctx.Deadline(); ok { // 1,检查context提供的deadline，如果提供了值并且程序时间超过这个时间则就会返回预设的 context.DeadlineExceeded 错误
		if deadline.Sub(time.Now().Add(1*time.Minute)) <= 0 {
			return "", context.DeadlineExceeded
		}
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(1 * time.Minute):
	}
	return "EN/US", nil
}

func printFarewall(done <-chan interface{}) error {
	farewell, err := genFarewell(done)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", farewell)
	return nil
}

func genFarewell(done <-chan interface{}) (string, error) {
	switch locale, err := locale(done); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func printGreeting(done <-chan interface{}) error {
	greeting, err := genGreeting(done)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", greeting)
	return nil
}

func genGreeting(done <-chan interface{}) (string, error) {
	switch locale, err := locale(done); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func locale(done <-chan interface{}) (string, error) {
	select {
	case <-done:
		return "", fmt.Errorf("canceled")
	case <-time.After(5 * time.Second):
	}
	return "EN/US", nil
}
