package dto

import (
	"context"

	domainInternship "github.com/omkar273/codegeeky/internal/domain/internship"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/omkar273/codegeeky/internal/validator"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

// InternshipCreateDTO is used for creating a new internship via API request
type CreateInternshipRequest struct {
	Title              string                `json:"title" binding:"required,min=3,max=255"`
	Description        string                `json:"description" binding:"required,min=10"`
	LookupKey          string                `json:"lookup_key" binding:"required"`
	Skills             []string              `json:"skills,omitempty"`
	Level              types.InternshipLevel `json:"level" binding:"required"`
	Mode               types.InternshipMode  `json:"mode" binding:"required"`
	DurationInWeeks    int                   `json:"duration_in_weeks,omitempty" binding:"gte=0"`
	LearningOutcomes   []string              `json:"learning_outcomes,omitempty"`
	Prerequisites      []string              `json:"prerequisites,omitempty"`
	Benefits           []string              `json:"benefits,omitempty"`
	Currency           string                `json:"currency" binding:"required,len=3"`
	Price              decimal.Decimal       `json:"price" binding:"required"`
	FlatDiscount       *decimal.Decimal      `json:"flat_discount,omitempty" binding:"omitempty"`
	PercentageDiscount *decimal.Decimal      `json:"percentage_discount,omitempty" binding:"omitempty"`
	CategoryIDs        []string              `json:"category_ids,omitempty"`
}

func (req *CreateInternshipRequest) Validate() error {
	err := validator.ValidateRequest(req)
	if err != nil {
		return ierr.WithError(err).
			WithHint("invalid internship create request").
			Mark(ierr.ErrValidation)
	}

	if req.FlatDiscount != nil && req.PercentageDiscount != nil {
		return ierr.NewError("both flat discount and percentage discount cannot be provided").
			WithHint("please provide only one of flat discount or percentage discount").
			Mark(ierr.ErrValidation)
	}

	// validate level
	if err := req.Level.Validate(); err != nil {
		return ierr.WithError(err).
			WithHint("invalid level").
			Mark(ierr.ErrValidation)
	}

	// validate mode
	if err := req.Mode.Validate(); err != nil {
		return ierr.WithError(err).
			WithHint("invalid mode").
			Mark(ierr.ErrValidation)
	}

	if req.Price.LessThan(decimal.Zero) {
		return ierr.NewError("price cannot be less than zero").
			WithHint("please provide a valid price").
			Mark(ierr.ErrValidation)
	}

	if req.FlatDiscount != nil && req.FlatDiscount.LessThan(decimal.Zero) {
		return ierr.NewError("flat discount cannot be less than zero").
			WithHint("please provide a valid flat discount").
			Mark(ierr.ErrValidation)
	}

	if req.PercentageDiscount != nil && req.PercentageDiscount.LessThan(decimal.Zero) {
		return ierr.NewError("percentage discount cannot be less than zero").
			WithHint("please provide a valid percentage discount").
			Mark(ierr.ErrValidation)
	}

	return nil
}

func (req *CreateInternshipRequest) ToInternship(ctx context.Context) *domainInternship.Internship {
	total := req.Price
	if req.FlatDiscount != nil {
		total = total.Sub(lo.FromPtr(req.FlatDiscount))
	}
	if req.PercentageDiscount != nil {
		discountAmount := total.Mul(lo.FromPtr(req.PercentageDiscount)).Div(decimal.NewFromInt(100))
		total = total.Sub(discountAmount)
	}

	if total.LessThan(decimal.Zero) {
		total = decimal.Zero
	}

	return &domainInternship.Internship{
		ID:                 types.GenerateUUIDWithPrefix(types.UUID_PREFIX_INTERNSHIP),
		Title:              req.Title,
		Description:        req.Description,
		LookupKey:          req.LookupKey,
		Skills:             req.Skills,
		Level:              req.Level,
		Mode:               req.Mode,
		DurationInWeeks:    req.DurationInWeeks,
		LearningOutcomes:   req.LearningOutcomes,
		Prerequisites:      req.Prerequisites,
		Benefits:           req.Benefits,
		Currency:           req.Currency,
		Price:              req.Price,
		FlatDiscount:       req.FlatDiscount,
		PercentageDiscount: req.PercentageDiscount,
		Categories: lo.Map(req.CategoryIDs, func(id string, _ int) *domainInternship.Category {
			return &domainInternship.Category{
				ID: id,
			}
		}),
		Subtotal:  req.Price,
		Total:     total,
		BaseModel: types.GetDefaultBaseModel(ctx),
	}
}

type InternshipResponse struct {
	domainInternship.Internship
}

func (i *InternshipResponse) FromDomain(internship *domainInternship.Internship) *InternshipResponse {
	return &InternshipResponse{
		Internship: *internship,
	}
}

// UpdateInternshipRequest is used for updating an existing internship via API request
type UpdateInternshipRequest struct {
	Title              *string                `json:"title,omitempty" binding:"omitempty,min=3,max=255"`
	Description        *string                `json:"description,omitempty" binding:"omitempty,min=10"`
	LookupKey          *string                `json:"lookup_key,omitempty"`
	Skills             []string               `json:"skills,omitempty"`
	Level              *types.InternshipLevel `json:"level,omitempty"`
	Mode               *types.InternshipMode  `json:"mode,omitempty"`
	DurationInWeeks    *int                   `json:"duration_in_weeks,omitempty" binding:"omitempty,gte=0"`
	LearningOutcomes   []string               `json:"learning_outcomes,omitempty"`
	Prerequisites      []string               `json:"prerequisites,omitempty"`
	Benefits           []string               `json:"benefits,omitempty"`
	Currency           *string                `json:"currency,omitempty" binding:"omitempty,len=3"`
	Price              *decimal.Decimal       `json:"price,omitempty" binding:"omitempty"`
	FlatDiscount       *decimal.Decimal       `json:"flat_discount,omitempty" binding:"omitempty"`
	PercentageDiscount *decimal.Decimal       `json:"percentage_discount,omitempty" binding:"omitempty"`
	CategoryIDs        []string               `json:"category_ids,omitempty"`
}

func (i *UpdateInternshipRequest) Validate() error {
	err := validator.ValidateRequest(i)
	if err != nil {
		return ierr.WithError(err).
			WithHint("invalid internship update request").
			Mark(ierr.ErrValidation)
	}

	if i.FlatDiscount != nil && i.PercentageDiscount != nil {
		return ierr.NewError("both flat discount and percentage discount cannot be provided").
			WithHint("please provide only one of flat discount or percentage discount").
			Mark(ierr.ErrValidation)
	}

	// validate level if provided
	if i.Level != nil {
		if err := validator.ValidateEnums([]types.InternshipLevel{*i.Level}, types.InternshipLevels, "level"); err != nil {
			return ierr.WithError(err).
				WithHint("invalid level").
				Mark(ierr.ErrValidation)
		}
	}

	// validate mode if provided
	if i.Mode != nil {
		if err := validator.ValidateEnums([]types.InternshipMode{*i.Mode}, types.InternshipModes, "mode"); err != nil {
			return ierr.WithError(err).
				WithHint("invalid mode").
				Mark(ierr.ErrValidation)
		}
	}

	return nil
}

type ListInternshipResponse = types.ListResponse[*InternshipResponse]
