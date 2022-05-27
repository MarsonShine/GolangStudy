package cyclicbarrier

import (
	"math/rand"
	"sort"
	"sync"
	"testing"
	"time"
)

func TestWaterFactory(t *testing.T) {
	//用来存放水分子结果的channel
	var ch chan string
	releaseHydrogen := func() {
		ch <- "H"
	}
	releaseOxygen := func() {
		ch <- "O"
	}
	// 300个原子，300个goroutine,每个goroutine并发的产生一个原子
	var N = 100
	ch = make(chan string, N*3)
	h2o := New()
	// 用来等待所有的goroutine完成
	var wg sync.WaitGroup
	wg.Add(N * 3)
	// 200个氢原子goroutine
	for i := 0; i < 2*N; i++ {
		go func() {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			h2o.hydrogen(releaseHydrogen)
			wg.Done()
		}()
	}
	// 100个氧原子goroutine
	for i := 0; i < N; i++ {
		go func() {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			h2o.oxygen(releaseOxygen)
			wg.Done()
		}()
	}
	//等待所有的goroutine执行完
	wg.Wait()
	// 结果中肯定是300个原子
	if len(ch) != N*3 {
		t.Fatalf("expect %d atom but got %d", N*3, len(ch))
	}
	// 每三个原子一组，分别进行检查。要求这一组原子中必须包含两个氢原子和一个氧原子，这样才能
	var s = make([]string, 3)
	for i := 0; i < N; i++ {
		s[0] = <-ch
		s[1] = <-ch
		s[2] = <-ch
		sort.Strings(s)
		water := s[0] + s[1] + s[2]
		if water != "HHO" {
			t.Fatalf("expect a water molecule but got %s", water)
		}
	}
}
