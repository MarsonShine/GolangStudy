// 代码引用自：https://github.com/chebyrash/promise/blob/master/promise.go

package promise

import (
	"errors"
	"fmt"
	"sync"
)

type Promise[T any] struct {
	result  T
	err     error
	pending bool        // 表示异步函数执行的状态
	mutex   *sync.Mutex // 控制对共享状态的并发更新
	wg      *sync.WaitGroup
}

func New[T any](f func(resolve func(T), reject func(error))) *Promise[T] {
	if f == nil {
		panic("f not null")
	}
	p := &Promise[T]{
		pending: true,
		mutex:   &sync.Mutex{},
		wg:      &sync.WaitGroup{},
	}
	p.wg.Add(1)
	go func() {
		defer p.handlePanic()
		f(p.resolve, p.reject)
	}()

	return p
}

func (p *Promise[T]) resolve(resolution T) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if !p.pending {
		return
	}

	p.result = resolution
	p.pending = false
	p.wg.Done()
}

func (p *Promise[T]) reject(err error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if !p.pending {
		return
	}
	p.err = err
	p.pending = false
	p.wg.Done()
}

func (p *Promise[T]) handlePanic() {
	err := recover()
	if validErr, ok := err.(error); ok {
		p.reject(fmt.Errorf("panic recovery: %w", validErr))
	} else {
		p.reject(fmt.Errorf("panic recovery: %+v", err))
	}
}

func (p *Promise[T]) Awaiter() (T, error) {
	p.wg.Wait()
	return p.result, p.err
}

func Then[A, B any](promise *Promise[A], resolveA func(data A) B) *Promise[B] {
	return New(func(resolveB func(B), reject func(error)) {
		result, err := promise.Awaiter()
		if err != nil {
			reject(err)
			return
		}
		resolveB(resolveA(result))
	})
}

func Catch[T any](promise *Promise[T], rejection func(err error) error) *Promise[T] {
	return New(func(resolve func(T), reject func(error)) {
		result, err := promise.Awaiter()
		if err != nil {
			reject(rejection(err))
			return
		}
		resolve(result)
	})
}

// 直接返回给定值
func Resolve[T any](resolution T) *Promise[T] {
	return &Promise[T]{
		result:  resolution,
		pending: false,
		mutex:   &sync.Mutex{},
		wg:      &sync.WaitGroup{},
	}
}

//
func Reject[T any](err error) *Promise[T] {
	return &Promise[T]{
		err:     err,
		pending: false,
		mutex:   &sync.Mutex{},
		wg:      &sync.WaitGroup{},
	}
}

// 定义元组
type tuple[T1, T2 any] struct {
	_1 T1
	_2 T2
}

// 当所有的任务都resolve时返回
// 当所有的任务中有拒绝的，会立即reject返回
func All[T any](promises ...*Promise[T]) *Promise[[]T] {
	if len(promises) == 0 {
		return nil
	}

	return New(func(resolve func([]T), reject func(error)) {
		valsChan := make(chan tuple[T, int], len(promises))
		errsChan := make(chan error, 1)

		for i, p := range promises {
			index := i
			// 等待所有resolve返回data
			_ = Then(p, func(data T) T {
				valsChan <- tuple[T, int]{_1: data, _2: index}
				return data
			})
			// 一旦发生错误就返回
			_ = Catch(p, func(err error) error {
				errsChan <- err
				return err
			})
		}

		resolutions := make([]T, len(promises))
		for i := 0; i < len(promises); i++ {
			select {
			case val := <-valsChan:
				resolutions[val._2] = val._1
			case err := <-errsChan:
				reject(err)
				return
			}
		}

		resolve(resolutions)
	})
}

func Any[T any](promises ...*Promise[T]) *Promise[T] {
	if len(promises) == 0 {
		return nil
	}
	return New(func(resolve func(T), reject func(error)) {
		valsChan := make(chan T, 1)
		errsChan := make(chan tuple[error, int], len(promises))

		for idx, p := range promises {
			idx := idx // https://golang.org/doc/faq#closures_and_goroutines
			_ = Then(p, func(data T) T {
				valsChan <- data
				return data
			})
			_ = Catch(p, func(err error) error {
				errsChan <- tuple[error, int]{_1: err, _2: idx}
				return err
			})
		}

		errs := make([]error, len(promises))
		for idx := 0; idx < len(promises); idx++ {
			select {
			case val := <-valsChan:
				resolve(val)
				return
			case err := <-errsChan:
				errs[err._2] = err._1
			}
		}

		errCombo := errs[0]
		for _, err := range errs[1:] {
			errCombo = errors.New(err.Error() + errCombo.Error()) // errors.Wrap(err, errCombo.Error())
		}
		reject(errCombo)
	})
}
