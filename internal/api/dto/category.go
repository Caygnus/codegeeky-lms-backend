package dto

import (
	"context"

	domainInternship "github.com/omkar273/codegeeky/internal/domain/internship"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/omkar273/codegeeky/internal/validator"
)

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"omitempty"`
	LookupKey   string `json:"lookup_key" binding:"required"`
}

func (c *CreateCategoryRequest) Validate() error {
	if err := validator.ValidateRequest(c); err != nil {
		return ierr.WithError(err).
			WithHint("invalid category create request").
			Mark(ierr.ErrValidation)
	}

	return nil
}

func (c *CreateCategoryRequest) ToCategory(ctx context.Context) *domainInternship.Category {

	return &domainInternship.Category{
		ID:          types.GenerateUUIDWithPrefix(types.UUID_PREFIX_CATEGORY),
		Name:        c.Name,
		Description: c.Description,
		LookupKey:   c.LookupKey,
		BaseModel:   types.GetDefaultBaseModel(ctx),
	}
}

type CategoryResponse struct {
	domainInternship.Category
}

// ListCategoryResponse represents the response for listing categories
type ListCategoryResponse = types.ListResponse[*CategoryResponse]

type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"omitempty"`
	Description string `json:"description" binding:"omitempty"`
	LookupKey   string `json:"lookup_key" binding:"omitempty"`
}

func (c *UpdateCategoryRequest) Validate() error {
	if err := validator.ValidateRequest(c); err != nil {
		return ierr.WithError(err).
			WithHint("invalid category update request").
			Mark(ierr.ErrValidation)
	}

	return nil
}
