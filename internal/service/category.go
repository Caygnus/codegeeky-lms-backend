package service

import (
	"context"

	"github.com/omkar273/codegeeky/internal/api/dto"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
)

type CategoryService interface {
	Create(ctx context.Context, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetByID(ctx context.Context, id string) (*dto.CategoryResponse, error)
	Update(ctx context.Context, id string, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter *types.CategoryFilter) (*dto.ListCategoryResponse, error)
}

type categoryService struct {
	ServiceParams
}

func NewCategoryService(
	params ServiceParams,
) CategoryService {
	return &categoryService{
		ServiceParams: params,
	}
}

func (s *categoryService) Create(ctx context.Context, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {

	if err := req.Validate(); err != nil {
		return nil, err
	}

	category := req.ToCategory(ctx)

	err := s.CategoryRepo.Create(ctx, category)
	if err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		Category: *category,
	}, nil
}

func (s *categoryService) GetByID(ctx context.Context, id string) (*dto.CategoryResponse, error) {
	category, err := s.CategoryRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		Category: *category,
	}, nil
}

func (s *categoryService) Update(ctx context.Context, id string, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {

	if err := req.Validate(); err != nil {
		return nil, err
	}

	category, err := s.CategoryRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// if no changes are made, return the category
	if category.Name == req.Name && category.Description == req.Description && category.LookupKey == req.LookupKey {
		return &dto.CategoryResponse{
			Category: *category,
		}, nil
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.LookupKey != "" {
		category.LookupKey = req.LookupKey
	}

	err = s.CategoryRepo.Update(ctx, category)
	if err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		Category: *category,
	}, nil
}

func (s *categoryService) Delete(ctx context.Context, id string) error {

	if id == "" {
		return ierr.NewError("Category ID is required").
			WithHint("Category ID is required").
			Mark(ierr.ErrValidation)
	}

	_, err := s.CategoryRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	err = s.CategoryRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *categoryService) List(ctx context.Context, filter *types.CategoryFilter) (*dto.ListCategoryResponse, error) {
	if filter == nil {
		filter = types.NewCategoryFilter()
	}

	if err := filter.Validate(); err != nil {
		return nil, err
	}

	count, err := s.CategoryRepo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Get categories
	categories, err := s.CategoryRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Build response
	response := &dto.ListCategoryResponse{
		Items:      make([]*dto.CategoryResponse, count),
		Pagination: types.NewPaginationResponse(count, filter.GetLimit(), filter.GetOffset()),
	}

	// Add items to response
	for i, category := range categories {
		response.Items[i] = &dto.CategoryResponse{Category: *category}
	}

	return response, nil
}
