// SPDX-License-Identifier: Apache-2.0

package share

type Set struct {
	m map[string]string
}

func NewSet() *Set {
	s := &Set{}
	s.m = make(map[string]string)

	return s
}

func (s *Set) Add(value string) {
	s.m[value] = value
}

func (s *Set) Remove(value string) {
	delete(s.m, value)
}

func (s *Set) Contains(value string) bool {
	_, c := s.m[value]

	return c
}
