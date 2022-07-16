package utils

import "fmt"

type ST interface {
	int | int64 | string
}

type Set[T ST] struct {
	m map[T]struct{}
}

func NewSet[T ST]() *Set[T] {
	return &Set[T]{m: make(map[T]struct{})}
}

func (s *Set[T]) Add(e T) {
	s.m[e] = struct{}{}
}

func (s *Set[T]) Find(e T) bool {
	_, ok := s.m[e]
	return ok
}

func (s *Set[T]) Del(e T) {
	delete(s.m, e)
}

func (s *Set[T]) Range() []T {
	var ks []T
	for key, _ := range s.m {
		ks = append(ks, key)
	}
	return ks
}

func (s *Set[T]) String() string {
	return fmt.Sprint(s.Range())
}
