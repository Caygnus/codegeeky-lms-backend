package service

import (
	"context"
	"time"

	"github.com/omkar273/codegeeky/internal/api/dto"
	"github.com/omkar273/codegeeky/internal/domain/internship"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/shopspring/decimal"
)

type DiscountService interface {
	Create(ctx context.Context, req *dto.CreateDiscountRequest) (*dto.DiscountResponse, error)
	GetByID(ctx context.Context, id string) (*dto.DiscountResponse, error)
	Update(ctx context.Context, id string, req *dto.UpdateDiscountRequest) (*dto.DiscountResponse, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter *types.DiscountFilter) (*dto.ListDiscountResponse, error)
	GetByCode(ctx context.Context, code string) (*dto.DiscountResponse, error)
	ValidateDiscountCode(ctx context.Context, code string, internship *internship.Internship) error
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

func (s *discountService) ValidateDiscountCode(ctx context.Context, code string, internship *internship.Internship) error {
	if internship == nil {
		return ierr.NewError("internship not found").
			WithHint("Internship not found").
			Mark(ierr.ErrNotFound)
	}

	discount, err := s.ServiceParams.DiscountRepo.GetByCode(ctx, code)
	if err != nil {
		return err
	}

	if discount == nil {
		return ierr.NewError("discount not found").
			WithHint("Discount not found").
			Mark(ierr.ErrNotFound)
	}

	// Check if discount is active
	if !discount.IsActive || discount.Status != types.StatusPublished {
		return ierr.NewError("discount is not active").
			WithHint("Discount is not active").
			Mark(ierr.ErrBadRequest)
	}

	// Check if discount is within valid time range
	now := time.Now()
	if discount.ValidFrom.After(now) {
		return ierr.NewError("discount is not yet valid").
			WithHint("Discount is not yet valid").
			Mark(ierr.ErrBadRequest)
	}

	if discount.ValidUntil != nil && discount.ValidUntil.Before(now) {
		return ierr.NewError("discount has expired").
			WithHint("Discount has expired").
			Mark(ierr.ErrBadRequest)
	}

	// Check minimum order value requirement
	if discount.MinOrderValue != nil && discount.MinOrderValue.GreaterThan(decimal.Zero) {
		if discount.MinOrderValue.GreaterThan(internship.Price) {
			return ierr.NewError("order value does not meet minimum requirement for discount").
				WithHint("Order value does not meet minimum requirement for discount").
				Mark(ierr.ErrBadRequest)
		}
	}

	// TODO: Implement max uses validation by counting actual usage from payments/enrollments
	// For now, we'll skip this validation to avoid the UsedCount error
	// if discount.MaxUses != nil && *discount.MaxUses > 0 {
	//     actualUsage := s.getDiscountUsageCount(ctx, discount.ID)
	//     if actualUsage >= *discount.MaxUses {
	//         return ierr.NewError("discount has reached the maximum number of uses").
	//             WithHint("Discount has reached the maximum number of uses").
	//             Mark(ierr.ErrBadRequest)
	//     }
	// }

	return nil
}
