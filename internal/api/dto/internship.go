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
	Price              decimal.Decimal       `json:"price" binding:"required,gt=0"`
	FlatDiscount       *decimal.Decimal      `json:"flat_discount,omitempty" binding:"omitempty,gt=0"`
	PercentageDiscount *decimal.Decimal      `json:"percentage_discount,omitempty" binding:"omitempty,gt=0,lt=100"`
	CategoryIDs        []string              `json:"category_ids,omitempty"`
}

func (i *CreateInternshipRequest) Validate() error {
	err := validator.ValidateRequest(i)
	if err != nil {
		return ierr.WithError(err).
			WithHint("invalid internship create request").
			Mark(ierr.ErrValidation)
	}

	if i.FlatDiscount != nil && i.PercentageDiscount != nil {
		return ierr.NewError("both flat discount and percentage discount cannot be provided").
			WithHint("please provide only one of flat discount or percentage discount").
			Mark(ierr.ErrValidation)
	}

	// validate level
	if err := validator.ValidateEnums([]types.InternshipLevel{i.Level}, types.InternshipLevels, "level"); err != nil {
		return ierr.WithError(err).
			WithHint("invalid level").
			Mark(ierr.ErrValidation)
	}

	// validate mode
	if err := validator.ValidateEnums([]types.InternshipMode{i.Mode}, types.InternshipModes, "mode"); err != nil {
		return ierr.WithError(err).
			WithHint("invalid mode").
			Mark(ierr.ErrValidation)
	}

	return nil
}

func (i *CreateInternshipRequest) ToInternship(ctx context.Context) *domainInternship.Internship {
	return &domainInternship.Internship{
		ID:                 types.GenerateUUIDWithPrefix(types.UUID_PREFIX_INTERNSHIP),
		Title:              i.Title,
		Description:        i.Description,
		LookupKey:          i.LookupKey,
		Skills:             i.Skills,
		Level:              i.Level,
		Mode:               i.Mode,
		DurationInWeeks:    i.DurationInWeeks,
		LearningOutcomes:   i.LearningOutcomes,
		Prerequisites:      i.Prerequisites,
		Benefits:           i.Benefits,
		Currency:           i.Currency,
		Price:              i.Price,
		FlatDiscount:       lo.FromPtr(i.FlatDiscount),
		PercentageDiscount: lo.FromPtr(i.PercentageDiscount),
		Categories: lo.Map(i.CategoryIDs, func(id string, _ int) *domainInternship.Category {
			return &domainInternship.Category{
				ID: id,
			}
		}),
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
