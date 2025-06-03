package validator

import (
	"errors"
	"fmt"
	"sync"

	"github.com/go-playground/validator/v10"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/samber/lo"
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

func ValidateEnums[T ~string](field []T, valid []T, fieldName string) error {
	for _, val := range field {
		if !lo.Contains(valid, val) {
			return ierr.NewErrorf("invalid %s: %s", fieldName, val).
				WithReportableDetails(map[string]any{
					"valid_" + fieldName + "s": valid,
				}).
				WithHint(fmt.Sprintf("Invalid %s", fieldName)).
				Mark(ierr.ErrValidation)
		}
	}
	return nil
}
