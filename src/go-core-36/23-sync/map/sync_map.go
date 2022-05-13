package main

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// 并发map，key 类型一定要合规
// 所以可以通过设置自定义结构
// 通过统一的入口来达到key类型检测的目的

type (
	IntStringMap struct {
		m sync.Map
	}

	// 注意：需要下载最新gopls以支持泛型编程
	GenericConcurrentMap[K int | int64, V any] struct {
		m sync.Map
	}

	// 非泛型的Map版本
	// 因为map对编译器来说不是类型安全，不会对key value最类型检查
	// 所以我们需要包装一个类来实现这个目的

	ConcurrentMap struct {
		m         sync.Map
		keyType   reflect.Type
		valueType reflect.Type
	}
)

func NewConcurrentMap(keyType, valueType reflect.Type) (*ConcurrentMap, error) {
	if keyType == nil {
		return nil, errors.New("nil key type")
	}
	if !keyType.Comparable() {
		return nil, fmt.Errorf("incomparable key type: %s", keyType)
	}
	if valueType == nil {
		return nil, errors.New("nil value type")
	}
	cMap := &ConcurrentMap{
		keyType:   keyType,
		valueType: valueType,
	}
	return cMap, nil
}

func (iMap *IntStringMap) Load(key int) (value string, ok bool) {
	v, ok := iMap.m.Load(key)
	if v != nil {
		value = v.(string)
	}
	return
}

func (iMap *IntStringMap) Range(f func(key int, value string) bool) {
	f1 := func(key, value interface{}) bool {
		return f(key.(int), value.(string))
	}
	iMap.m.Range(f1)
}

func (iMap *IntStringMap) LoadOrStore(key int, value string) (actual interface{}, loaded bool) {
	a, loaded := iMap.m.LoadOrStore(key, value)
	actual = a.(string)
	return
}

func (iMap *IntStringMap) Delete(key int) {
	iMap.m.Delete(key)
}

func (iMap *IntStringMap) Store(key int, value string) {
	iMap.m.Store(key, value)
}

func (iMap *GenericConcurrentMap[K, V]) Load(key int) (value string, ok bool) {
	v, ok := iMap.m.Load(key)
	if v != nil {
		value = v.(string)
	}
	return
}

func (iMap *GenericConcurrentMap[K, V]) Range(f func(key K, value V) bool) {
	f1 := func(key, value interface{}) bool {
		return f(key.(K), value.(V))
	}
	iMap.m.Range(f1)
}

func (iMap *GenericConcurrentMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	a, loaded := iMap.m.LoadOrStore(key, value)
	actual = a.(V)
	return
}

func (iMap *GenericConcurrentMap[K, V]) Delete(key K) {
	iMap.m.Delete(key)
}

func (iMap *GenericConcurrentMap[K, V]) Store(key K, value V) {
	iMap.m.Store(key, value)
}

func (cmap *ConcurrentMap) Load(key interface{}) (value interface{}, ok bool) {
	// 做类型检查
	if reflect.TypeOf(key) != cmap.keyType {
		return
	}
	return cmap.m.Load(key)
}

func (cmap *ConcurrentMap) Store(key int, value string) {
	// 类型检查
	if reflect.TypeOf(key) != cmap.keyType {
		panic(fmt.Errorf("wrong key type: %v", reflect.TypeOf(key)))
	}
	if reflect.TypeOf(value) != cmap.valueType {
		panic(fmt.Errorf("wrong value type: %v", reflect.TypeOf(value)))
	}
	cmap.m.Store(key, value)
}

func (cmap *ConcurrentMap) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	if reflect.TypeOf(key) != cmap.keyType {
		panic(fmt.Errorf("wrong key type: %v", reflect.TypeOf(key)))
	}
	if reflect.TypeOf(value) != cmap.valueType {
		panic(fmt.Errorf("wrong value type: %v", reflect.TypeOf(value)))
	}
	actual, loaded = cmap.m.LoadOrStore(key, value)
	return
}

func (cmap *ConcurrentMap) Range(f func(key, value interface{}) bool) {
	cmap.m.Range(f)
}

func main() {
	// gm := GenericConcurrentMap[int, string]{}
	// gm.Store(int64(64), "marsonshine") // 类型错误，编译期类型安全

}
