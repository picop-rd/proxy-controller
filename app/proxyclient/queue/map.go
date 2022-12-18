package queue

import (
	"sync"
)

type Map[T any] struct {
	q sync.Map
}

func NewMap[T any]() *Map[T] {
	return &Map[T]{
		q: sync.Map{},
	}
}

func (m *Map[T]) Get(k string) (T, bool) {
	v, ok := m.q.Load(k)
	if !ok {
		return *new(T), ok
	}
	vT, ok := v.(T)
	if !ok {
		return *new(T), ok
	}
	return vT, ok
}

func (m *Map[T]) Set(k string, v T) {
	m.q.Store(k, v)
}

func (m *Map[T]) Del(k string) {
	m.q.Delete(k)
}

func (m *Map[T]) Pop(k string) (T, bool) {
	v, ok := m.q.LoadAndDelete(k)
	if !ok {
		return *new(T), ok
	}
	vT, ok := v.(T)
	if !ok {
		return *new(T), ok
	}
	return vT, ok
}

func (m *Map[T]) Range(f func(k string, v T) bool) {
	m.q.Range(func(k, v any) bool {
		kStr, ok := k.(string)
		if !ok {
			return false
		}
		vT, ok := v.(T)
		if !ok {
			return false
		}
		return f(kStr, vT)
	})
}
