package stack

type simpleStack struct {
	list []interface{}
}

func (s *simpleStack) Push(d interface{}) {
	s.list = append(s.list, d)
}

func (s *simpleStack) Pop() (data interface{}) {
	s.list, data = s.list[:s.Len()-1], s.list[s.Len()-1]
	return
}

func (s *simpleStack) Len() int {
	return len(s.list)
}

func (s *simpleStack) Top() (data interface{}) {
	s.list, data = s.list[1:], s.list[0]
	return
}

func (s *simpleStack) IsEmpty() bool {
	return s.Len() <= 0
}

