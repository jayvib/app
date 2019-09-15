package validator

import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
	entranslations "gopkg.in/go-playground/validator.v9/translations/en"
)

var defaultValidator Validator

func init() {
	defaultValidator, _ = newDefaultValidator()
}

type Validator interface {
	Var(field interface{}, tag string) error
	Struct(obj interface{}) error
	StructExcept(obj interface{}, fields ...string) error
}

func New() Validator {
	return defaultValidator
}

// Struct validates the required field of an struct.
// If there's an missing field it will return a
// ValidationErr that consist of the validation
// field err occur and its value.
func Struct(obj interface{}) error {
	return defaultValidator.Struct(obj)
}

// StructExcept validate the obj and takes the fields argument
// that will be excluding during the validation checking. It returns
// nil if the validation is successful otherwise it will ValidationErr.
func StructExcept(obj interface{}, fields ...string) error {
	return defaultValidator.StructExcept(obj, fields...)
}

func StructPartial(s interface{}, fields ...string) error {
	return errors.New("not implemented")
}

func VarWithValue(field interface{}, other interface{}, tag string) error {
	return errors.New("not implemented")
}

// Var takes a field and its tag to be validated. It returns
// nil error for a valid field and ValidationErr if there
// an validation exception.
func Var(field interface{}, tag string) error {
	return defaultValidator.Var(field, tag)
}

func newDefaultValidator() (Validator, error) {
	v := validator.New()
	fv := &fallbackValidator{v: v}
	if err := fv.registerTranslation(); err != nil {
		return nil, err
	}
	return fv, nil
}

type fallbackValidator struct {
	v     *validator.Validate
	trans ut.Translator
}

func (d *fallbackValidator) registerTranslation() (err error) {
	enTranslator := en.New()
	uni := ut.New(enTranslator, enTranslator)
	trans, _ := uni.GetTranslator("en")
	d.trans = trans
	if err := entranslations.RegisterDefaultTranslations(d.v, trans); err != nil {
		return err
	}
	return d.v.RegisterTranslation("required", trans, RegisterRequired, TranslateRequired)
}

func (d *fallbackValidator) Struct(obj interface{}) error {
	err := d.v.Struct(obj)
	if err != nil {
		return d.validationErr(err)
	}
	return nil
}

func (d *fallbackValidator) StructExcept(obj interface{}, fields ...string) error {
	err := d.v.StructExcept(obj, fields...)
	if err != nil {
		return d.validationErr(err)
	}
	return nil
}

func (d *fallbackValidator) Var(field interface{}, tag string) error {
	err := d.v.Var(field, tag)
	if err != nil {
		return d.validationErr(err)
	}
	return nil
}

func (d *fallbackValidator) validationErr(err error) error {
	if verr, ok := err.(validator.ValidationErrors); ok {
		validationErr := make(ValidationErr)
		for _, e := range verr {
			validationErr[e.Field()] = e.Translate(d.trans)
		}
		return validationErr
	}
	return err
}
