package service

import (
	"context"

	"github.com/omkar273/codegeeky/internal/api/dto"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
)

type DiscountService interface {
	Create(ctx context.Context, req *dto.CreateDiscountRequest) (*dto.DiscountResponse, error)
	GetByID(ctx context.Context, id string) (*dto.DiscountResponse, error)
	Update(ctx context.Context, id string, req *dto.UpdateDiscountRequest) (*dto.DiscountResponse, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter *types.DiscountFilter) (*dto.ListDiscountResponse, error)
	GetByCode(ctx context.Context, code string) (*dto.DiscountResponse, error)
}

type discountService struct {
	ServiceParams
}

func NewDiscountService(params ServiceParams) DiscountService {
	return &discountService{
		ServiceParams: params,
	}
}

func (s *discountService) Create(ctx context.Context, req *dto.CreateDiscountRequest) (*dto.DiscountResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	discount := req.ToDiscount(ctx)

	if err := s.ServiceParams.DiscountRepo.Create(ctx, discount); err != nil {
		return nil, err
	}

	return &dto.DiscountResponse{Discount: *discount}, nil
}

func (s *discountService) GetByID(ctx context.Context, id string) (*dto.DiscountResponse, error) {
	discount, err := s.ServiceParams.DiscountRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.DiscountResponse{Discount: *discount}, nil
}

func (s *discountService) Update(ctx context.Context, id string, req *dto.UpdateDiscountRequest) (*dto.DiscountResponse, error) {
	discount, err := s.ServiceParams.DiscountRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if discount == nil {
		return nil, ierr.NewError("discount not found").
			WithHint("Discount not found").
			Mark(ierr.ErrNotFound)
	}

	if req.Description != "" {
		discount.Description = req.Description
	}

	if req.ValidFrom != nil {
		discount.ValidFrom = *req.ValidFrom
	}

	if req.ValidUntil != nil {
		discount.ValidUntil = req.ValidUntil
	}

	if req.IsActive != nil {
		discount.IsActive = *req.IsActive
	}

	if req.IsCombinable != nil {
		discount.IsCombinable = *req.IsCombinable
	}

	if req.MaxUses != nil {
		discount.MaxUses = req.MaxUses
	}

	if req.MinOrderValue != nil {
		discount.MinOrderValue = req.MinOrderValue
	}

	if req.Metadata != nil {
		discount.Metadata = types.Metadata(*req.Metadata)
	}

	if err := s.ServiceParams.DiscountRepo.Update(ctx, discount); err != nil {
		return nil, err
	}

	return &dto.DiscountResponse{Discount: *discount}, nil
}

func (s *discountService) Delete(ctx context.Context, id string) error {
	return s.ServiceParams.DiscountRepo.Delete(ctx, id)
}

func (s *discountService) List(ctx context.Context, filter *types.DiscountFilter) (*dto.ListDiscountResponse, error) {
	if filter == nil {
		filter = types.NewDiscountFilter()
	}

	if err := filter.Validate(); err != nil {
		return nil, err
	}

	count, err := s.ServiceParams.DiscountRepo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	discounts, err := s.ServiceParams.DiscountRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	response := &dto.ListDiscountResponse{
		Items:      make([]*dto.DiscountResponse, count),
		Pagination: types.NewPaginationResponse(count, filter.GetLimit(), filter.GetOffset()),
	}

	for i, discount := range discounts {
		response.Items[i] = &dto.DiscountResponse{Discount: *discount}
	}

	return response, nil
}

func (s *discountService) GetByCode(ctx context.Context, code string) (*dto.DiscountResponse, error) {
	discount, err := s.ServiceParams.DiscountRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	return &dto.DiscountResponse{Discount: *discount}, nil
}
