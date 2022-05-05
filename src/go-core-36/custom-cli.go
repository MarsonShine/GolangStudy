package main

import (
	"flag"
	"fmt"
	"os"
)

var name string
var cmdLine *flag.FlagSet

func init() {
	cmdLine = flag.NewFlagSet("question", flag.ExitOnError)
	cmdLine.StringVar(&name, "name", "everyone", "请输入-name=value")
}

func main() {
	cmdLine.Parse(os.Args[1:])
	fmt.Printf("Hello, %s!\n", name)
}
