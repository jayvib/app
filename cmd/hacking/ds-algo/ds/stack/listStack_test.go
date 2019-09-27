package stack

import (
	"github.com/Workiva/go-datastructures/list"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListStack(t *testing.T){
	t.Run("Push", func(t *testing.T){
		listStack := newEmptyListStack()

		data := "awesome"
		want := data

		listStack.Push(data)
		assert.Equal(t, 1, listStack.Len())
		got, _ := listStack.l.Get(0)
		assert.Equal(t, want, got)
	})

	t.Run("Pop", func(t *testing.T){
		listStack := newEmptyListStack()
		data := "awesome"
		listStack.Push(data)
		data = "day"
		listStack.Push(data)
		got := listStack.Pop()
		want := data
		assert.Equal(t, want, got)
	})
}

func newEmptyListStack() *ListStack {
	return &ListStack{l: list.Empty}
}
