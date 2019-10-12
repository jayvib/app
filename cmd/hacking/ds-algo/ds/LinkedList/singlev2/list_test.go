package singlev2

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// An element can be inserted into a linked list
// 1. Insertion of an element at the start of linked list
// 2. Insertion of an element at the end of linked list
// 3. Insertion of an element at the Nth position in linked list
// 4. Insert element in sorted order in linked list

func TestList_AddHead(t *testing.T) {
	t.Run("Insertion of an element at the start of linked list", func(t *testing.T){
		list := List{}
		list.AddHead(1)

		want := &Node{value: 1}
		t.Run("After adding a new node the length of the list must be 1", func(t *testing.T){
			assert.Equal(t, 1, list.length)
		})

		t.Run("After adding a new node the new head of the list must be the added value", func(t *testing.T){
			assert.Equal(t, want, list.head)
		})

		t.Run("After adding a new node the underlying node of the new node must be nil", func(t *testing.T){
			assert.Nil(t, list.head.next)
		})
	})
}

func TestList_AddTail(t *testing.T) {

	t.Run("Adding a new node on the tail of the list the length must be 1", func(t *testing.T){
		l := &List{}
		l.AddTail(2)
		assert.Equal(t, 1,l.length)
	})

	t.Run("After adding a new node on the tail with empty List the tail must be head node", func(t *testing.T){
		l := &List{}
		l.AddTail(2)
		want := &Node{value: 2}
		assert.Equal(t, want, l.head)
	})

	t.Run("After adding a new node on the tail with head on the list the tail must be the added value", func(t *testing.T){
		newL := &List{}
		newL.AddHead(1)
		newL.AddTail(2)

		got := newL.head.next
		want := &Node{value:2}
		assert.Equal(t, want, got)
	})

	t.Run("After adding a new node on the tail with variable number of existing node the tail must be the added value", func(t *testing.T){
		l := &List{}
		l.AddHead(1)
		l.AddHead(2)
		l.AddHead(3)
		l.AddHead(4)

		l.AddTail(5)

		want := &Node{value: 5}
		got := l.tail

		assert.Equal(t, want, got)
	})
}

func TestList_SortedInsert(t *testing.T){
	l := &List{}
	l.AddHead(1)
	l.AddHead(3)

	l.SortedInsert(2)

	want := &Node{value: 2}
	var got *Node
	curr := l.head
	for curr.next != nil {
		if curr.value == want.value {
			got = curr
			break
		}
		curr = curr.next
	}
	assert.Equal(t, want.value, got.value)
}