package main

import (
	"fmt"
	"reflect"
	"sync"
	"time"
	"unsafe"
)

// 线程安全的计数器
type Counter struct {
	mu    sync.Mutex
	count uint64
}

type Obj struct {
	b bool
	a int64
	c byte
	d int64
}
type Obj2 struct {
	a int64
	d int64
	c byte
	b bool
}

func (c *Counter) Incr() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}
func (c *Counter) Count() uint64 {
	return c.count
}

func main() {
	var o Obj = Obj{
		true, 5555, 'a', 88888888,
	}
	v := unsafe.Pointer(&o)
	ip := (*bool)(v)
	fmt.Printf("%v", &ip)
	typ := reflect.TypeOf(Obj{})
	fmt.Printf("Struct is %d bytes long\n", typ.Size())

	typ2 := reflect.TypeOf(Obj2{})
	fmt.Printf("Struct is %d bytes long\n", typ2.Size())

	counter()

}

func counter() {
	var counter Counter
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			time.Sleep(time.Second)
			counter.Incr()
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(counter.Count())
}

/*

 */
