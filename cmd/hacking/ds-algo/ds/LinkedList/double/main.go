package main

import "fmt"

type Node struct {
	value int
	next *Node
	prev *Node
}

type DoublyLinkedList struct {
	count int
	head *Node
	tail *Node
}

func (l *DoublyLinkedList) IsEmpty() bool {
	if l.count == 0 {
		return true
	}
	return false
}

func (l *DoublyLinkedList) Peek() (int, bool) {
	if l.IsEmpty() {
		return 0, false
	}
	return l.head.value, true
}

func (l *DoublyLinkedList) AddHead(value int) {
	newNode := &Node{value: value}
	if l.IsEmpty() {
		l.tail, l.head = newNode, newNode
	} else {
		l.head.prev = newNode
		newNode.next = l.head
		l.head = newNode
	}
	l.count++
}

func (l *DoublyLinkedList) AddTail(value int) {
	newNode := &Node{value: value}
	if l.IsEmpty() {
		l.head, l.tail = newNode, newNode
	} else {
		newNode.prev = l.tail
		l.tail.next = newNode
		l.tail = newNode
	}
	l.count++
}

func (l *DoublyLinkedList) RemoveHead() (int, bool) {
	if l.IsEmpty() {
		return 0, false
	}

	value := l.head.value
	l.head = l.head.next
	if l.head == nil {
		l.tail = nil
	} else {
		l.head.prev = nil
	}
	l.count--
	return value, true
}

func main() {
	dlist := new(DoublyLinkedList)
	dlist.AddHead(1)
	dlist.AddHead(2)
	dlist.AddTail(3)
	fmt.Println(dlist.count)
	fmt.Println(dlist.head.value, dlist.head.next.value, dlist.head.next.next.value)
	v, ok := dlist.RemoveHead()
	if ok {
		fmt.Println("head value removed:", v)
	}
	fmt.Println(dlist.head.value)
}