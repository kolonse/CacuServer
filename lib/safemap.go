package lib

import (
	"sync"
)

type SafeMap struct {
	m   map[interface{}]interface{}
	mtx sync.Mutex
}

func (sm *SafeMap) MapIndex(key interface{}) (interface{}, bool) {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()
	v, ok := sm.m[key]
	return v, ok
}

func (sm *SafeMap) SetMapIndex(key interface{}, value interface{}) {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()
	sm.m[key] = value
}

func (sm *SafeMap) Keys() []interface{} {
	var ret []interface{}
	for k, _ := range sm.m {
		ret = append(ret, k)
	}
	return ret
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		m: make(map[interface{}]interface{}),
	}
}
