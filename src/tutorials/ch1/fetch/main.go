package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	hello(2)
	for _, url := range os.Args[1:] {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}
		b, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}
		fmt.Printf("%s", b)
	}
}

func hello2() int {
	a := 2
	b := 3
	return a + b
}

func hello(a int) int {
	c := a + 2
	return c
}
