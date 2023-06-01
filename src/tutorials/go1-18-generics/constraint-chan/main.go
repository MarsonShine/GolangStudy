package main

import (
	"context"
)

// https://github.com/golang/go/commit/070951c5dcc47c9cff2ad4c1ac6170a4060a4d0c
// 取消了该特性
func makeChan[T chan E, E ~int](ctx context.Context, arr []T) T {
	ch := make(T)
	go func() {
		defer close(ch)
		for _, v := range arr {
			select {
			case <-ctx.Done():
				return
			default:
			}
			ch <- v
		}
	}()
	return ch
}

func main() {

}
