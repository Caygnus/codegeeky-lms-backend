package types

import (
	ierr "github.com/omkar273/codegeeky/internal/errors"
)

type DiscountType string

const (
	DiscountTypeFlat       DiscountType = "flat"
	DiscountTypePercentage DiscountType = "percentage"
)

func (d DiscountType) String() string {
	return string(d)
}

func (d DiscountType) Validate() error {
	switch d {
	case DiscountTypeFlat, DiscountTypePercentage:
		return nil
	default:
		return ierr.NewError("invalid discount type").
			WithHint("discount type must be either flat or percentage").
			Mark(ierr.ErrValidation)
	}
}
