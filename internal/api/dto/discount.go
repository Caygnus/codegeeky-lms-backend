package dto

import (
	"context"
	"time"

	domainDiscount "github.com/omkar273/codegeeky/internal/domain/discount"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/omkar273/codegeeky/internal/validator"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

type CreateDiscountRequest struct {
	Code          string             `json:"code" validate:"required"`
	Description   string             `json:"description" validate:"omitempty"`
	DiscountType  types.DiscountType `json:"discount_type" validate:"required"`
	DiscountValue decimal.Decimal    `json:"discount_value" validate:"required"`
	ValidFrom     *time.Time         `json:"valid_from" validate:"omitempty"`
	ValidUntil    *time.Time         `json:"valid_until" validate:"omitempty"`
	IsActive      *bool              `json:"is_active" validate:"omitempty"`
	IsCombinable  bool               `json:"is_combinable" validate:"omitempty"`
	MaxUses       *int               `json:"max_uses" validate:"omitempty"`
	MinOrderValue *decimal.Decimal   `json:"min_order_value" validate:"omitempty"`
	Metadata      types.Metadata     `json:"metadata" validate:"omitempty"`
}

func (r *CreateDiscountRequest) Validate() error {
	if err := validator.ValidateRequest(r); err != nil {
		return err
	}

	if r.ValidFrom != nil && r.ValidUntil != nil && r.ValidFrom.After(lo.FromPtr(r.ValidUntil)) {
		return ierr.NewError("valid_from must be before valid_until").
			WithHint("Valid from must be before valid until").
			Mark(ierr.ErrValidation)
	}

	if r.ValidUntil != nil && r.ValidUntil.Before(time.Now()) {
		return ierr.NewError("valid_until must be in the future").
			WithHint("Valid until must be in the future").
			Mark(ierr.ErrValidation)
	}

	if err := r.DiscountType.Validate(); err != nil {
		return r.DiscountType.Validate()
	}

	if r.DiscountValue.LessThan(decimal.Zero) {
		return ierr.NewError("discount_value must be greater than zero").
			WithHint("Discount value must be greater than zero").
			Mark(ierr.ErrValidation)
	}

	if r.MaxUses != nil && lo.FromPtr(r.MaxUses) < 0 {
		return ierr.NewError("max_uses must be greater than zero").
			WithHint("Max uses must be greater than zero").
			Mark(ierr.ErrValidation)
	}

	if r.MinOrderValue != nil && r.MinOrderValue.LessThan(decimal.Zero) {
		return ierr.NewError("min_order_value must be greater than zero").
			WithHint("Min order value must be greater than zero").
			Mark(ierr.ErrValidation)
	}

	return nil
}

func (r *CreateDiscountRequest) ToDiscount(ctx context.Context) *domainDiscount.Discount {
	return &domainDiscount.Discount{
		ID:            types.GenerateUUIDWithPrefix(types.UUID_PREFIX_DISCOUNT),
		Code:          r.Code,
		Description:   r.Description,
		DiscountType:  r.DiscountType,
		DiscountValue: r.DiscountValue,
		ValidFrom:     lo.FromPtr(r.ValidFrom),
		ValidUntil:    r.ValidUntil,
		IsActive:      lo.FromPtr(r.IsActive),
		MaxUses:       r.MaxUses,
		MinOrderValue: r.MinOrderValue,
		IsCombinable:  r.IsCombinable,
		Metadata:      r.Metadata,
		BaseModel:     types.GetDefaultBaseModel(ctx),
	}
}

type UpdateDiscountRequest struct {
	Description   string           `json:"description" validate:"omitempty"`
	ValidFrom     *time.Time       `json:"valid_from" validate:"omitempty"`
	ValidUntil    *time.Time       `json:"valid_until" validate:"omitempty"`
	IsActive      *bool            `json:"is_active" validate:"omitempty"`
	IsCombinable  *bool            `json:"is_combinable" validate:"omitempty"`
	MaxUses       *int             `json:"max_uses" validate:"omitempty"`
	MinOrderValue *decimal.Decimal `json:"min_order_value" validate:"omitempty"`
	Metadata      *types.Metadata  `json:"metadata" validate:"omitempty"`
}

type DiscountResponse struct {
	domainDiscount.Discount
}

type ListDiscountResponse = types.ListResponse[*DiscountResponse]
