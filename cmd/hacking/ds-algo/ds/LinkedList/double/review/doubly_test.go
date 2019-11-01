package list

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
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

	t.Run("Adding a tail in a non-empty list", func(t *testing.T){
		l := &List{}
		l.AddTail(3)
		l.AddTail(2)
		l.AddTail(4)
		assertHeadAndTail(t, 3, 4, l)
		l.Print()
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
	t.Run("When list size is more then one", func(t *testing.T){
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
