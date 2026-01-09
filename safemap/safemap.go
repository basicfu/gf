package safemap

import (
	"sync"
)

type Map[K comparable, V any] struct {
	mp    map[K]V
	mutex *sync.RWMutex
}

func New[K comparable, V any](cap ...int) *Map[K, V] {
	c := 0
	if len(cap) > 0 {
		c = cap[0]
	}
	return &Map[K, V]{
		mp:    make(map[K]V, c),
		mutex: &sync.RWMutex{},
	}
}

func (a *Map[K, V]) Set(k K, v V) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.mp[k] = v
}
func (a *Map[K, V]) Get(k K) (v V) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	v, _ = a.mp[k]
	return
}
func (a *Map[K, V]) GetExists(k K) (v V, ok bool) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	v, ok = a.mp[k]
	return
}
func (a *Map[K, V]) Exists(k K) (ok bool) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	_, ok = a.mp[k]
	return
}

func (a *Map[K, V]) Delete(k K) (existed bool) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if _, existed = a.mp[k]; existed {
		delete(a.mp, k)
	}
	return
}
func (a *Map[K, V]) Len() int {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return len(a.mp)
}

// 注意是在锁内循环的
func (a *Map[K, V]) Range(run func(k K, v V) bool) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	for k, v := range a.mp {
		if !run(k, v) {
			return
		}
	}
}
func (a *Map[K, V]) Map() map[K]V {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	data := make(map[K]V, len(a.mp))
	for k, v := range a.mp { //浅copy注意
		data[k] = v
	}
	return data
}
func (a *Map[K, V]) GetSet(k K, v V) (old V, existed bool) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	old, existed = a.mp[k]
	a.mp[k] = v
	return
}
