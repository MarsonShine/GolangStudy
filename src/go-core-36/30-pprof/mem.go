package main

import (
	"errors"
	"example/gocore36/30-pprof/common"
	"example/gocore36/30-pprof/common/op"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
)

var (
	memProfileName = "memprofile.out"
	memProfileRate = 8
)

func mem() {
	f, err := common.CreateFile("", memProfileName)
	if err != nil {
		fmt.Printf("memory profile creation error: %v\n", err)
		return
	}
	defer f.Close()
	startMemProfile()
	if err = common.Execute(op.MemProfile, 10); err != nil {
		fmt.Printf("execute error: %v\n", err)
		return
	}
	if err := stopMemProfile(f); err != nil {
		fmt.Printf("memory profile stop error: %v\n", err)
		return
	}
}

func startMemProfile() {
	runtime.MemProfileRate = memProfileRate
}

func stopMemProfile(f *os.File) error {
	if f == nil {
		return errors.New("nil file")
	}
	return pprof.WriteHeapProfile(f)
}
