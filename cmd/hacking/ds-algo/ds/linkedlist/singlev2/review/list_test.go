package list

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestList_AddHead(t *testing.T) {

	t.Run("adding a node should increment the size by 1", func(t *testing.T) {
		l := &List{}
		l.AddHead(1)
		assert.Equal(t, 1, l.Size())
	})

	t.Run("adding a head node on empty list", func(t *testing.T) {
		l := &List{}
		l.AddHead(1)
		assertHeadNode(t, l, 1)
	})

	t.Run("adding a head with 1-item list", func(t *testing.T) {
		l := &List{}
		l.AddHead(1)
		l.AddHead(2)

		assertHeadNode(t, l, 2)
		assertTailNode(t, l, 1)
		assertHeadNextNode(t, l, 1)
	})

	t.Run("asserting the tail after adding a new head node to the list", func(t *testing.T) {
		l := &List{}
		l.AddHead(1)
		assertTailNode(t, l, 1)
	})
}

func TestList_AddTail(t *testing.T) {
	t.Run("adding a tail on an empty list", func(t *testing.T) {
		l := &List{}
		l.AddTail(3)
		assertHeadNode(t, l, 3)
		assertTailNode(t, l, 3)
	})

	t.Run("adding a tail on an 1 item list", func(t *testing.T) {
		l := &List{}
		l.AddHead(1)
		l.AddTail(2)
		assertTailNode(t, l, 2)
	})
}

func TestList_PrintOutput(t *testing.T) {
	l := &List{}
	l.AddHead(1)
	l.AddHead(2)
	l.AddHead(3)
	l.AddTail(4)

	out := bytes.Buffer{}
	l.PrintOut(&out)
	got := out.String()
	want := "3\n2\n1\n4\n"
	assert.Equal(t, want, got)
}

func TestList_RemoveHead(t *testing.T) {
	t.Run("After removing a head of a 2-item-sized list", func(t *testing.T) {
		l := &List{}
		l.AddHead(1)
		l.AddHead(2)

		l.RemoveHead()
		assert.Equal(t, 1, l.Size())
	})

	t.Run("Removing a head on a 2-item-sized list", func(t *testing.T) {
		l := &List{}
		l.AddHead(1)
		l.AddHead(2)

		val, _ := l.RemoveHead()
		assert.Equal(t, 2, val)
	})

	t.Run("Removing a head on an empty list", func(t *testing.T) {
		l := &List{}
		_, err := l.RemoveHead()
		assert.Equal(t, ErrEmptyList, err)
	})
}

func TestRemoveNode(t *testing.T) {
	t.Run("after removing a node size should be decrement by 1", func(t *testing.T) {
		l := &List{}
		l.AddHead(1)
		l.AddHead(2)
		l.RemoveNode(2)
		assert.Equal(t, 1, l.Size())
	})

	t.Run("removing the node in the list", func(t *testing.T) {
		l := &List{}
		l.AddHead(1)
		l.AddHead(2)
		l.AddHead(3)
		l.AddTail(10)
		l.RemoveNode(2)
		assert.False(t, l.IsPresent(2))
	})

	t.Run("removing the non-existing node in the list", func(t *testing.T) {
		l := &List{}
		l.AddHead(1)
		l.AddHead(2)
		l.AddHead(3)
		err := l.RemoveNode(4)
		assert.Equal(t, ErrNotFound, err)
	})
}

func TestRemoveNodes(t *testing.T) {
	l := &List{}
	l.AddHead(1)
	l.AddHead(2)
	l.AddHead(2)
	l.AddHead(3)

	l.RemoveNodes(2)
	assert.False(t, l.IsPresent(2))
	l.Print()
}

func TestReverse(t *testing.T) {
	l := &List{}
	l.AddHead(1)
	l.AddHead(2)
	l.AddHead(3)

	l.Reverse()
	var b bytes.Buffer
	l.PrintOut(&b)
	want := "1\n2\n3\n"
	assert.Equal(t, want, b.String())
}

func assertTailNodeIterative(t *testing.T, l *List, want int) {
	t.Helper()
	var gotTail *Node
	l.iterate(func(n *Node) bool {
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
