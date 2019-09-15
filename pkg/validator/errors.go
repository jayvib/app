package validator

import "strings"

type ValidationErr map[string]string // field and the message.

func (v ValidationErr) Error() string {
	var builder strings.Builder
	prefix := "validator:"
	builder.WriteString(prefix)
	builder.WriteString(" Validation Error")

	for _, message := range v {
		builder.WriteString("\n\t")
		builder.WriteString(message)
	}
	return builder.String()
}

func IsValidationErr(err error) bool {
	if _, ok := err.(ValidationErr); ok {
		return true
	}
	return false
}
