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

func (s *Stack) Top() (value interface{}) {
	return s.l.Front().Value
}

func (s *Stack) Len() int {
	return s.l.Len()
}

func(s *Stack) IsEmpty() bool {
	return s.Len() == 0
}

func InfixToPostFix2(exp string) string {
	return ""
}

func InfixToPostfix(expn string) string {
	fmt.Println(expn)
	stk := &Stack{l: list.New()}
	output := ""
	for _, ch := range expn {
		if ch <= '9' && ch >= '0' {
			output = output + string(ch)
		} else {
			switch ch {
			case '+', '-', '*', '/', '%', '^':
				for stk.IsEmpty() == false && precedence(ch) <=
					precedence(stk.Top().(rune)) {
					out := stk.Pop().(rune)
					output = output + " " + string(out)
				}
				stk.Push(ch)
				output = output + " "
			case '(':
				stk.Push(ch)
			case ')':
				for stk.IsEmpty() == false && stk.Top() != '(' {
					out := stk.Pop().(rune)
					output = output + " " + string(out) + " "
				}
				if stk.IsEmpty() == false && stk.Top() == '(' {
					stk.Pop()
				}
			}
		}
	}
	for stk.IsEmpty() == false {
		out := stk.Pop().(rune)
		output = output + string(out) + " "
	}
	return output
}

func precedence(x rune) int {
	if x == '(' {
		return (0)
	}
	if x == '+' || x == '-' {
		return (1)
	}
	if x == '*' || x == '/' || x == '%' {
		return (2)
	}
	if x == '^' {
		return (3)
	}
	return (4)
}

func main() {
	expn := "10+((3))*5/(16-4)"
	value := InfixToPostfix(expn)
	fmt.Println("Infix Expn: ", expn)
	fmt.Println("Postfix Expn: ", value)
}