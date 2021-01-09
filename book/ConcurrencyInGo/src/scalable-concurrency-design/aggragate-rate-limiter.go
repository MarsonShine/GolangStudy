package main

import (
	"context"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func ratelimitExample2() {
	defer log.Printf("Done.")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	apiConnection := Open2()
	var wg sync.WaitGroup
	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ReadFile(context.Background())
			if err != nil {
				log.Printf("cannot ReadFile: %v", err)
			}
			log.Printf("ReadFile")
		}()
	}
	log.Printf("ReadFile")

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConnection.ResolveAddress(context.Background())
			if err != nil {
				log.Printf("cannot ResolveAddress: %v", err)
			}
			log.Printf("ResolveAddress")
		}()
	}
	log.Printf("ResolveAddress")
	wg.Wait()
}

// 聚合速率限制器，方便管理
type RateLimiter interface {
	Wait(context.Context) error
	Limit() rate.Limit
}

type multilimiter struct {
	limiters []RateLimiter
}

func MultiLimiter(limiters ...RateLimiter) *multilimiter {
	byLimit := func(i, j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	}
	sort.Slice(limiters, byLimit) // 2 按照每个RateLimiter的 Limit() 行排序。
	return &multilimiter{
		limiters: limiters,
	}
}

func (l *multilimiter) Wait(ctx context.Context) error {
	for _, l := range l.limiters {
		if err := l.Wait(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (l *multilimiter) Limit() rate.Limit {
	return l.limiters[0].Limit() // 3 因为我们在multiLimiter实例化时对子RateLimiter实例进行排序，所以我们可以简单地返回限制性最高的limit，这将是切片中的第一个元素。
}

// 对每秒和每分钟都进行限流
func Open2() *APIConnection2 {
	secondLimit := rate.NewLimiter(Per(2, time.Second), 1)   // 1 每秒的极限
	minuteLimit := rate.NewLimiter(Per(10, time.Minute), 10) // 2 每分钟的极限设置为 10，每秒的限制将确保我们不会因请求而使系统过载。
	return &APIConnection2{
		rateLimiter: MultiLimiter(secondLimit, minuteLimit),
	}
}

type APIConnection2 struct {
	rateLimiter RateLimiter
}

func (a *APIConnection2) ReadFile(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	// Pretend we do work here
	return nil
}
func (a *APIConnection2) ResolveAddress(ctx context.Context) error {
	if err := a.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	// Pretend we do work here
	return nil
}
