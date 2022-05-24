package main

import (
	"fmt"
	"reflect"
	"time"
)

// channel 编排模式之fanin扇入模式
// 多个输入，一个是输出；即：多个channel输入，一个目的channel输出
// fanout扇出模式
// 一个输入，多个目标输出，常用于观察者模式；一个对象的状态发生变化，其它所有依赖于这个对象都会收到通知
func main() {
	start := time.Now()
	<-fanInReflect(
		sig(10*time.Second),
		sig(20*time.Second),
		sig(30*time.Second),
		sig(40*time.Second),
		sig(50*time.Second),
		sig(01*time.Minute),
	)
	// 执行到这里，说明orDone已经被close掉了
	fmt.Printf("done after %v", time.Since(start))
}

func fanInReflect(chans ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		// 构造SelectCase slice
		var cases []reflect.SelectCase
		for _, c := range chans {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}

		// 循环，从cases中选择一个可用的
		for len(cases) > 0 {
			i, v, ok := reflect.Select(cases)
			if !ok {
				// remove
				cases = append(cases[:i], cases[i+1:]...)
				continue
			}
			out <- v.Interface()
		}
	}()
	return out
}

func fanInRec(chans ...<-chan interface{}) <-chan interface{} {
	switch len(chans) {
	case 0:
		c := make(chan interface{})
		close(c)
		return c
	case 1:
		return chans[0]
	case 2:
		return mergeTwo(chans[0], chans[1])
	default:
		m := len(chans) / 2
		return mergeTwo(
			fanInRec(chans[:m]...),
			fanInRec(chans[m:]...))
	}
}
func mergeTwo(a, b <-chan interface{}) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		close(c)
		for a != nil || b != nil { //只要还有可读的chan
			select {
			case v, ok := <-a:
				if !ok { // a已关闭，设置为nil
					a = nil
					continue
				}
				c <- v
			case v, ok := <-b:
				if !ok {
					b = nil
					continue
				}
				c <- v
			}
		}
	}()
	return c
}

func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func fanOut(ch <-chan interface{}, out []chan interface{}, async bool) {
	go func() {
		defer func() {
			for i := 0; i < len(out); i++ {
				close(out[i])
			}
		}()

		for v := range ch {
			v := v // 避免延迟绑定
			for i := 0; i < len(out); i++ {
				i := i
				if async {
					go func() {
						out[i] <- v // 放入到输出chan中
					}()
				} else {
					out[i] <- v
				}
			}
		}
	}()
}

func fanOutReflect(ch <-chan interface{}, out []chan interface{}, async bool) {
	go func() {
		defer func() {
			for i := 0; i < len(out); i++ {
				close(out[i])
			}
		}()
		// 构造SelectCase slice
		var cases []reflect.SelectCase
		for v := range ch {
			for _, c := range out {
				cases = append(cases, reflect.SelectCase{
					Dir:  reflect.SelectSend,
					Chan: reflect.ValueOf(c),
					Send: reflect.ValueOf(v),
				})
			}
		}

		for len(cases) > 0 {
			if async {
				go func() {
					i, _, ok := reflect.Select(cases)
					if ok {
						// remove
						cases = append(cases[:i], cases[i+1:]...)
					}
				}()
			} else {
				i, _, ok := reflect.Select(cases)
				if ok {
					// remove
					cases = append(cases[:i], cases[i+1:]...)
				}
			}
		}
	}()
}
