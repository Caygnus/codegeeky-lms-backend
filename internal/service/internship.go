package service

import (
	"context"

	"github.com/omkar273/codegeeky/internal/api/dto"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
)

type InternshipService interface {
	Create(ctx context.Context, req *dto.CreateInternshipRequest) (*dto.InternshipResponse, error)
	GetByID(ctx context.Context, id string) (*dto.InternshipResponse, error)
	Update(ctx context.Context, id string, req *dto.UpdateInternshipRequest) (*dto.InternshipResponse, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter *types.InternshipFilter) (*dto.ListInternshipResponse, error)
}

type internshipService struct {
	ServiceParams
}

func NewInternshipService(
	params ServiceParams,
) InternshipService {
	return &internshipService{
		ServiceParams: params,
	}
}

func (s *internshipService) Create(ctx context.Context, req *dto.CreateInternshipRequest) (*dto.InternshipResponse, error) {

	if err := req.Validate(); err != nil {
		return nil, err
	}

	internship := req.ToInternship(ctx)

	err := s.InternshipRepo.Create(ctx, internship)
	if err != nil {
		return nil, err
	}

	return &dto.InternshipResponse{
		Internship: *internship,
	}, nil
}

func (s *internshipService) GetByID(ctx context.Context, id string) (*dto.InternshipResponse, error) {

	internship, err := s.InternshipRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if internship == nil {
		return nil, ierr.ErrNotFound
	}

	return &dto.InternshipResponse{
		Internship: *internship,
	}, nil
}

func (s *internshipService) Update(ctx context.Context, id string, req *dto.UpdateInternshipRequest) (*dto.InternshipResponse, error) {

	if err := req.Validate(); err != nil {
		return nil, err
	}

	existingInternship, err := s.InternshipRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		existingInternship.Title = lo.FromPtr(req.Title)
	}
	if req.Description != nil {
		existingInternship.Description = lo.FromPtr(req.Description)
	}
	if req.Level != nil {
		existingInternship.Level = lo.FromPtr(req.Level)
	}
	if req.Mode != nil {
		existingInternship.Mode = lo.FromPtr(req.Mode)
	}
	if req.DurationInWeeks != nil {
		existingInternship.DurationInWeeks = lo.FromPtr(req.DurationInWeeks)
	}
	if req.Currency != nil {
		existingInternship.Currency = lo.FromPtr(req.Currency)
	}
	if req.Price != nil {
		existingInternship.Price = lo.FromPtr(req.Price)
	}
	if req.FlatDiscount != nil {
		existingInternship.FlatDiscount = req.FlatDiscount
	}
	if req.PercentageDiscount != nil {
		existingInternship.PercentageDiscount = req.PercentageDiscount
	}
	if req.Skills != nil {
		existingInternship.Skills = req.Skills
	}
	if req.LearningOutcomes != nil {
		existingInternship.LearningOutcomes = req.LearningOutcomes
	}
	if req.Prerequisites != nil {
		existingInternship.Prerequisites = req.Prerequisites
	}
	if req.Benefits != nil {
		existingInternship.Benefits = req.Benefits
	}

	err = s.InternshipRepo.Update(ctx, existingInternship)
	if err != nil {
		return nil, err
	}

	return &dto.InternshipResponse{
		Internship: *existingInternship,
	}, nil
}

func (s *internshipService) Delete(ctx context.Context, id string) error {

	_, err := s.InternshipRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	err = s.InternshipRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *internshipService) List(ctx context.Context, filter *types.InternshipFilter) (*dto.ListInternshipResponse, error) {

	if filter == nil {
		filter = types.NewInternshipFilter()
	}

	if err := filter.Validate(); err != nil {
		return nil, err
	}

	count, err := s.InternshipRepo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	internships, err := s.InternshipRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	response := &dto.ListInternshipResponse{
		Items:      make([]*dto.InternshipResponse, count),
		Pagination: types.NewPaginationResponse(count, filter.GetLimit(), filter.GetOffset()),
	}

	for i, internship := range internships {
		response.Items[i] = &dto.InternshipResponse{Internship: *internship}
	}

	return response, nil
}
