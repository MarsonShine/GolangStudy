package main

import "sync"

var lock sync.Mutex

func runtime_mutex() {
	lock.Lock()
	defer lock.Unlock()
}
