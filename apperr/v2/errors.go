package apperr

import (
	"strings"
)

var Separator = "\n\t"

type Op string

type Kind uint8

func (k Kind) String() string {
	switch k {
	case IO:
		return "I/O error"
	case Other:
		return "unknown error kind"
	case Invalid:
		return "invalid operation"
	}
	return "unknown error kind"
}

const (
	Other Kind = iota + 1
	IO
	Invalid
)

// Error is a custom error used throughout the application.
type Error struct {
	Err  error
	Op   Op
	Kind Kind
}

func (e Error) Error() string {
	var b strings.Builder
	if !isEmpty(e.Op) {
		b.WriteString(string(e.Op))
		pad(&b, ": ")
	}
	if !isEmpty(e.Kind) {
		b.WriteString(e.Kind.String())
	}
	if e.Err != nil {
		switch prevErr := e.Err.(type) {
		case Error:
			if !prevErr.isZero() {
				separate(&b, Separator)
				b.WriteString(prevErr.Error())
			}
		default:
			if !isEmpty(e.Kind) {
				pad(&b, ": ")
			}
			b.WriteString(e.Err.Error())
		}
	}
	return b.String()
}

func (e Error) Cause() error {
	return e.Err
}

func (e Error) Unwrap() error {
	return e.Err
}

func (e Error) isZero() bool {
	return isEmpty(e.Op) && isEmpty(e.Kind) && isEmpty(e.Err)
}

func isEmpty(t interface{}) bool {
	const (
		emptyString = ""
		emptyKind   = 0
	)
	switch v := t.(type) {
	case Op:
		if v == emptyString {
			return true
		}
	case Kind:
		if v == emptyKind {
			return true
		}
	case Error:
		if v.isZero() {
			return true
		}
	case error:
		switch e := v.(type) {
		case Error:
			if e.isZero() {
				return true
			}
		default:
			if e.Error() == emptyString {
				return true
			}
		}
	case nil:
		return true
	}
	return false
}

// pad writes 'str' to 'b'.
func pad(b *strings.Builder, str string) {
	b.WriteString(str)
}

func separate(b *strings.Builder, str string) {
	b.WriteString(str)
}
