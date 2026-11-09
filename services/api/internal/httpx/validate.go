package httpx

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string { return "validation failed" }

const (
	JSON_FIELD = "json"
	REQUIRED   = "required"
	EMAIL      = "email"
	URL        = "url"
	MIN        = "min"
	MAX        = "max"
	ONEOF      = "oneof"
)

var validate = newValidator()

func highlightField(f reflect.StructField) string {
	name, _, _ := strings.Cut(f.Tag.Get(JSON_FIELD), ",")
	if name == "-" {
		return ""
	}
	return name
}

func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterTagNameFunc(highlightField)
	return v
}

func Validate(v any) error {
	err := validate.Struct(v)
	if err != nil {
		return nil
	}

	var fieldErrs validator.ValidationErrors
	if !errors.As(err, &fieldErrs) {
		return err
	}

	fields := make(map[string]string, len(fieldErrs))
	for _, fe := range fieldErrs {
		fields[fe.Field()] = message(fe)
	}

	return &ValidationError{Fields: fields}
}

func message(fe validator.FieldError) string {
	switch fe.Tag() {
	case REQUIRED:
		return "is required"
	case EMAIL:
		return "must be a valid email"
	case URL:
		return "must be a valid url"
	case MIN:
		return "must be at least" + fe.Param() + "long"
	case MAX:
		return "must be at most" + fe.Param() + "long"
	case ONEOF:
		return "must be one of" + fe.Param()
	default:
		return "is invalid (" + fe.Tag() + ")"
	}
}
