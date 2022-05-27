package cyclicbarrier

import (
	"context"

	"github.com/marusama/cyclicbarrier"
	"golang.org/x/sync/semaphore"
)

type H2O struct {
	semaH *semaphore.Weighted         // 氢原子的信号量
	semaO *semaphore.Weighted         // 氧原子的信号量
	b     cyclicbarrier.CyclicBarrier // 循环栅栏，用来控制合成
}

func New() *H2O {
	return &H2O{
		semaH: semaphore.NewWeighted(2),
		semaO: semaphore.NewWeighted(1),
		b:     cyclicbarrier.New(3),
	}
}

// 氧原子处理
func (h20 *H2O) oxygen(releaseOsygen func()) {
	h20.semaO.Acquire(context.Background(), 1)
	releaseOsygen()
	h20.b.Await(context.Background()) // 循环栅栏，等待其它处理到达一起放行
	h20.semaO.Release(1)              // 释放氢原子
}

func (h2o *H2O) hydrogen(releaseHydrogen func()) {
	h2o.semaH.Acquire(context.Background(), 1)
	releaseHydrogen()
	h2o.b.Await(context.Background())
	h2o.semaH.Release(1)
}
