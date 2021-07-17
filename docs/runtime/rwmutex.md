# RWMutex —— 细粒度的读写锁

我们之前有讲过 [Mutex 互斥锁](https://github.com/MarsonShine/GolangStudy/blob/master/docs/runtime/mutex.md)。这是在任何时刻下只允许一个 goroutine 执行的串行化的锁。而现在这个 RWMutex 就是在 Mutex 的基础上进行了拓展能支持多个 goroutine 持有读锁，而在尝试持有写锁时就会如 Mutex 一样就会陷入等待锁的释放。它是一种细粒度的锁。虽然可以允许多次持有读锁，但是 Go 团队还特意嘱咐，**为了确保锁的可用性，不能用于递归读锁。一个阻塞的锁要排除正在持有锁的新读。**

那么上面说到的这些功能，RWMutex 是如何实现的呢？首先我们来看它的内部结构：

```go
type RWMutex struct {
	w           Mutex  // held if there are pending writers
	writerSem   uint32 // semaphore for writers to wait for completing readers
	readerSem   uint32 // semaphore for readers to wait for completing writers
	readerCount int32  // number of pending readers
	readerWait  int32  // number of departing readers
}
```

只有 5 个对象，其中最重要的就是 Mutex 锁的字段 w，它就是实现写锁的关键。

- writerSem 是写等待读完成的信号量
- readerSem 是读等待写完成的信号量
- readerCount 正处于读锁的个数
- readerWait 尝试获取写锁时读等待的个数（这个怎么理解？）

其中还有一个全局的常数变量 rwmutexMaxReaders，表示最多的读操作。

我们先来看写锁

## Lock/UnLock 写锁/解锁

```go
func (rw *RWMutex) Lock() {
	...
	rw.w.Lock()
	r := atomic.AddInt32(&rw.readerCount, -rwmutexMaxReaders) + rwmutexMaxReaders
	// Wait for active readers.
	if r != 0 && atomic.AddInt32(&rw.readerWait, r) != 0 {
		runtime_SemacquireMutex(&rw.writerSem, false, 0)
	}
	if race.Enabled {
		race.Enable()
		race.Acquire(unsafe.Pointer(&rw.readerSem))
		race.Acquire(unsafe.Pointer(&rw.writerSem))
	}
}
```

这里直接用到了 Mutex 互斥锁来保证只有一个 goroutine 能进来。**接下来就会判断在获取写锁的时候如果还存在其他的读锁没有释放，那么这个时候就会陷入睡眠进入等待者队列中等待所有的读锁被释放之后唤醒**。

> 可能有些人对这个限制有些不懂，其实这就是为了保证锁的区间的读的值顺序性的正确性。因为在获取写的时候，目的就是进行写操作，所谓我就必须要在此时还存在其他可能会读这个变量的读锁全部释放才行。

而释放写锁就是 UnLock 操作了。如果调用此操作时，本就没有上锁那么就会直接抛异常。

```go
func (rw *RWMutex) Unlock() {
	...
	// Announce to readers there is no active writer.
	r := atomic.AddInt32(&rw.readerCount, rwmutexMaxReaders)
	if r >= rwmutexMaxReaders {
		race.Enable()
		throw("sync: Unlock of unlocked RWMutex")
	}
	// Unblock blocked readers, if any.
	for i := 0; i < int(r); i++ {
		runtime_Semrelease(&rw.readerSem, false, 0)
	}
	// Allow other writers to proceed.
	rw.w.Unlock()
	...
}
```

如果还存在读锁时，那么就会进入 [runtime.Semrelease](https://github.com/MarsonShine/go/blob/aa4e0f528e1e018e2847decb549cfc5ac07ecf20/src/runtime/sema.go#L159) 对那些阻塞的读锁解锁（找到对应的信号量等待者队列然后弹出唤醒）。最后释放 w 锁。

## RLock/RUnlock 读锁/解锁

```go
func (rw *RWMutex) RLock() {
	...
	if atomic.AddInt32(&rw.readerCount, 1) < 0 {
		// A writer is pending, wait for it.
		runtime_SemacquireMutex(&rw.readerSem, false, 0)
	}
	...
}
```

读锁就非常简单了，仅仅只是对 readerCount 字段自增。这里的判断要注意，这个判断成立说明有协程调用了 rw.Lock 获取了写锁。所以就要等待其它协程的释放。

知道读锁的机制，那么就能想到释放读锁其实就是撤销读锁，将 readerCount 字段减1即可。

```
func (rw *RWMutex) RUnlock() {
	...
	if r := atomic.AddInt32(&rw.readerCount, -1); r < 0 {
		// Outlined slow-path to allow the fast-path to be inlined
		rw.rUnlockSlow(r)
	}
	...
}
```

同样在释放读锁时会判断 r 是否为负数，如果为负数就说明有其它协程获取了写锁，就会进入 rUnlockSlow 方法。

```go
func (rw *RWMutex) rUnlockSlow(r int32) {
	if r+1 == 0 || r+1 == -rwmutexMaxReaders {
		race.Enable()
		throw("sync: RUnlock of unlocked RWMutex")
	}
	// A writer is pending.
	if atomic.AddInt32(&rw.readerWait, -1) == 0 {
		// The last reader unblocks the writer.
		runtime_Semrelease(&rw.writerSem, false, 1)
	}
}
```

如果锁状态已经是解锁状态则抛异常。

**如果是只剩下一个读等待，则释放写信号量通知其他正在尝试持有写锁的协程上锁。**

## 关于信号量的细节

我们上面分析了读写锁的上锁与解锁的过程，其实有一个点不知道大家有没有注意。就是关于信号量的操作对象的细节。

1. 调用 Lock 获取写锁，会持有 writerSem 信号
2. 调用 Unlock 释放写锁时，会释放 readerSem 信号
3. 调用 RLock 获取读锁时，会持有 readerSem 信号
4. 调用 RUnlock 释放读锁时，会释放 writerSem 信号

大家有没有发现其中的规律，这么做的目的是什么呢？

也就是说：我们在获取写锁之前，会先等待读锁的释放操作。而在获取读锁时，会先等待写锁的释放操作。

我们用反证法来假设这个场景：我这里有一个连续的写操作；那么也就是说我要连续反复的调用 Lock + Unlock 操作。如果没有上面的信号量的互相牵制，那么就很容易出现读操作没法执行的问题，也就是说会”饿死“。

所以 RWMutex 加入读写信号量的机制是为了更好达到 RW 的目的，而不是一直 W。

## 总结

- 在调用 Lock 获取写锁时，会先等待 RUnlock 将其 readerCount 置为 0，然后成功获取写锁。
  - 还有一个操作是将 readerCount - rwmutexMaxReaders，其目的是为了阻塞后续的 RLock 操作。即在读取写锁其他任何读写操作都不允许了。
- 在调用 Unlock 释放写锁时，会通知所有读操作，解锁那些阻塞的读锁，然后成功释放写锁。

