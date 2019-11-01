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
	length int
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
	for tmp := l.head; tmp != nil; tmp = tmp.next {
		fmt.Fprintln(out, tmp.value)
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
