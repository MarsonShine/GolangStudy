package main

import (
	"fmt"
	"os/exec"
)

func main() {
	// dateCmd := exec.Command("date")

	// dateOut, err := dateCmd.Output()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("> date")
	// fmt.Println(string(dateOut))

	// grepCmd := exec.Command("grep", "hello")

	// // 获取命令行管道中的输入和输出
	// grepIn, _ := grepCmd.StdinPipe()
	// grepOut, _ := grepCmd.StdoutPipe()
	// grepCmd.Start()
	// grepIn.Write([]byte("hello grep\ngoodbye grep"))
	// grepIn.Close()
	// grepBytes, _ := ioutil.ReadAll(grepOut)
	// grepCmd.Wait()

	// fmt.Println("> grep hello")
	// fmt.Println(string(grepBytes))

	// lsCmd := exec.Command("bash", "-c", "ls -a -l -h")
	lsCmd := exec.Command("ls")
	lsOut, err := lsCmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println("> ls -a -l -h")
	fmt.Println(string(lsOut))
}
