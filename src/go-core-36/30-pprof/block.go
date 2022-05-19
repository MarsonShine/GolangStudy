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
	blockProfileName = "blockprofile.out"
	blockProfileRate = 2
	debug            = 0
)

func block() {
	f, err := common.CreateFile("", blockProfileName)
	if err != nil {
		fmt.Printf("block profile creation error: %v\n", err)
		return
	}
	defer f.Close()
	startBlockProfile()
	if err = common.Execute(op.BlockProfile, 10); err != nil {
		fmt.Printf("execute error: %v\n", err)
		return
	}
	if err := stopBlockProfile(f); err != nil {
		fmt.Printf("block profile stop error: %v\n", err)
		return
	}
}

func startBlockProfile() {
	runtime.SetBlockProfileRate(blockProfileRate)
}

func stopBlockProfile(f *os.File) error {
	if f == nil {
		return errors.New("nil file")
	}
	return pprof.Lookup("block").WriteTo(f, debug)
}
