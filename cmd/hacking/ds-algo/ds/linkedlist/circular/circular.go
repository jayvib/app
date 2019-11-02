package circular

import (
	"errors"
	"fmt"
)

var ErrEmptyList = errors.New("empty list")

type Node struct {
	value int
	next  *Node
}

type CircularLinkedList struct {
	tail  *Node
	count int
}

func (c *CircularLinkedList) Size() int {
	return c.count
}

func (c *CircularLinkedList) IsEmpty() bool {
	return c.count == 0
}

func (c *CircularLinkedList) Peek() (val int, err error) {
	if c.IsEmpty() {
		return 0, ErrEmptyList
	}
	return c.tail.next.value, nil
}

func (c *CircularLinkedList) AddHead(val int) {
	node := &Node{value: val} // first, we create node with given value and its next pointing to null.
	switch {
	case c.addNodeWhenEmpty(node): // If the list is empty then tail of the list will point to it. In addition, the next
	default: // If the list is not empty then the next of the new node will be next of the tail.
		node.next = c.tail.next
		c.tail.next = node
	}
	c.count++
}

func (c *CircularLinkedList) AddTail(val int) {
	node := &Node{value: val}
	switch {
	case c.addNodeWhenEmpty(node):
	default:
		node.next = c.tail.next // new tails link the existing head
		c.tail.next = node      // the existing tail will be the previous value of the new tail
		c.tail = node
	}
	c.count++
}

func (c *CircularLinkedList) RemoveHead() (val int, err error) {
	if c.tail == nil {
		return 0, ErrEmptyList
	}

	val = c.tail.next.value

	if c.tail == c.tail.next {
		c.tail = nil
	} else {
		c.tail.next = c.tail.next.next
	}
	c.count--
	return val, nil
}

func (c *CircularLinkedList) IsPresent(v int) (found bool) {
	if c.tail == nil {
		return false
	}
	found = false
	c.iterate(func(n *Node)(stop bool){
		if n.value == v {
			found = true
			stop = true
		}
		return
	})

	return found
}

func (c *CircularLinkedList) Reset() {
	c.tail = nil
	c.count = 0
}

func (c *CircularLinkedList) RemoveNode(key int) (ok bool) {

	head, tail := c.tail.next, c.tail
	// key is in the head
	if head.value == key {
		if head == tail {
			c.tail = nil
		}
		c.count--
		return true
	}

	// somewhere in the body
	for curr, prev := head.next, tail.next; curr != nil; curr, prev = curr.next, curr {
		if curr.value == key {
			if curr != tail {
				prev.next = curr.next
			} else {
				prev.next = head
				c.tail = prev
			}
			c.count--
			return true
		}
		if curr == tail {
			break
		}
	}

	return false
}

func (c *CircularLinkedList) CopyListReversed() (*CircularLinkedList) {

	reversedList := new(CircularLinkedList)

	c.iterate(func(n *Node) bool {
		reversedList.AddHead(n.value)
		return false
	})

	return reversedList
}

func (c *CircularLinkedList) CopyList() (*CircularLinkedList) {
	copyL := new(CircularLinkedList)

	c.iterate(func(n *Node) bool {
		copyL.AddTail(n.value)
		return false
	})
	return copyL
}

func (c *CircularLinkedList) iterate(fn func(n *Node) (stop bool)) {
	if c.IsEmpty() {
		return
	}
	for n := c.tail.next; c != nil; n = n.next {
		stop := fn(n)

		if stop {
			return
		}
		if n == c.tail {
			break
		}
	}
}

func(c *CircularLinkedList) Print() {
	if c.IsEmpty() {
		return
	}
	c.iterate(func(n *Node) bool {
		fmt.Println(n.value)
		return false
	})
}

func (c *CircularLinkedList) addNodeWhenEmpty(node *Node) (ok bool) {
	if c.IsEmpty() {
		c.tail = node
		c.tail.next = node
		return true
	}
	return
}