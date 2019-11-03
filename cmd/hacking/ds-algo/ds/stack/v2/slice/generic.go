package slice

import "errors"

var ErrEmptyStack = errors.New("empty stack")

type GenericStack struct {
	v []interface{}
}

func (s *GenericStack) Push(v interface{}) {
	s.v = append(s.v, v)
}

func (s *GenericStack) Pop() (value interface{}, err error) {
	switch {
	case s.len() == 0:
		err = ErrEmptyStack
	case s.len() == 1:
		value = s.pop()
		s.v = s.v[:0]
	default:
		value = s.pop()
		s.v = s.v[:s.len()-1]
	}
	return
}

func (s *GenericStack) pop() (interface{}) {
	return s.v[s.len()-1]
}

func (s *GenericStack) len() int {
	return len(s.v)
}
