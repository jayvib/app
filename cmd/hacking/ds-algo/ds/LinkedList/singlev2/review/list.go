package list

import "fmt"

type Node struct {
	value int
	next *Node
}

type List struct {
	head, tail *Node
	count int
}

func (l *List) Size() int {
	return l.count
}

func (l *List) IsEmpty() bool {
	return l .count == 0
}

func (l *List) AddHead(val int) {
	var newNode *Node
	if l.IsEmpty() {
		newNode = &Node{value: val}
		l.tail = newNode
	} else {
		newNode = &Node{value: val, next: l.head}
	}
	l.head = newNode
	l.count++
}

func (l *List) AddTail(val int) {
	if l.IsEmpty() {
		l.head = &Node{value: val}
	} else {
		var tail *Node
		l.iterate(func(n *Node)bool{
			if n == nil {
				return true
			} else {
				tail = n
				fmt.Println("got tail:", tail.value)
				return false
			}
		})
		tail = &Node{value:val}
	}
	l.count++
}

func (l *List) iterate(fn func(n *Node)(stop bool)) {
	if l.IsEmpty(){
		fmt.Println("WARNING: Empty List")
		return
	}

	h := l.head
	for n := h.next;; n = n.next {
		stop := fn(n)
		if stop {
			return
		}
	}
}

