package main

import "time"

// 超时和取消的最佳实践

var reallyLongCalculation = func(done <-chan interface{}, value interface{}) interface{} {
	intermediateResult := longCalculation(value)
	select {
	case <-done:
		return nil
	default:
	}
	return longCalculation(intermediateResult)
}

func longCalculation(value interface{}) interface{} {
	time.Sleep(1 * time.Minute)
	return struct{}{}
}
