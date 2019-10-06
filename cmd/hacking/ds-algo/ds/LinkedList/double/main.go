package main

import "fmt"

type Node struct {
	value int
	next  *Node
	prev  *Node
}

type DoublyLinkedList struct {
	count int
	head  *Node
	tail  *Node
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

	value := l.head.value // get the value of the head
	l.head = l.head.next // overwrite it with the next node
	if l.head == nil {
		l.tail = nil
	} else {
		l.head.prev = nil
	}
	l.count--
	return value, true
}

func (l *DoublyLinkedList) RemoveNode(key int) bool {
	curr := l.head // get the head as the starting point
	if curr == nil { // hen the current head is empty
		return false // then return false
	}
	if curr.value == key { // when the key is the current key of the head
		curr := curr.next // assign the next node as the current
		l.count-- // decrement the list count
		if curr != nil { // when the new current is not empty
			l.head = curr // set the current as the new head
			l.head.prev = nil // since the head don't have previous set it to empty
		} else {
			l.tail = nil // no more value in the list
		}
		return true
	}
	if curr.next != nil { // traversing in the list
		if curr.next.value == key { // standing in the current, peek the next node's value, then its equal
			curr.next = curr.next.next // then set the next node as the next next node since the next node will be remove
			if curr.next == nil { // when the new next node is nil
				l.tail = curr // set the current node as the tail of the list.
			} else {
				curr.next.prev = curr // assign the prev node in the next node as the current.
			}
			l.count--
			return true
		}
		curr = curr.next
	}
	return false
}

func main() {
	dlist := new(DoublyLinkedList)
	dlist.AddHead(1)
	dlist.AddHead(2)
	dlist.AddTail(3)
	fmt.Println(dlist.count)
	fmt.Println(dlist.head.value, dlist.head.next.value, dlist.head.next.next.value)
	//v, ok := dlist.RemoveHead()
	//if ok {
	//	fmt.Println("head value removed:", v)
	//}
	//fmt.Println(dlist.head.value)
	dlist.RemoveNode(1)
	fmt.Println(dlist.head.next.value)
	fmt.Println(dlist.head.next.prev.value)
}
