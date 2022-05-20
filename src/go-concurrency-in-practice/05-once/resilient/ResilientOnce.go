package resilient

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type ResilientOnce struct {
	done uint32
	m    sync.Mutex
}

func (o *ResilientOnce) Do(f func() error) error {
	if atomic.LoadUint32(&o.done) == 1 {
		return nil
	}
	return o.slowDo(f)
}

func (o *ResilientOnce) slowDo(f func() error) error {
	o.m.Lock()
	defer o.m.Unlock()
	var err error
	if o.done == 0 {
		err = f()
		if err == nil { // 只有初始化成功才将标记置为已初始化
			atomic.StoreUint32(&o.done, 1)
		}
	}
	return err
}

func (o *ResilientOnce) Done() bool {
	return atomic.LoadUint32((*uint32)(unsafe.Pointer(o))) == 1
}
