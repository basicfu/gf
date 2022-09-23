package gmap

import (
	"github.com/basicfu/gf/internal/rwmutex"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Map     = AnyAnyMap // Map is alias of AnyAnyMap.
	HashMap = AnyAnyMap // HashMap is alias of AnyAnyMap.
)
type TMap[K primitive.ObjectID | string, V any] struct {
	mu   rwmutex.RWMutex
	data map[K]V
}

func New[K primitive.ObjectID | string, V any](safe ...bool) *TMap[K, V] {
	return &TMap[K, V]{
		mu:   rwmutex.Create(safe...),
		data: make(map[K]V),
	}
}

func (m *TMap[K, V]) Set(k K, v V) {
	defer m.mu.Unlock()
	m.mu.Lock()
	if m.data == nil {
		m.data = make(map[K]V)
	}
	m.data[k] = v
}

func (m *TMap[K, V]) Get(k K) (v V) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.data != nil {
		v, _ = m.data[k]
	}
	return
}
func (m *TMap[K, V]) GetExists(k K) (v V, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.data != nil {
		v, ok = m.data[k]
	}
	return
}
func (m *TMap[K, V]) Map() map[K]V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if !m.mu.IsSafe() {
		return m.data
	}
	data := make(map[K]V, len(m.data))
	for k, v := range m.data {
		data[k] = v
	}
	return data
}

func (m *TMap[K, V]) Remove(k K) (v V, ok bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data != nil {
		if v, ok = m.data[k]; ok {
			delete(m.data, k)
		}
	}
	return
}
func (m *TMap[K, V]) Exists(k K) bool {
	var ok bool
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.data != nil {
		_, ok = m.data[k]
	}
	return ok
}
