package store

import (
	"testing"
)

type testUser struct {
	id   uint64
	name string
}

func newTestUser(id uint64, name string) testUser {
	return testUser{
		id:   id,
		name: name,
	}
}

func (u testUser) StoreKey() uint64 {
	return u.id
}

func TestStore(t *testing.T) {
	s := New[uint64, testUser]()

	t.Run("store", func(t *testing.T) {
		s.Store(newTestUser(0, "Alex"))
		s.Store(newTestUser(1, "Marcus"))
	})

	t.Run("load", func(t *testing.T) {
		u, ok := s.Load(0)
		if !ok {
			t.Fail()
			return
		}
		if u.name != "Alex" {
			t.Fail()
		}
	})

	t.Run("each", func(t *testing.T) {
		passed := 2
		s.Each(func(_ testUser) {
			passed--
		})
		if passed != 0 {
			t.Fail()
		}
	})

	t.Run("filter", func(t *testing.T) {
		r := s.Filter(func(u testUser) bool {
			return len(u.name) > 4
		})
		if len(r) != 1 {
			t.Fail()
		}
	})

	t.Run("filter-one", func(t *testing.T) {
		u, ok := s.FilterOne(func(u testUser) bool {
			return u.name == "Marcus"
		})
		if !ok {
			t.Fail()
		}
		if u.name != "Marcus" {
			t.Fail()
		}
	})

	t.Run("len", func(t *testing.T) {
		if s.Len() != 2 {
			t.Fail()
		}
	})

	t.Run("values", func(t *testing.T) {
		r := s.Values()
		if len(r) != 2 {
			t.Fail()
		}
	})

	t.Run("delete", func(t *testing.T) {
		s.Delete(0)
		s.Delete(1)
		if s.Len() != 0 {
			t.Fail()
		}
	})
}
