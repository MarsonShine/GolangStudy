package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
)

func errorPassingExample() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)
	err := runJob("1")
	if err != nil {
		msg := "There was an unexpected issue; please report this as a bug."
		if _, ok := err.(IntermediateErr); ok { //1
			msg = err.Error()
		}
		handleError(1, err, msg) //2
	}
}

type MyError struct {
	Inner      error
	Message    string
	StackTrace string
	Misc       map[string]interface{}
}

func wrapError(err error, messageF string, msgArgs ...interface{}) MyError {
	return MyError{
		Inner:      err, // 1, 包装的错误，当发生错误时希望能看到最低级的错误
		Message:    fmt.Sprintf(messageF, msgArgs...),
		StackTrace: string(debug.Stack()),        // 2, 记录创建错误时的堆栈跟踪
		Misc:       make(map[string]interface{}), // 3, 可以存放一些其他字段数据，可以存储并发ID，以及其他上下文信息
	}
}

func (err MyError) Error() string {
	return err.Message
}

type LowLevelErr struct {
	error
}

func isGloballyExec(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, LowLevelErr{wrapError(err, err.Error(), nil)} // 1, 自定义错误来封装os.Stat的原始错误
	}
	return info.Mode().Perm()&0100 == 0100, nil
}

type IntermediateErr struct {
	error
}

func runJob(id string) error {
	const jobBinPath = "/bad/job/binary"
	isExecutable, err := isGloballyExec(jobBinPath)
	if err != nil {
		return err //1
	} else if isExecutable == false {
		return wrapError(nil, "job binary is not executable")
	}
	return exec.Command(jobBinPath, "--id="+id).Run() // 1, 没有用自定义错误信息封装起来，这会产生问题
}

func runJob2(id string) error {
	const jobBinPath = "/bad/job/binary"
	isExecutable, err := isGloballyExec(jobBinPath)
	if err != nil {
		return IntermediateErr{wrapError(err, "cannot run job %q: requisite binaries not available", id)} // 1, 使用自定义错误。我们想隐藏工作未运行原因的底层细节，因为这对于用户并不重 要。
	} else if isExecutable == false {
		return wrapError(nil, "cannot run job %q: requisite binaries are not executable", id)
	}
	return exec.Command(jobBinPath, "--id="+id).Run()
}

func handleError(key int, err error, message string) {
	log.SetPrefix(fmt.Sprintf("[logID: %v]: ", key))
	log.Printf("%#v", err) //3
	fmt.Printf("[%v] %v", key, message)
}
