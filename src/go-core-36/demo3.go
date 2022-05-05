package main

import (
	"flag"
	"fmt"
	"os"
)

var name string

func init() {
	flag.CommandLine = flag.NewFlagSet("", flag.PanicOnError)
	flag.CommandLine.Usage = func() {
		fmt.Fprintf(os.Stderr, "Useage of %s:\n", "question")
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {
	fmt.Printf("Hello, %s!\n", name)
}
