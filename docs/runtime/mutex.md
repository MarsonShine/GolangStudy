# Mutex 互斥锁

## 概要描述

mutex 是 go 提供的同步原语。用于多个协程之间的同步协作。在大多数底层框架代码中都会用到这个锁。

mutex 总过有三个状态

- mutexLocked: 表示占有锁
- mutexWoken: 表示唤醒
- mutexStarving: 表示等待锁的饥饿状态（从正常模式进入饥饿状态）

## 具体实现

首先得清楚 Mutex 的结构

```go
type Mutex struct {
	state int32
	sema  uint32
}
```

### Lock 持有锁

```go
func (m *Mutex) Lock() {
	// Fast path: grab unlocked mutex.
	if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
		if race.Enabled {
			race.Acquire(unsafe.Pointer(m))
		}
		return
	}
	// Slow path (outlined so that the fast path can be inlined)
	m.lockSlow()
}
```

比较锁的状态，如果是初始状态 0，则获取当前锁并将锁的状态置为 mutexLocked 之后返回，即 1。其他协程则在释放锁之前走 `m.lockSlow` 代码分支。

m.lockSlow 代码很复杂，里面有一些比较的变量为了方便理解，我标记起来：

```go
func (m *Mutex) lockSlow() {
	var waitStartTime int64
    starving := false	// goruotine饥饿标识
    awoke := false	// 唤醒标记
    iter := 0	// 自旋的次数
    ...
}
```

- waitStartTime: 表示等待唤醒的时间，如果等待的时间超过了 1ms 则会走锁的饥饿模式（后面会讲到）
- starving: 表示是饥饿标识
- awoke: 是否被唤醒的标识
- iter: 等待锁的次数

在申明这些变量之后，程序就会进入一个循环阶段，在这个循环体内，runtime 会根据当前锁的状态 `m.state` 以及等待锁的次数以及系统的自身运行情况来判断程序是否进入自旋等待。

```go
// lockSlow
for { // 无论是新协程还是老协程都在循环获取锁
    	// 锁是普通状态，锁还没有被释放，则自旋
		if old&(mutexLocked|mutexStarving) == mutexLocked && runtime_canSpin(iter) {
			if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 &&
				atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
				awoke = true
			}
			runtime_doSpin()
			iter++
			old = m.state	// 再次获取锁的状态，如此每次都要检查锁的状态
			continue
		}
  	...
}
```

上面的第一个判断语句就是判断了当前锁的状态是否为普通模式（只有普通模式才会进入自旋）。除此之外 [runtime_canSpin](https://github.com/golang/go/blob/master/src/runtime/proc.go#L6358) 还会判断当前自旋的次数是否超过 4 次、当前操作系统是否为多核 cpu 、 `GOMAXPROCS > 1` 以及至少有一个正在运行的 P 并且 P 下的局部队列 `runq` 是空闲的。

在没有自旋的情况下，则正常往下执行，就会执行 [runtime_SemacquireMutex](https://github.com/golang/go/blob/2ebe77a2fda1ee9ff6fd9a3e08933ad1ebaea039/src/runtime/sema.go#L98) 方法，目的就是获取当前锁的 `m.sema` 变量的地址作为信号量将当前的等待者进入 `semaRoot` 的等待优先级队列（链表结构）中。并且根据传的变量 `queueLifo` 来控制是传入到队列尾部还是头部。如果 `queueLifo=true` 说明之前已经处理等待状态，则就会插入到队列的头部等待唤醒，否则就进入队列的尾部。在这个 [runtime_SemacquireMutex](https://github.com/golang/go/blob/2ebe77a2fda1ee9ff6fd9a3e08933ad1ebaea039/src/runtime/sema.go#L98) 方法中会不断尝试获取锁并进入休眠状态等待信号量释放。一旦获取锁立即返回执行后续的代码。

```go
// lockSlow
if atomic.CompareAndSwapInt32(&m.state, old, new) {	// 成功设置锁的新状态
    // 原来的锁已经释放，并且不是饥饿状态，正常请求锁并返回
    if old&(mutexLocked|mutexStarving) == 0 {
      break // locked the mutex with CAS
    }
    // 往下开始处理饥饿状态
    // 如果以前就在队列里面，则加入到队列头
    queueLifo := waitStartTime != 0
    if waitStartTime == 0 {
      waitStartTime = runtime_nanotime()
    }
    // 阻塞等待
    runtime_SemacquireMutex(&m.sema, queueLifo, 1)
    // 唤醒之后检查锁是否处于饥饿状态
    starving = starving || runtime_nanotime()-waitStartTime > starvationThresholdNs
    old = m.state
    // 如果处于饥饿状态，直接占有锁，返回
    if old&mutexStarving != 0 {
      if old&(mutexLocked|mutexWoken) != 0 || old>>mutexWaiterShift == 0 {
        throw("sync: inconsistent mutex state")
      }
      delta := int32(mutexLocked - 1<<mutexWaiterShift)
      if !starving || old>>mutexWaiterShift == 1 { // 最后一个waiter或者是已经不饥饿了，则清除饥饿标记
        delta -= mutexStarving
      }
      atomic.AddInt32(&m.state, delta)
      break
    }
    awoke = true
    iter = 0
  } else {
    old = m.state
  }
}

//runtime_SemacquireMutex
func semacquire1(addr *uint32, lifo bool, profile semaProfileFlags, skipframes int) {
  ...
  for {
      lockWithRank(&root.lock, lockRankRoot)
      // Add ourselves to nwait to disable "easy case" in semrelease.
      atomic.Xadd(&root.nwait, 1)
      // Check cansemacquire to avoid missed wakeup.
      if cansemacquire(addr) {
        atomic.Xadd(&root.nwait, -1)
        unlock(&root.lock)
        break
      }
      // Any semrelease after the cansemacquire knows we're waiting
      // (we set nwait above), so go to sleep.
      root.queue(addr, s, lifo)
      goparkunlock(&root.lock, waitReasonSemacquire, traceEvGoBlockSync, 4+skipframes)
      if s.ticket != 0 || cansemacquire(addr) {
        break
      }
	}
	if s.releasetime > 0 {
		blockevent(s.releasetime-t0, 3+skipframes)
	}
	releaseSudog(s)
}
```

如上面代码所示，在插入队列之后就会根据标识 `starving` 以及 `runtime_nanotime() - waitStartTime > 1ms` 来判断当前等待锁模式是否进入饥饿模式。

在进入饥饿模式之前，等待者优先级队列中的第一个 g 被唤醒，如果其他新到达的 g 的优先级要比等待队列中的要高，所以这个等待着就会获取锁失败，那么就会转而进入到等待队列中的头部等待下次唤醒。而此时如果等待唤醒的队列等待时间超过了 1ms，那么就会进入饥饿模式，此时新到来的 g 就不会参与占有锁，也不会自旋，而是进入到等待着优先级队列中的队尾等待唤醒。

只有在这个等待队列为空，或者最后一个等待着的等待唤醒时间小于 1ms，那就会推出饥饿模式。

所以这个饥饿模式就是为了防止等待队列中的 g 陷入无限的等待状态。

#### 如何根据 m.sema 获取等待者队列？

首先 rutime 会根据 m.sema 的地址通过哈希计算来生成一个 table，每个 mutex 的信号量地址对应一个表。每个表都有一个 semaRoot 的对象，这个对象包含一个 `treap *sudog` 链表结构队列。通过 `semaroot.g = getg()` 把当前的 g 绑定起来就进入到等待着队列了。

```
// mutex.go
func semroot(addr *uint32) *semaRoot {
	return &semtable[(uintptr(unsafe.Pointer(addr))>>3)%semTabSize].root
}
```

### Unlock 释放锁

知道了 Lock 的实现细节，那么释放锁我们就相对来说比较简单了。我们从源码的代码量就能知道要比 Lock 简单不少：

```go
func (m *Mutex) Unlock() {
	if race.Enabled {
		_ = m.state
		race.Release(unsafe.Pointer(m))
	}
	new := atomic.AddInt32(&m.state, -mutexLocked)
	if new != 0 {
		m.unlockSlow(new)
	}
}
```

首先就是通过原子操作将 m.state 减 1 操作让其回到初始状态 0。如果该锁还没有释放则进入 unlockSlow 操作

```
func (m *Mutex) unlockSlow(new int32) {
	if (new+mutexLocked)&mutexLocked == 0 {
		throw("sync: unlock of unlocked mutex")
	}
	if new&mutexStarving == 0 {
		old := new
		for {
			if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken|mutexStarving) != 0 {
				return
			}
			// Grab the right to wake someone.
			new = (old - 1<<mutexWaiterShift) | mutexWoken
			if atomic.CompareAndSwapInt32(&m.state, old, new) {
				runtime_Semrelease(&m.sema, false, 1)
				return
			}
			old = m.state
		}
	} else {
		runtime_Semrelease(&m.sema, true, 1)
	}
}
```

该操作仅仅做了当前锁状态是否为饥饿模式，如果不是就会正常唤醒等待队列中的 g，并通过 sema 信号量释放该信号量将所有权交给等待者。如果是饥饿模式，则直接将当前所有权交给队列中的下一个等待者。

## 总结

Mutex 在占有锁的时候会有两种模式，一种是正常模式，等待者会进入一个等待者队列等待被唤醒。其他 g 尝试获取锁时会根据系统自身情况以及是否有 runq 空闲的 P 来进入自旋。在自旋的过程中会计算锁的状态。由于正常模式下等待者队列中的 G 可能会无限等待下去，所以这个时候就会进入饥饿状态。会将锁的所有权直接交给被唤醒的等待者直到队列为空退出饥饿模式。

在进入等待者队列是会获取当前锁的信号量，这样保证了同一个时刻不会有多个 G 占有同一个锁。

解锁过程就是对应加锁过程的逆过程。

如果当前锁的状态不对，则直接抛异常；当互斥锁为饥饿模式时，直接将所有权交给等待者。当处于正常模式时，如果没有锁或者当前的锁状态不都为 0，则直接返回。否则就会释放信号量唤醒等待者队列中的等待者获取锁。