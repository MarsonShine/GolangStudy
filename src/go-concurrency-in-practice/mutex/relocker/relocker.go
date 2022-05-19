package relocker

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/petermattis/goid"
)

// mutex 不是可重入锁
// 可重入锁的实现
type RecursiveMutex struct {
	sync.Mutex
	owner     int64 // 当前持有锁的goroutine id
	recursion int32 // 允许可重入的次数
}

func (m *RecursiveMutex) Lock() {
	gid := goid.Get()
	if atomic.LoadInt64(&m.owner) == gid {
		m.recursion++
		return
	}
	m.Mutex.Lock()
	// 第一次获取锁
	atomic.StoreInt64(&m.owner, gid)
	m.recursion = 1
}

func (m *RecursiveMutex) Unlock() {
	gid := goid.Get()
	// 检查非持久锁尝试解锁，则panic
	if atomic.LoadInt64(&m.owner) != gid {
		panic(fmt.Sprintf("wrong the owner(%d): %d!", m.owner, gid))
	}
	// 更新计数
	m.recursion--
	if m.recursion != 0 { // 说明还没有完全释放
		return
	}
	atomic.StoreInt64(&m.owner, -1)
	m.Mutex.Unlock()
}

type TokenRecursionMutex struct {
	sync.Mutex
	token     int64
	recursion int32
}

func (m *TokenRecursionMutex) Lock(token int64) {
	if atomic.LoadInt64(&m.token) == token {
		m.recursion++
		return
	}
	m.Mutex.Lock()
	atomic.StoreInt64(&m.token, token)
	m.recursion = 1
}
func (m *TokenRecursionMutex) Unlock(token int64) {
	if atomic.LoadInt64(&m.token) != token {
		panic(fmt.Sprintf("wrong the token(%d): %d!", m.token, token))
	}
	m.recursion--
	if m.recursion != 0 {
		return
	}
	atomic.StoreInt64(&m.token, -1)
	m.Mutex.Unlock()
}

func GetGoroutineId() int64 {
	var buf [64]byte
	// 从堆栈中获取goroutine id
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return int64(id)
}

// doc
/*
第一步：首先，我们获取运行时的 g 指针，反解出对应的 g 的结构。每个运行的 goroutine 结构的 g 指针保存在当前 goroutine 的一个叫做 TLS 对象中
第二步：从TLS中获取 goroutine 结构的 g 指针
第三步：从 g 指针中取出 goroutine id
最后：上面三步获取 goroutine id 已经有组件做了这件事：https://github.com/petermattis/goid

除了上述通过反射hacker的方式，还可以通过传递凭证token的方式，通过token来判断groutine的归属
*/
