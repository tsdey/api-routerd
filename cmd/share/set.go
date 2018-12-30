// SPDX-License-Identifier: Apache-2.0

package share

type Set struct {
	m map[string]bool
}

func NewSet() *Set {
	s := &Set{}
	s.m = make(map[string]bool)

	return s
}

func (s *Set) Add(value string) {
	s.m[value] = true
}

func (s *Set) Remove(value string) {
	delete(s.m, value)
}

func (s *Set) Contains(value string) bool {
	_, c := s.m[value]

	return c
}
