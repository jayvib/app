package elasticsearch

import (
	"github.com/jayvib/app/log"
	"github.com/jayvib/app/pkg/validator"
)

// SetValidatorOpt use a custom validator to be use for the
// Search.
func SetValidatorOpt(v validator.Validator) func(se *Search) {
	return func(se *Search) {
		se.setValidator(v)
	}
}

func SetLoggerOpt(l log.Logger) func(se *Search) {
	return func(se *Search) {
		se.setLogger(l)
	}
}
