package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

var status int64

func condSync() {
	c := sync.NewCond(&sync.Mutex{})
	for i := 0; i < 10; i++ {
		go listen(c)
	}
	time.Sleep(5 * time.Second)
	go broadcast(c)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}

// broadcast signals the condition variable and wakes up all goroutines waiting on it.
//
// It takes a pointer to a sync.Cond object as a parameter.
// This function updates the value of the global variable "status" to 1 using the atomic.StoreInt64 function.
// Then, it calls the c.Broadcast() method to signal the condition variable and wake up all waiting goroutines.
func broadcast(c *sync.Cond) {
	// c.L.Lock()
	atomic.StoreInt64(&status, 1)
	c.Broadcast()
	// c.L.Unlock()
}

// listen waits until the condition variable is signaled.
//
// It takes a pointer to a sync.Cond object as a parameter.
// This function first locks the associated mutex using c.L.Lock().
// Then, it enters a loop and checks the value of the global variable "status".
// The loop continues until the value of "status" becomes 1, indicating that the condition has been met.
// Inside the loop, it calls c.Wait() to wait for a signal on the condition variable.
// Once the condition is met, it prints "listen" to the console.
// Finally, it unlocks the mutex using c.L.Unlock().
func listen(c *sync.Cond) {
	c.L.Lock()
	for atomic.LoadInt64(&status) != 1 {
		c.Wait()
	}
	fmt.Println("listen")
	c.L.Unlock()
}
