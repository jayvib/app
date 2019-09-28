package stack

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContainer(t *testing.T) {
	spy := &simpleStack{}
	t.Run("Push", func(t *testing.T) {
		container := Container{spy}
		t.Run("pushing string value to the container", func(t *testing.T) {
			input := "awesome"
			container.Push(input)
			assert.Equal(t, 1, container.Len())

			want := input
			assert.Equal(t, want, spy.list[0])
		})
		t.Run("pushing int value to the container", func(t *testing.T) {
			input := 13
			container.Push(input)
			assert.Equal(t, 2, container.Len())
			want := input
			assert.Equal(t, want, spy.list[len(spy.list)-1])
		})
	})

	t.Run("Pop", func(t *testing.T) {
		spy := &simpleStack{list: []interface{}{"awesome"}}
		container := Container{list: spy}
		t.Run("with existing value", func(t *testing.T) {
			got := container.Pop()
			want := "awesome"
			assert.Equal(t, want, got)

			assert.Equal(t, 0, container.Len(), "after popping the only item"+
				"there should no item in the container list")
		})
		t.Run("pushing new item", func(t *testing.T) {
			item := "day"
			container.Push(item)

			got := container.Pop()
			want := item
			assert.Equal(t, want, got)
		})
	})

	t.Run("Top", func(t *testing.T) {
		spy := &simpleStack{list: []interface{}{
			"awesome", "day",
		}}
		container := Container{spy}
		got := container.Top()
		want := "awesome"
		assert.Equal(t, want, got)
	})
}
