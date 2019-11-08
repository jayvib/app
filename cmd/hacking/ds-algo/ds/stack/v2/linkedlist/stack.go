package linkedlist

type Node struct {
	next  *Node
	value interface{}
}

type Stack struct {
	head  *Node
	count int
}

func (s *Stack) Push(v interface{}) {
	s.head = &Node{value: v, next: s.head}
	s.count++
}

func (s *Stack) Pop() (value interface{}) {

	if s.head == nil {
		return
	}

	value = s.head
	s.head = s.head.next
	s.count--
	return
}

func (s *Stack) Len() int {
	return s.count
}

func (s *Stack) IsEmpty() bool {
	return s.Len() == 0
}

func (s *Stack) Peek() (value interface{}, ok bool) {
	return
}
