package integers

import (
	"errors"
	"fmt"
)

var (
	ErrParameterOdd = errors.New("invalid odd parameter")
)

const (
	OddParameter = "OddParameter"
)

func NewBadParameter(t string, orig error) BadParameter {
	return BadParameter{
		Type: t,
		origErr: orig,
	}
}

type BadParameter struct {
	Type string
	origErr error
}

func (b BadParameter) Error() string {
	return fmt.Sprintf("%s: %s", b.Type, b.origErr)
}

func (b BadParameter) Unwrap() error {
	return b.origErr
}

// Add takes two integers and returns the sum of them.
func Add(x, y int) (int, error) {
	if isOddParam(x, y) {
		return 0, BadParameter{Type: OddParameter, origErr: ErrParameterOdd}
	}
	return x + y, nil
}


func isOddParam(x, y int) bool {
		if x % 2 != 0 {
		return true
	}

	if y & 2 != 0 {
		return true
	}
		return false
}
