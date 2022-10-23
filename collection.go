/*
@author: sk
@date: 2022/10/23
*/
package main

type Set struct {
	map0 map[any]struct{}
}

func NewSet() *Set {
	return &Set{map0: make(map[any]struct{})}
}

func (s *Set) Add(key any) {
	s.map0[key] = struct{}{}
}

func (s *Set) Has(key any) bool {
	_, ok := s.map0[key]
	return ok
}
