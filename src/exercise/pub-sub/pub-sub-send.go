package main

import (
	"fmt"
	"sync"
	"time"
)

type PubSub struct {
	mutex  sync.Mutex
	subs   map[string][]chan string
	closed bool
}

func NewPubSub() *PubSub {
	return &PubSub{
		subs: make(map[string][]chan string),
	}
}

func (ps *PubSub) Subscribe(topic string, ch chan string) {
	ps.mutex.Lock()
	ps.subs[topic] = append(ps.subs[topic], ch)
	ps.mutex.Unlock()
}

func (ps *PubSub) SubscribeWithBuffer(topic string) <-chan string {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()
	ch := make(chan string, 1)
	ps.subs[topic] = append(ps.subs[topic], ch)
	return ch
}

func (ps *PubSub) Publish(topic string, msg string) {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()
	if ps.closed {
		return
	}
	for _, ch := range ps.subs[topic] { // 这里的ch如果没有设置buffer，那么这里消费ch就会阻塞
		ch <- msg
	}
}

func (ps *PubSub) Publish2(topic string, msg string) {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()
	if ps.closed {
		return
	}

	for _, ch := range ps.subs[topic] {
		go func(channel chan string) { // 非阻塞，但请注意，开启的新的goroutine本身会因客户端慢消费而阻塞
			channel <- msg
		}(ch)
	}
}

func (ps *PubSub) Close() {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	if !ps.closed {
		ps.closed = true
		// 关闭要把剩下的subs全部关闭
		for _, sub := range ps.subs {
			for _, ch := range sub {
				close(ch)
			}
		}
	}
}

func main() {
	ps := NewPubSub()
	ch1 := ps.SubscribeWithBuffer("tech")
	ch2 := ps.SubscribeWithBuffer("travel")
	ch3 := ps.SubscribeWithBuffer("travel")

	listener := func(name string, ch <-chan string) {
		for i := range ch {
			fmt.Printf("[%s] got %s\n", name, i)
		}
		fmt.Printf("[%s] done\n", name)
	}

	go listener("1", ch1)
	go listener("2", ch2)
	go listener("3", ch3)

	pub := func(topic string, msg string) {
		fmt.Printf("Publishing @%s: %s\n", topic, msg)
		ps.Publish2(topic, msg)
		time.Sleep(1 * time.Millisecond)
	}

	time.Sleep(50 * time.Millisecond)

	pub("tech", "tablets")
	pub("health", "vitamins")
	pub("tech", "robots")
	pub("travel", "beaches")
	pub("travel", "hiking")
	pub("tech", "drones")

	time.Sleep(50 * time.Millisecond)
	ps.Close()
	time.Sleep(50 * time.Millisecond)
}
