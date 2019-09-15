// +build unit

package generateutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateID(t *testing.T) {
	id := GenerateID("user")
	assert.NotEmpty(t, id)
	//t.Log(id)
}
