package resilient

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type Once struct {
	sync.Once
}

func (o *Once) Done() bool {
	return atomic.LoadUint32((*uint32)(unsafe.Pointer(&o.Once))) == 1
}
