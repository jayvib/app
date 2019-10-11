package singlev2

type Node struct {
	value int
	next  *Node
}

type List struct {
	head  *Node
	length int
}

// AddHead create a new node with the value passed to the function as argument
//
// While creating the new node the reference stored in head is passed as
// argument to Node so that the next reference will start pointing to the node
// or null which is referenced y the head node.
func (l *List) AddHead(value int) {
	l.length++
	l.head = &Node{value: value, next: l.head}
}


func (l *List) AddTail(value int) {
	l.length++
	curr := l.head

	if curr == nil {
		l.head = &Node{value: value}
		return
	} else {
		for curr.next != nil {
			curr = curr.next
		}
		curr.next = &Node{value:value}
		return
	}
}

