package list

import (
	"fmt"
	"io"
	"os"
)

type Node struct {
	value    int
	next     *Node
	previous *Node
}

type List struct {
	head, tail *Node
	length     int
}

func (l *List) Len() int {
	return l.length
}

func (l *List) IsEmpty() bool {
	return l.length == 0
}

func (l *List) Peek() (value int, ok bool) {
	if l.IsEmpty() {
		return 0, false
	}
	return l.head.value, true
}

func (l *List) AddHead(value int) {
	node := &Node{value: value}
	if l.IsEmpty() {
		l.head = node
		l.tail = node
	} else {
		l.head.previous = node
		node.next = l.head
		l.head = node
	}
	l.length++
}

func (l *List) AddTail(value int) {
	newNode := &Node{value: value}
	if l.IsEmpty() {
		l.head = newNode
		l.tail = newNode
	} else {
		l.tail.next = newNode
		newNode.previous = l.tail
		l.tail = newNode
	}
	l.length++
}

func (l *List) Print() {
	l.PrintOut(os.Stdout)
}

func (l *List) PrintOut(out io.Writer) {
	l.Range(func(n *Node) bool {
		fmt.Fprintln(out, n.value)
		return false
	})
}

func (l *List) Range(fn func(*Node) bool) {
	for tmp := l.head; tmp != nil; tmp = tmp.next {
		if fn(tmp) {
			return
		}
	}
}

func (l *List) RemoveHead() (val int, ok bool) {
	if l.IsEmpty() {
		return 0, false
	}

	head := l.head
	l.head = l.head.next
	if l.head == nil {
		l.tail = nil
	} else {
		l.head.previous = nil
	}

	l.length--
	return head.value, true
}

func (l *List) IsExists(value int) bool {
	exists := false
	l.Range(func(n *Node) bool {
		if n.value == value {
			exists = true
			return true
		}
		return false
	})
	return exists
}

func (l *List) RemoveNode(value int) (ok bool) {

	// Is the value the head?
	if curr := l.head; curr.value == value {
		if curr.next == nil {
			l.head = nil
			l.tail = nil
		} else {
			l.head = curr.next
			l.head.previous = nil
		}
		l.length--
		return true
	}

	for curr := l.head; curr != nil; curr = curr.next {
		if curr.value == value {
			prev := curr.previous
			next := curr.next // check if the next if not nil.
			if next != nil {
				prev.next = next
				next.previous = prev
			} else {
				prev.next = nil
			}
			l.length--
		}
	}
	return
}

func (l *List) Reverse() {
	for curr, next := l.head, l.head.next; curr != nil; curr, next = next, next.next {
		curr.next = curr.previous	// Switch
		curr.previous = next
		if next == nil {
			l.tail = l.head
			l.head = curr
			return
		}
	}
}

