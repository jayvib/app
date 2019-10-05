package apperr

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

// String format for Error should be:
// Get: I/O error: network unreachable

func TestError(t *testing.T) {
	t.Run("simple nested error", func(t *testing.T){
		origErr := errors.New("this is an error")
		err :=  Error{Err: origErr, Kind: Other}
		want := `unknown error kind: this is an error`
		assertErrorString(t, err, want)
	})

	t.Run("simple nested error with no Kind", func(t *testing.T){
		origErr := errors.New("this is an error")
		err :=  Error{Err: origErr}
		want := `this is an error`
		assertErrorString(t, err, want)
	})

	t.Run("simple nested error with no Op", func(t *testing.T){
		origErr := errors.New("this is an error")
		err :=  Error{Err: origErr,Kind: Other}
		want := `unknown error kind: this is an error`
		assertErrorString(t, err, want)
	})

	t.Run("simple nested error with Op", func(t *testing.T){
		want := `Get: unknown error kind: this is an error`
		origErr := errors.New("this is an error")
		err :=  Error{Op: Op("Get"), Kind: Other, Err: origErr}
		assertErrorString(t, err, want)
	})

	t.Run("simple nested error witn Op and Kind", func(t *testing.T) {
		want := `Get: I/O error: this is an error`
		origErr := errors.New("this is an error")
		err :=  Error{Op: Op("Get"), Kind: IO, Err: origErr}
		assertErrorString(t, err, want)
		t.Log(err)
	})

	t.Run("nesting one Error", func(t *testing.T){
		want := `Get: I/O error
	Put: I/O error: can't Put item`

		nestedErr := Error{Op: Op("Put"), Kind: IO, Err: errors.New("can't Put item")}

		err :=  Error{Op: Op("Get"), Kind: IO, Err: nestedErr}
		assertErrorString(t, err, want)
		t.Log(err.Error())
	})

	t.Run("no Op value", func(t *testing.T) {
		want := `I/O error
	Put: I/O error: can't Put item`

		nestedErr := Error{Op: Op("Put"), Kind: IO, Err: errors.New("can't Put item")}

		err :=  Error{Kind: IO, Err: nestedErr}
		assertErrorString(t, err, want)
		t.Log(err.Error())
	})

	t.Run("no Kind value", func(t *testing.T){
		want := `Get: 
	Put: I/O error: can't Put item`

		nestedErr := Error{Op: Op("Put"), Kind: IO, Err: errors.New("can't Put item")}

		err :=  Error{Op: Op("Get"), Err: nestedErr}
		assertErrorString(t, err, want)
		t.Log(err.Error())
	})

	t.Run("no nested error", func(t *testing.T) {
		want := `Get: I/O error`

		err :=  Error{Op: Op("Get"), Kind: IO}
		assertErrorString(t, err, want)
		t.Log(err.Error())
	})

	t.Run("empty nested Error type", func(t *testing.T){
		want := `Get: unknown error kind`

		nestedErr := Error{}

		err :=  Error{Op: Op("Get"), Err: nestedErr, Kind: Other}
		assertErrorString(t, err, want)
		t.Log(err.Error())

	})
}

func assertErrorString(t *testing.T,err error, want string) {
	assert.Equal(t, want, err.Error())
}