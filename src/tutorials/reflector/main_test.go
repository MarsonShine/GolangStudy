package main_test

import (
	"fmt"
	"testing"
	"time"
	"tutorials/reflector/format"
)

func Test(t *testing.T) {
	// The pointer values are just examples, and may vary from run to run.
	//!+time
	var x int64 = 1
	var d time.Duration = 1 * time.Nanosecond
	fmt.Println(format.Any(x))                  // "1"
	fmt.Println(format.Any(d))                  // "1"
	fmt.Println(format.Any([]int64{x}))         // "[]int64 0x8202b87b0"
	fmt.Println(format.Any([]time.Duration{d})) // "[]time.Duration 0x8202b87e0"
	//!-time
}
