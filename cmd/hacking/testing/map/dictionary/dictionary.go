package dictionary

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound  = errors.New("word not found")
	ErrWordExist = errors.New("word exists")
)

type ErrType int

const (
	NotFound ErrType = iota
	WordExist
)

var sep = "::"

type Op string

type Err struct {
	t         ErrType
	operation Op
	origErr   error
}

func(e Err) Error() string {
	var errTypeMessage string
	switch e.t {
	case NotFound:
		errTypeMessage = "word not found"
	case WordExist:
		errTypeMessage = "word already exists"
	default:
		errTypeMessage = "unknown error"
	}

	var format string
	switch e.origErr {
	case nil:
		format = "%s: %s"
		return fmt.Sprintf(format, e.operation, errTypeMessage)
	}

	format = "%s: %s%s%s"
	return fmt.Sprintf(format, e.operation, errTypeMessage, sep, e.origErr)

}

// Dictionary is a key-value type look-up object.
type Dictionary map[string]string

// Search searches for key in the dictionary and returns
// its equivalent definition if found.
func (d Dictionary) Search(key string) (definition string, err error) {
	var ok bool
	definition, ok = d[key]
	if !ok {
		return "", ErrNotFound
	}
	return
}

func (d Dictionary) Add(key, definition string) error {
	const op = "Dictionary.Add"
	_, err := d.Search(key)
	switch err {
	case ErrNotFound:
		d[key] = definition
	case nil:
		return Err{operation: op, t: WordExist}
	default:
		return err
	}
	return nil
}
