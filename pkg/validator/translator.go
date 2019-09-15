package validator

import (
	"github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
)

func RegisterRequired(ut ut.Translator) error {
	return ut.Add("required", "{0} is a required field", true)
}

func TranslateRequired(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T("required", fe.Field())
	return t
}
