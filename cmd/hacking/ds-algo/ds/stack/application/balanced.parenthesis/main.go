package main

import (
	"container/list"
	"fmt"
)

type Stack struct {
	l *list.List
}

func (s *Stack) Push(v interface{}) {
	s.l.PushFront(v)
}

func (s *Stack) Pop() (value interface{}) {
	e := s.l.Front()
	value = s.l.Remove(e)
	return
}

func (s *Stack) Len() int {
	return s.l.Len()
}

func isBalancedParenthesis(expn string) bool {
	s := list.New()
	for _, c := range expn {
		switch c {
		case '{', '[', '(':
			s.PushFront(c)
		case '}':
			e := s.Front()
			v := s.Remove(e)
			if v != '{' {
				return false
			}
		case ']':
			e := s.Front()
			v := s.Remove(e)
			if v != '[' {
				return false
			}
		case ')':
			e := s.Front()
			v := s.Remove(e)
			if v != '(' {
				return false
			}
		}
	}
	return s.Len() == 0
}

func main() {
	in := "[hey you]"
	fmt.Println("is balance:", isBalancedParenthesis(in))
}
