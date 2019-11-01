package list

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrEmptyList = errors.New("empty list")
	ErrNotFound  = errors.New("node not found")
)

type Node struct {
	value int
	next  *Node
}

type List struct {
	head, tail *Node
	count      int
}

func (l *List) Size() int {
	return l.count
}

func (l *List) IsEmpty() bool {
	return l.count == 0
}

func (l *List) AddHead(val int) {
	newNode := &Node{value: val}
	if l.IsEmpty() {
		l.tail = newNode
	} else {
		newNode.next = l.head
	}
	l.head = newNode
	l.count++
}

func (l *List) IsPresent(d int) bool {
	var found bool
	l.iterate(func(n *Node) bool {
		if n == nil {
			return true
		}
		if n.value == d {
			found = true
			return true
		}
		return false
	})
	return found
}

func (l *List) AddTail(val int) {
	newNode := &Node{value: val}
	if l.IsEmpty() {
		l.head = newNode
		l.tail = newNode
	} else {
		tail := l.tail
		tail.next = newNode
		l.tail = tail.next
	}
	l.count++
}

func (l *List) Print() {
	l.PrintOut(os.Stdout)
}

func (l *List) PrintOut(w io.Writer) {
	l.iterate(func(n *Node) bool {
		if n == nil {
			return true
		}

		fmt.Fprintln(w, n.value)
		return false
	})
}

func (l *List) RemoveHead() (val int, err error) {
	if l.IsEmpty() {
		return 0, ErrEmptyList
	}
	l.count--
	var head *Node
	head = l.head
	l.head = head.next
	return head.value, nil
}

func (l *List) RemoveNode(val int) error {
	l.count--

	// check first the head
	head := l.head
	if head.value == val {
		l.head = head.next
		return nil
	}

	for next, prev := head.next, head; next != nil; next, prev = next.next, next {
		if next.value == val {
			prev.next = next.next
			return nil
		}
	}

	return ErrNotFound
}

func (l *List) RemoveNodes(val int) {
	// check first the head
	if l.head.value == val {
		l.head = l.head.next
	}

	hasFound := false
	for current, prev := l.head.next, l.head; current != nil; current, prev = current.next, current {
		if current.value == val {
			if !hasFound {
				hasFound = true
			}
			prev.next = current.next
			current = prev
			l.count--
		}
	}
	if !hasFound {
		return  // TODO: Use an error
	}
}

func (l *List) Reverse() {
	curr := l.head
	var next, prev *Node
	for curr != nil {
		next = curr.next
		curr.next = prev
		prev = curr
		curr = next
	}

	l.head = prev
}

func (l *List) getTail() *Node {
	var currentTail *Node
	l.iterate(func(n *Node) bool {
		if n == nil {
			return true
		} else {
			currentTail = n
			return false
		}
	})
	return currentTail
}

func (l *List) FreeList() {
	l.head = nil
	l.count = 0
}

func (l *List) iterate(fn func(n *Node) (stop bool)) {
	if l.IsEmpty() {
		fmt.Println("WARNING: Empty List")
		return
	}

	h := l.head
	for n := h; ; n = n.next {
		stop := fn(n)
		if stop {
			return
		}
	}
}
