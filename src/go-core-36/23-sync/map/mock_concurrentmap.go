package main

import "sync"

type (
	// 通过map和rwmutex实现并发安全的map
	MockConcurrentMap struct {
		m  map[interface{}]interface{}
		mu sync.RWMutex
	}
)

func NewMockConcurrentMap() *MockConcurrentMap {
	return &MockConcurrentMap{
		m: make(map[interface{}]interface{}),
	}
}

func (m *MockConcurrentMap) Delete(key interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.m, key)
}

func (m *MockConcurrentMap) Load(key interface{}) (value interface{}, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok = m.m[key]
	return
}

func (m *MockConcurrentMap) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	actual, loaded = m.m[key]
	if loaded {
		return
	}
	m.m[key] = value
	actual = value
	return
}

func (m *MockConcurrentMap) Range(f func(key, value interface{}) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.m {
		if !f(k, v) {
			break
		}
	}
}

func (m *MockConcurrentMap) Store(key, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m[key] = value
}
