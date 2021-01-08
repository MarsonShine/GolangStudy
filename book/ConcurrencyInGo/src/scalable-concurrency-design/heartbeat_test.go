package main

import (
	"testing"
	"time"
)

func DoWork(done <-chan interface{}, nums ...int) (<-chan interface{}, <-chan int) {
	heartbeat := make(chan interface{}, 1)
	intStream := make(chan int)
	go func() {
		defer close(heartbeat)
		defer close(intStream)
		time.Sleep(2 * time.Second) // 1 模拟日常网络延迟
		for _, n := range nums {

			select {
			case heartbeat <- struct{}{}:
			default:
			}
			select {
			case <-done:
				return
			case intStream <- n:
			}
		}
	}()
	return heartbeat, intStream
}

// 测试因为外部因素（网络异常）很难正确调试
func TestDoWork_GeneratesAllNumbers(t *testing.T) {
	done := make(chan interface{})
	defer close(done)
	intSlice := []int{0, 1, 2, 3, 5}
	_, results := DoWork(done, intSlice...)
	for i, expected := range intSlice {
		select {
		case r := <-results:
			if r != expected {
				t.Errorf(
					"index %v: expected %v, but received %v,", i,
					expected, r,
				)
			}
		case <-time.After(1 * time.Second): // 1 设置超时，防止死锁
			t.Fatal("test timed out")
		}
	}
}

func TestDoWork_GeneratesAllNumbers2(t *testing.T) {
	done := make(chan interface{})
	defer close(done)
	intSlice := []int{0, 1, 2, 3, 5}
	heartbeat, results := DoWork(done, intSlice...)
	<-heartbeat //1
	i := 0
	for r := range results {
		if expected := intSlice[i]; r != expected {
			t.Errorf("index %v: expected %v, but received %v,", i, expected, r)
		}
		i++
	}
}

// 基于间隔的心跳
func DoWorkBaseInterval(done <-chan interface{}, pulseInterval time.Duration, nums ...int) (<-chan interface{}, <-chan int) {
	heartbeat := make(chan interface{}, 1)
	intStream := make(chan int)
	go func() {
		defer close(heartbeat)
		defer close(intStream)
		time.Sleep(2 * time.Second)
		pulse := time.Tick(pulseInterval)
	numLoop: //2
		for _, n := range nums {
			for { //1
				select {
				case <-done:
					return
				case <-pulse:
					select {
					case heartbeat <- struct{}{}:
					default:
					}
				case intStream <- n:
					continue numLoop //3
				}
			}
		}
	}()
	return heartbeat, intStream
}

func TestDoWork_GeneratesAllNumbers3(t *testing.T) {
	done := make(chan interface{})
	defer close(done)
	intSlice := []int{0, 1, 2, 3, 5}
	const timeout = 2 * time.Second
	heartbeat, results := DoWorkBaseInterval(done, timeout/2, intSlice...)
	<-heartbeat //4
	i := 0
	for {
		select {
		case r, ok := <-results:
			if ok == false {
				return
			} else if expected := intSlice[i]; r != expected {
				t.Errorf(
					"index %v: expected %v, but received %v,", i,
					expected, r,
				)
			}
			i++
		case <-heartbeat: //5
		case <-time.After(timeout):
			t.Fatal("test timed out")
		}
	}
}
