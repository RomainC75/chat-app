package typedsyncmap

import "sync"

type TSyncMap[K any, V any] struct {
	sm sync.Map
}

func NewSyncMap[K any, V any]() *TSyncMap[K, V] {
	return &TSyncMap[K, V]{}
}

func (tsm *TSyncMap[K, V]) Store(key K, value V) {
	tsm.sm.Store(key, value)
}

func (tsm *TSyncMap[K, V]) Delete(key K) {
	tsm.sm.Delete(key)
}

func (tsm *TSyncMap[K, V]) Load(key K) (V, bool) {
	value, ok := tsm.sm.Load(key)
	var emptyValue V
	if !ok {
		return emptyValue, ok
	}
	converted, ok := value.(V)
	return converted, ok
}

func (tsm *TSyncMap[K, V]) Range(f func(key K, value V) bool) {
	tsm.sm.Range(func(key, value any) bool {
		convKey, _ := key.(K)
		convValue, _ := value.(V)
		f(convKey, convValue)
		return true
	})
}

func (tsm *TSyncMap[K, V]) DeleteAll() {
	tsm.Range(func(key K, value V) bool {
		tsm.sm.Delete(key)
		return true
	})
}
