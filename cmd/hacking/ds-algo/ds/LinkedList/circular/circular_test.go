package circular

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPeek(t *testing.T) {
	t.Run("empty list", func(t *testing.T){
		l := &CircularLinkedList{}

		want := 0
		got, err := l.Peek()
		assert.Error(t, err)
		assert.Equal(t, ErrEmptyList, err)
		assert.Equal(t, want, got)
	})
}

func TestAddHead(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := &CircularLinkedList{}
		l.AddHead(3)

		assertCount(t, l, 1)
		assertHeadAndTail(t, l, 3, 3)
	})

	t.Run("there's an existing item in the list", func(t *testing.T) {
		l := &CircularLinkedList{}
		l.AddHead(3)
		l.AddHead(2)

		assertCount(t, l, 2)
		assertHeadAndTail(t, l, 2, 3)
	})

	t.Run("list with exiting 2 items", func(t *testing.T) {
		l := &CircularLinkedList{}
		l.AddHead(3)
		l.AddHead(2)
		l.AddHead(1)

		assertHeadAndTail(t, l, 1, 3)
	})
}

func TestListAddTail(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := &CircularLinkedList{}
		l.AddTail(3)

		assertCount(t, l, 1)
		assertHeadAndTail(t, l, 3, 3)
	})

	t.Run("list with existing 1 item", func(t *testing.T) {
		l := &CircularLinkedList{}
		l.AddTail(3)
		l.AddTail(2)

		assertCount(t, l, 2)
		assertHeadAndTail(t, l, 3, 2)
	})

	t.Run("list with exiting 2 items", func(t *testing.T) {
		l := &CircularLinkedList{}
		l.AddTail(3)
		l.AddTail(2)
		l.AddTail(1)

		assertHeadAndTail(t, l, 3, 1)
	})
}

func TestListRemoveHead(t *testing.T) {
	t.Run("removing head on a 1-item list", func(t *testing.T){
		l := &CircularLinkedList{}
		l.AddHead(1)

		got, _ := l.RemoveHead()

		want := 1
		assertCount(t, l, 0)
		assert.Equal(t, want, got)
	})

	t.Run("after removing the head of an 1-item list when removing head again it should return an error", func(t *testing.T){
		l := &CircularLinkedList{}
		l.AddHead(1)
		l.RemoveHead()
		_, err := l.RemoveHead()
		assert.Equal(t, ErrEmptyList, err)
	})

	t.Run("removing a head on a 2-item list", func(t *testing.T){
		l := &CircularLinkedList{}
		l.AddHead(3)
		l.AddHead(2)
		l.AddHead(1)

		got, _ := l.RemoveHead()

		want := 1

		assert.Equal(t, want, got)

		gotNewHead, _ := l.Peek()
		wantNewHead := 2
		assert.Equal(t, wantNewHead, gotNewHead)
	})

}

func assertCount(t *testing.T, l *CircularLinkedList, want int) {
	t.Helper()
	assert.Equal(t, l.Size(), want)
}

func assertHeadAndTail(t *testing.T, l *CircularLinkedList, wantHead, wantTail int) {
	t.Helper()
	// this is the next of the head
	gotTailValue := l.tail.value
	gotHeadValue, _ := l.Peek()
	assert.Equal(t, wantHead, gotHeadValue)
	assert.Equal(t, wantTail, gotTailValue)
}
