package selection_sort

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSort(t *testing.T) {
	input := []int{3, 2, 8, 9, 7, 6, 1, 0}
	Sort(input)

	want := []int{0, 1, 2, 3, 6, 7 ,8, 9}
	assert.Equal(t, want, input)
}
