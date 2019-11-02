package single

type List struct {
	next *Node
	count int
}

type Node struct {
	next *Node
	value int
}

func (l *List) Size() int {
	return l.count
}

func (l *List) IsEmpty() bool {
	return l.count == 0
}

func (l *List) AddHead(val int) {
	l.count++
}
