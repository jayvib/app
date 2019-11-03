package slice

import (
	"fmt"
	"io"
	"os"
)

type StackInt struct {
	s []int
}

func (s *StackInt) Pop() int {
	if s.IsEmpty() {
		return 0
	}
	v := s.s[s.Len()]
	s.s = s.s[:s.Len()-1]
	return v
}

func (s *StackInt) Push(value int) {
	s.s = append(s.s, value)
}

func (s *StackInt) Len() int {
	return len(s.s)
}

func (s *StackInt) IsEmpty() bool {
	return s.Len() == 0
}

func (s *StackInt) Print() {
	s.Fprint(os.Stdout)
}

func (s *StackInt) Fprint(w io.Writer) {
	for _, e := range s.s {
		fmt.Fprintln(w, e)
	}
}
