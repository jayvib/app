package list

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestList_AddHead(t *testing.T) {

	t.Run("adding a node should increment the size by 1", func(t *testing.T){
		l := &List{}
		l.AddHead(1)
		assert.Equal(t, 1, l.Size())
	})

	t.Run("adding a head node on empty list", func(t *testing.T){
		l := &List{}
		l.AddHead(1)
		assertHeadNode(t, l, 1)
	})

	t.Run("adding a head with 1-item list", func(t *testing.T){
		l := &List{}
		l.AddHead(1)
		l.AddHead(2)

		assertHeadNode(t, l, 2)
		assertTailNode(t, l, 1)
		assertHeadNextNode(t, l, 1)
	})

	t.Run("asserting the tail after adding a new head node to the list", func(t *testing.T){
		l := &List{}
		l.AddHead(1)
		assertTailNode(t, l, 1)
	})
}

func TestList_AddTail(t *testing.T) {
	t.Run("adding a tail on an empty list", func(t *testing.T){
		l := &List{}
		l.AddTail(3)
		assertHeadNode(t, l, 3)
	})

	t.Run("adding a tail on an 1 item list", func(t *testing.T){
		l := &List{}
		l.AddHead(1)
		l.AddTail(2)
		assertTailNodeIterative(t, l, 2)
	})
}

func assertTailNodeIterative(t *testing.T, l *List, want int) {
	t.Helper()
	var gotTail *Node
	l.iterate(func(n *Node)bool{
		if n == nil {
			return true
		} else {
			gotTail = n
			return false
		}
	})
	assert.Equal(t, want, gotTail.value)
}

func assertHeadNextNode(t *testing.T, l *List, want int) {
	t.Helper()
	require.NotNil(t, l.head)
	headsNext := l.head.next
	assert.Equal(t, want, headsNext.value)
}

func assertTailNode(t *testing.T, l *List, want int) {
	t.Helper()
	gotTail := l.tail
	require.NotNil(t, gotTail)
	assert.Equal(t, want, gotTail.value)
}


func assertHeadNode(t *testing.T, l *List, want int) {
	t.Helper()
	gotHead := l.head
	require.NotNil(t, gotHead)

	got := gotHead.value
	assert.Equal(t, want, got)
}
