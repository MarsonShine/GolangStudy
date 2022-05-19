package main

import (
	"fmt"
	"sync"
	"time"
)

// 具有环形依赖就会导致死锁
func deadlock() {
	var mu sync.RWMutex

	go func() {
		time.Sleep(200 * time.Microsecond)
		mu.Lock()
		fmt.Println("Lock")
		time.Sleep(100 * time.Microsecond)
		mu.Unlock()
		fmt.Println("Unlock")
	}()

	go func() {
		factorial(&mu, 10)
	}()

	select {}
}

func factorial(m *sync.RWMutex, n int) int {
	if n < 1 {
		return 0
	}
	fmt.Println("RLock")
	m.RLock()
	defer func() {
		fmt.Println("RUnlock")
		m.RUnlock()
	}()
	time.Sleep(100 * time.Microsecond)
	return factorial(m, n-1) * n
}
