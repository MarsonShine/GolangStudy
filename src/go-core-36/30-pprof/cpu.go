package main

import (
	"errors"
	"example/gocore36/30-pprof/common"
	"example/gocore36/30-pprof/common/op"
	"fmt"
	"os"
	"runtime/pprof"
)

var cpuProfileName = "cpuprofile.out"

func cpu() {
	f, err := common.CreateFile("", cpuProfileName)
	if err != nil {
		fmt.Printf("CPU profile creation error: %v\n", err)
		return
	}
	defer f.Close()
	if err := startCPUProfile(f); err != nil {
		fmt.Printf("CPU profile start error: %v\n", err)
		return
	}
	if err = common.Execute(op.CPUProfile, 10); err != nil {
		fmt.Printf("execute error: %v\n", err)
		return
	}
	stopCPUProfile()
}

func startCPUProfile(f *os.File) error {
	if f == nil {
		return errors.New("nil file")
	}
	return pprof.StartCPUProfile(f)
}

func stopCPUProfile() {
	pprof.StopCPUProfile()
}
