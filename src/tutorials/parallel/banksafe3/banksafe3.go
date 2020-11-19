package banksafe3

import "sync"

// goroutine 在调用 mu.Lock 就会获得锁
// 其他的 goroutine 在调用之前会一直阻塞，直到那个协程调用完 mu.Unlock

var (
	mu      sync.Mutex // 互斥锁
	balance int
)

func Deposit(amount int) {
	mu.Lock()
	balance += amount
	mu.Unlock()
}

func Balance() int {
	mu.Lock()
	b := balance
	mu.Unlock()
	return b
}

func BalanceBetter() int {
	mu.Lock()
	defer mu.Unlock()
	return balance
}

func WithDraw(amount int) bool {
	Deposit(-amount)
	if BalanceBetter() < 0 {
		Deposit(amount)
		return false
	}
	return true
}

// 这会导致死锁，因为调用方法 WithDrawWithLock 时会首先调用 mu.Lock() 占有锁
// 当继续执行 BalanceBetter 时会再此尝试占有锁，但是由于锁对象 mu 已经被占据，所以会等待 mu 的释放，而此时该线程却在尝试占有锁，这就造成了死锁
func WithDrawWithLock(amount int) bool {
	mu.Lock()
	defer mu.Unlock()
	Deposit(-amount)
	if BalanceBetter() < 0 {
		Deposit(amount)
		return false
	}
	return true
}
