package stack

type spyList struct {
	list []interface{}
}

func (s *spyList) Push(d interface{}) {
	s.list = append(s.list, d)
}

func (s *spyList) Pop() (data interface{}) {
	s.list, data = s.list[:s.Len()-1], s.list[s.Len()-1]
	return
}

func (s *spyList) Len() int {
	return len(s.list)
}

func (s *spyList) Top() (data interface{}) {
	s.list, data = s.list[1:], s.list[0]
	return
}

func (s *spyList) IsEmpty() bool {
	return s.Len() <= 0
}

