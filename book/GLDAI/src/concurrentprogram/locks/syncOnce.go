// once 只执行一次，可以实现单例
// Once.Do 方法，如果传入的参数已经执行，则会直接返回；如果没有执行过，就会调用 sync.Once.doSlow 执行传入的函数
package locks

import (
	"fmt"
	"sync"
)

func single() {
	o := &sync.Once{}
	for i := 0; i < 10; i++ {
		o.Do(func() {
			// 只会执行一次
			fmt.Println("only once")
		})
	}
}
