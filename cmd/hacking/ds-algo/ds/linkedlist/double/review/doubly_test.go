package list

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList_AddHead(t *testing.T) {
	t.Run("Adding head should increment the list length", func(t *testing.T) {
		l := &List{}
		l.AddHead(1)
		assertListLength(t, 1, l)
	})

	t.Run("Adding head in an empty list", func(t *testing.T) {
		l := &List{}
		l.AddHead(2)
		t.Run("head and tail should not nil", func(t *testing.T) {
			assertHeadAndTailNotNil(t, l)
		})

		t.Run("head and tail should contain the added value", func(t *testing.T) {
			assertHeadAndTail(t, 2, 2, l)
		})
	})

	t.Run("Adding head on an non-empty list", func(t *testing.T) {
		l := &List{}
		l.AddHead(1)
		l.AddHead(3)
		l.AddHead(2)
		assertHeadAndTail(t, 2, 1, l)
	})
}

func TestList_AddTail(t *testing.T) {
	t.Run("Adding tail should increment its length", func(t *testing.T) {
		l := &List{}
		l.AddTail(3)
		assertListLength(t, 1, l)
	})

	t.Run("Adding a tail in an empty list", func(t *testing.T) {
		l := &List{}
		l.AddTail(3)
		assertHeadAndTail(t, 3, 3, l)
	})

	t.Run("Adding a tail in a non-empty list", func(t *testing.T) {
		l := &List{}
		l.AddTail(3)
		l.AddTail(2)
		l.AddTail(4)
		assertHeadAndTail(t, 3, 4, l)
	})
}

func TestList_PrintOut(t *testing.T) {
	l := &List{}
	l.AddHead(1)
	l.AddHead(2)
	l.AddHead(3)
	var buff bytes.Buffer
	l.PrintOut(&buff)
	want := "3\n2\n1\n"
	got := buff.String()
	assert.Equal(t, want, got)
}

func TestList_RemoveHead(t *testing.T) {
	t.Run("When list size is more then one", func(t *testing.T) {
		l := &List{}
		l.AddHead(1)
		l.AddHead(2)
		val, _ := l.RemoveHead()
		assert.Equal(t, 2, val)
		assert.Equal(t, 1, l.head.value)
		assert.Nil(t, l.head.previous)
	})

	t.Run("When list size is one", func(t *testing.T) {
		l := &List{}
		l.AddHead(1)
		val, _ := l.RemoveHead()
		assert.Equal(t, 1, val)
		assertNilHeadAndTail(t, l)
	})

	t.Run("Size should decrement by 1", func(t *testing.T) {
		l := &List{}
		l.AddHead(1)
		l.RemoveHead()
		assertListLength(t, 0, l)
	})

	t.Run("When list is empty", func(t *testing.T) {
		l := &List{}
		_, ok := l.RemoveHead()
		assert.False(t, ok)
	})
}

func TestList_RemoveNode(t *testing.T) {
	t.Run("Removing a node should decrement the list size by 1", func(t *testing.T) {
		l := &List{}
		l.AddHead(1)
		l.AddHead(2)
		l.AddHead(3)
		l.RemoveNode(2)
		assertListLength(t, 2, l)
	})

	t.Run("Removing an existing node on an non-empty list", func(t *testing.T) {
		l := &List{}
		l.AddHead(1)
		l.AddHead(2)
		l.AddHead(3)
		l.RemoveNode(2)
		assertValueNotExists(t, 2, l)
		wantOrder := "3\n1\n"
		assertListPrintOut(t, wantOrder, l)
		assertListLength(t, 2, l)
	})

	t.Run("Removing an existing head node on an non-empty list", func(t *testing.T) {
		l := &List{}
		l.AddHead(1)
		l.AddHead(2)
		l.AddHead(3)
		l.RemoveNode(3)
		assertValueNotExists(t, 3, l)
		wantOrder := "2\n1\n"
		assertListPrintOut(t, wantOrder, l)
		assertListLength(t, 2, l)
	})

	t.Run("Removing an existing tail node on an non-empty list", func(t *testing.T) {
		l := &List{}
		l.AddHead(1)
		l.AddHead(2)
		l.AddHead(3)
		l.RemoveNode(1)
		assertValueNotExists(t, 1, l)
		wantOrder := "3\n2\n"
		assertListPrintOut(t, wantOrder, l)
		assertListLength(t, 2, l)
	})

	// TODO: Test when the list is empty
	// TODO: Test when the value is not exists
}

func TestList_Reverse(t *testing.T) {
	l := &List{}
	l.AddHead(1)
	l.AddHead(2)
	l.AddHead(3)
	l.Reverse()
	wantOrder := "1\n2\n3\n"
	assertListPrintOut(t, wantOrder, l)
}

func TestList_Sort(t *testing.T) {
	l := &List{}
	l.AddHead(1)
	l.AddHead(2)
	l.AddHead(3)
	l.Sort()
	wantOrder := "1\n2\n3\n"
	assertListPrintOut(t, wantOrder, l)
}

func assertValueNotExists(t *testing.T, want int, l *List) {
	assert.False(t, l.IsExists(want))
}

func assertListPrintOut(t *testing.T, want string, l *List) {
	var out bytes.Buffer
	l.PrintOut(&out)
	assert.Equal(t, want, out.String())
}

func assertListLength(t *testing.T, want int, l *List) {
	t.Helper()
	assert.Equal(t, want, l.Len())
}

func assertHeadAndTailNotNil(t *testing.T, l *List) {
	t.Helper()
	assert.NotNil(t, l.head)
	assert.NotNil(t, l.tail)
}

func assertNilHeadAndTail(t *testing.T, l *List) {
	t.Helper()
	assert.Nil(t, l.head)
	assert.Nil(t, l.tail)
}

func assertHeadAndTail(t *testing.T, wantHead, wantTail int, l *List) {
	t.Helper()
	gotHead, gotTail := l.head, l.tail
	assert.Equal(t, wantHead, gotHead.value)
	assert.Equal(t, wantTail, gotTail.value)
}
