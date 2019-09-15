package elasticsearch

import (
	"github.com/jayvib/clean-architecture/log"
	"github.com/jayvib/clean-architecture/pkg/validator"
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
