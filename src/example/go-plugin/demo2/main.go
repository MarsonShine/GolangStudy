package main

import (
	"fmt"

	"github.com/demo2/pkg"
)

func main() {
	fmt.Println("try to LoadPlugin...")
	err := pkg.LoadPlugin("../demo2-plugins/plugin1.so")
	if err != nil {
		fmt.Println("LoadPlugin error:", err)
		return
	}
	fmt.Println("LoadPlugin ok")
	err = pkg.LoadPlugin("../demo2-plugins/plugin1.so")
	if err != nil {
		fmt.Println("Re-LoadPlugin error:", err)
		return
	}
	fmt.Println("Re-LoadPlugin ok")
}
