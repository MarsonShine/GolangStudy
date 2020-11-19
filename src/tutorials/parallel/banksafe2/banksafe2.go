package banksafe2

var (
	sema    = make(chan struct{}, 1)
	balance int
)

func Deposit(amount int) {
	sema <- struct{}{} // 获取令牌
	balance += amount
	<-sema //释放锁
}

func Balance() int {
	sema <- struct{}{} //获取令牌
	b := balance
	<-sema //释放锁
	return b
}
