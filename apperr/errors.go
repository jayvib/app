package apperr

import (
	"errors"
	"fmt"
	"runtime"
)

// Blog for error handling:
// https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html

var (
	ItemNotFound         = errors.New("item not found")
	ItemExist            = errors.New("item already exist")
	BadParamInput        = errors.New("bad parameter input")
	EmailAlreadyExist    = errors.New("email already exist")
	UsernameAlreadyExist = errors.New("username  already exist")
	EmptyItemID          = errors.New("empty id")
)

const (
	InternalError = "InternalError"
	EmptyID       = "EmptyID"
	NoItemFound   = "NoItemFound"
	BadParameter  = "BadParameter"
	ValidationErr = "ValidationError"
)

// Error is a snapshot of the error handling of
// AWS Go SDK
type Error interface {
	error
	Code() string
	Message() string
	OrigErrors() []error
	ExtraInfo() map[string]string
}

func New(code string, message string, origErr error) error {
	err := appError{
		code:    code,
		message: message,
	}
	if origErr != nil {
		err.errs = []error{origErr}
	}

	// Get the information where the error occur.
	// including the filename and the line number.
	_, file, line, _ := runtime.Caller(1)
	err.extraInfo = map[string]string{
		"File": fmt.Sprintf("%s:%d", file, line),
	}
	return err
}

type appError struct {
	code      string
	message   string
	errs      []error
	extraInfo map[string]string
}

func (b appError) Error() string {
	size := len(b.errs)
	if size > 0 {
		return SprintError(b.code, b.message, b.extraInfo, errorList(b.errs))
	}
	return SprintError(b.code, b.message, b.extraInfo, nil)
}

func (b appError) Code() string {
	return b.code
}

func (b appError) Message() string {
	return b.message
}

func (b appError) OrigErrors() []error {
	return b.errs
}

func (b appError) ExtraInfo() map[string]string {
	return b.extraInfo
}

type errorList []error

func (e errorList) Error() string {
	msg := ""
	if l := len(e); l > 0 {
		for i := 0; i < l; i++ {
			msg += e[i].Error()
			if i+1 < l {
				msg += "\n"
			}
		}
	}
	return msg
}

func SprintError(code, message string, extra map[string]string, origErr error) string {
	msg := fmt.Sprintf("%s: %s", code, message)
	if extra != nil {
		for k, v := range extra {
			msg += fmt.Sprintf("\n\t%s: %s", k, v)
		}
	}
	if origErr != nil {
		fmt.Println(origErr)
		msg = fmt.Sprintf("%s\ncaused by %v", msg, origErr)
	}
	return msg
}

func Is(err error) bool {
	if _, ok := err.(Error); ok {
		return true
	}
	return false
}

func AddInfo(err error, key, value string) error {
	if !Is(err) {
		return err
	}
	ae := err.(appError)
	if ae.extraInfo == nil {
		ae.extraInfo = map[string]string{
			key: value,
		}
	}

	ae.extraInfo[key] = value
	return ae
}

func AddInfos(err error, keyValue ...string) error {
	if !Is(err) {
		return err
	}
	ae := err.(appError)
	if ae.extraInfo == nil {
		ae.extraInfo = make(map[string]string)
	}
	for len(keyValue) > 2 {
		curKey := keyValue[0]
		curValue := keyValue[1]
		keyValue = keyValue[2:]
		ae.extraInfo[curKey] = curValue
	}

	if len(keyValue) < 2 {
		return err
	}

	ae.extraInfo[keyValue[0]] = keyValue[1]
	return ae
}
