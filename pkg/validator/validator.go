package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Single validation failure.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []FieldError

func (e ValidationErrors) Error() string {
	msgs := make([]string, len(e))
	for i, fe := range e {
		msgs[i] = fmt.Sprintf("%s: %s", fe.Field, fe.Message)
	}
	return strings.Join(msgs, "; ")
}

var validate = validator.New()

func init() {
	// Use JSON tag names in error messages instead of struct field names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "" || name == "-" {
			return ""
		}
		return strings.Split(name, ",")[0]
	})
}

// Runs struct validation and returns typed FieldErrors or nil.
func Validate(s any) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}
	var errs ValidationErrors
	for _, ve := range err.(validator.ValidationErrors) {
		errs = append(errs, FieldError{
			Field:   ve.Field(),
			Message: humanize(ve),
		})
	}
	return errs
}

func humanize(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s characters", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", fe.Param())
	case "gt":
		return fmt.Sprintf("must be greater than %s", fe.Param())
	case "gte":
		return fmt.Sprintf("must be greater than or equal to %s", fe.Param())
	case "lte":
		return fmt.Sprintf("must be less than or equal to %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("must be one of: %s", strings.ReplaceAll(fe.Param(), " ", ", "))
	case "uuid4":
		return "must be a valid UUID"
	default:
		return fmt.Sprintf("failed validation: %s", fe.Tag())
	}
}
