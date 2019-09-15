// +build unit

package apperr

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddInfos(t *testing.T) {
	err := New(InternalError, "some internal error", nil)
	err = AddInfos(err, "test1", "test1Value", "test2", "teset2Value", "test3", "test3value")
	ae := err.(appError)
	assert.Len(t, ae.extraInfo, 4)
}
