package slice

type StackInt struct {
	s []int
}

func (s *StackInt) Pop() int {
	return 0
}

func (s *StackInt) Push(value int) {

}

func (s *StackInt) Len() int {
	return len(s.s)
}

func (s *StackInt) IsEmpty() bool {
	return s.Len() == 0
}
