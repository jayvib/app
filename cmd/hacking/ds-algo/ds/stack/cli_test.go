package stack

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	reader := strings.NewReader("awesome")
	stack := &simpleStack{}

	cli := CLI{stack, reader}
	cli.Run()

	assert.Equal(t, stack.Len(), 1)
}
