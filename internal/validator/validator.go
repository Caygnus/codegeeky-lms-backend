package validator

import (
	"errors"
	"sync"

	"github.com/go-playground/validator/v10"
	ierr "github.com/omkar273/police/internal/errors"
)

var (
	validate *validator.Validate
	once     sync.Once
)

// initValidator initializes the validator exactly once
func initValidator() {
	once.Do(func() {
		validate = validator.New()
	})
}

func NewValidator() *validator.Validate {
	initValidator()
	return validate
}

func GetValidator() *validator.Validate {
	initValidator()
	return validate
}

func ValidateRequest(req interface{}) error {
	initValidator()

	if err := validate.Struct(req); err != nil {
		details := make(map[string]any)
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			for _, err := range validateErrs {
				details[err.Field()] = err.Error()
			}
		}
		return ierr.WithError(err).
			WithHint("Request validation failed").
			WithReportableDetails(details).
			Mark(ierr.ErrValidation)
	}
	return nil
}
