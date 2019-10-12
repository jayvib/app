package singlev2

type Node struct {
	value int
	next  *Node
}

type List struct {
	head  *Node
	length int
	tail *Node // in order to avoid traversing the entire list
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

// AddTail creates a new node with the value inside it.
//
// When the list is empty. Next of new node is nil. And head will
// store the reference to the newly created node.
//
// If list is not empty then traverse until the end of the list
// and the new node is added to the end of the list.
func (l *List) AddTail(value int) {
	l.length++
	curr := l.head

	switch curr {
	case nil:
		l.head = &Node{value: value}
	default:
		for curr.next != nil {
			curr = curr.next
		}

		newNode := &Node{value: value}
		curr.next =  newNode
		l.tail = newNode
	}
}

func (l *List) SortedInsert(value int) {
	newNode := &Node{value: value}
	curr := l.head

	// for the head
	if curr == nil || curr.value > value{
		newNode.next = l.head
		l.head = newNode
		return
	}

	// so the new value is somewhere in the body
	for curr.next != nil || curr.next.value < value {
		curr = curr.next
	}

	newNode.next = curr.next
	curr.next = newNode
}
