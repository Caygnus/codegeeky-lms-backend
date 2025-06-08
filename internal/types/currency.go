package types

import (
	"github.com/bojanz/currency"
	ierr "github.com/omkar273/codegeeky/internal/errors"
)

// Currency represents an ISO 4217 currency code like "INR", "USD", etc.
type Currency string

func (c Currency) String() string {
	return string(c)
}

func (c Currency) Validate() error {
	if !currency.IsValid(string(c)) {
		return ierr.NewError("invalid currency").
			WithHint("Currency is invalid").
			WithReportableDetails(map[string]any{
				"currency": c,
			}).
			Mark(ierr.ErrValidation)
	}
	return nil
}


