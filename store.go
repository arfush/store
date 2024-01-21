package store

import (
	"sync"
)

type Storer[K comparable] interface {
	StoreKey() K
}

type Store[K comparable, V Storer[K]] struct {
	mtx sync.Mutex
	m   map[K]V
}

func New[K comparable, V Storer[K]]() *Store[K, V] {
	return &Store[K, V]{
		m: make(map[K]V),
	}
}

func (s *Store[K, V]) Store(v V) bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if _, ok := s.m[v.StoreKey()]; ok {
		return false
	}
	s.m[v.StoreKey()] = v
	return true
}

func (s *Store[K, V]) Load(k K) (V, bool) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	v, ok := s.m[k]
	return v, ok
}

func (s *Store[K, V]) Each(f func(V)) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for _, v := range s.m {
		f(v)
	}
}

func (s *Store[K, V]) Filter(f func(V) bool) (r []V) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for _, v := range s.m {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

func (s *Store[K, V]) FilterOne(f func(V) bool) (V, bool) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for _, v := range s.m {
		if f(v) {
			return v, true
		}
	}
	return *new(V), false
}

func (s *Store[K, V]) Len() int {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return len(s.m)
}

func (s *Store[K, V]) Values() []V {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	r := make([]V, 0, len(s.m))
	for _, v := range s.m {
		r = append(r, v)
	}
	return r
}

func (s *Store[K, V]) Delete(k K) bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if _, ok := s.m[k]; !ok {
		return false
	}
	delete(s.m, k)
	return true
}
