package single

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestAddHead(t *testing.T) {
	l := new(List)

	l.AddHead(2)

	gotSize := l.Size()
	wantSize := 1
	assert.Equal(t, wantSize, gotSize)

}
