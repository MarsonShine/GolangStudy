package racecondition

import "fmt"

var balance int

func Deposit(amount int) { balance = balance + amount }
func Balance() int       { return balance }

func main() {
	// Alice:
	go func() {
		Deposit(200)                // A1
		fmt.Println("=", Balance()) // A2
	}()
	// Bob:
	go Deposit(100)
	// 存在并发问题，共同修改一个变量
}
