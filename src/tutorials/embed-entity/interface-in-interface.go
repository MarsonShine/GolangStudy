package main

// 接口嵌入接口
type Reader interface {
	Read(p []byte) (n int, err error)
}
type Writer interface {
	Write(p []byte) (n int, err error)
}

// 嵌入接口，可以做各种组合
// 具体案例可以详见 package io 包
type ReaderWriter interface {
	Reader
	Writer
}
