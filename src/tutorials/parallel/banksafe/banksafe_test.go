package banksafe_test

import (
	"fmt"
	"testing"
	"tutorials/parallel/banksafe"
)

// 利用非缓存 channel 实现线程安全
func TestBank(t *testing.T) {
	done := make(chan struct{})

	// Alice
	go func() {
		banksafe.Deposit(200)
		fmt.Println("=", banksafe.Balance())
		done <- struct{}{}
	}()

	// Bob
	go func() {
		banksafe.Deposit(100)
		done <- struct{}{}
	}()

	// Wait for both transactions.
	<-done
	<-done

	if got, want := banksafe.Balance(), 300; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}
}
