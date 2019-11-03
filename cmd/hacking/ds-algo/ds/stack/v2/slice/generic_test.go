package slice

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestGenericStack_Push(t *testing.T) {
	s := &GenericStack{}
	s.Push("hello")
	assert.Equal(t, 1, s.len())
	assert.Equal(t, "hello", s.v[0])
}

func TestGenericStack_Pop(t *testing.T) {
	t.Run("stack size is 1", func(t *testing.T){
		s := &GenericStack{}
		s.Push("hello")
		v, _ := s.Pop()
		assert.Equal(t, "hello", v)
		assert.Equal(t, 0, s.len())
	})

	t.Run("stack size is more than 1", func(t *testing.T){
		s := &GenericStack{}
		s.Push("hello")
		s.Push("world")
		v, _ := s.Pop()
		assert.Equal(t, "world", v)
		assert.Equal(t, 1, s.len())
	})

	t.Run("stack is empty", func(t *testing.T){
		s := &GenericStack{}
		_, err := s.Pop()
		assert.Equal(t, ErrEmptyStack, err)
	})
}
