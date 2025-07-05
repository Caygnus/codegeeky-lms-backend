package service

import (
	"context"

	"github.com/omkar273/codegeeky/internal/api/dto"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
)

type InternshipBatchService interface {
	Create(ctx context.Context, req *dto.CreateInternshipBatchRequest) (*dto.InternshipBatchResponse, error)
	Get(ctx context.Context, id string) (*dto.InternshipBatchResponse, error)
	Update(ctx context.Context, id string, req *dto.UpdateInternshipBatchRequest) (*dto.InternshipBatchResponse, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter *types.InternshipBatchFilter) (*dto.ListInternshipBatchResponse, error)
}

type internshipBatchService struct {
	ServiceParams
}

func NewInternshipBatchService(
	params ServiceParams,
) InternshipBatchService {
	return &internshipBatchService{
		ServiceParams: params,
	}
}

func (s *internshipBatchService) Create(ctx context.Context, req *dto.CreateInternshipBatchRequest) (*dto.InternshipBatchResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	batch := req.ToInternshipBatch(ctx)

	err := s.InternshipBatchRepo.Create(ctx, batch)
	if err != nil {
		return nil, err
	}

	return &dto.InternshipBatchResponse{
		InternshipBatch: *batch,
	}, nil
}

func (s *internshipBatchService) Get(ctx context.Context, id string) (*dto.InternshipBatchResponse, error) {
	batch, err := s.InternshipBatchRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.InternshipBatchResponse{
		InternshipBatch: *batch,
	}, nil
}

func (s *internshipBatchService) Update(ctx context.Context, id string, req *dto.UpdateInternshipBatchRequest) (*dto.InternshipBatchResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	existingBatch, err := s.InternshipBatchRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		existingBatch.Name = lo.FromPtr(req.Name)
	}
	if req.Description != nil {
		existingBatch.Description = lo.FromPtr(req.Description)
	}
	if req.StartDate != nil {
		existingBatch.StartDate = lo.FromPtr(req.StartDate)
	}
	if req.EndDate != nil {
		existingBatch.EndDate = lo.FromPtr(req.EndDate)
	}
	if req.BatchStatus != nil {
		existingBatch.BatchStatus = lo.FromPtr(req.BatchStatus)
	}

	err = s.InternshipBatchRepo.Update(ctx, existingBatch)
	if err != nil {
		return nil, err
	}

	return &dto.InternshipBatchResponse{
		InternshipBatch: *existingBatch,
	}, nil
}

func (s *internshipBatchService) Delete(ctx context.Context, id string) error {
	_, err := s.InternshipBatchRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	err = s.InternshipBatchRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *internshipBatchService) List(ctx context.Context, filter *types.InternshipBatchFilter) (*dto.ListInternshipBatchResponse, error) {
	if filter == nil {
		filter = types.NewInternshipBatchFilter()
	}

	if err := filter.Validate(); err != nil {
		return nil, err
	}

	count, err := s.InternshipBatchRepo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	batches, err := s.InternshipBatchRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	response := &dto.ListInternshipBatchResponse{
		Items:      make([]*dto.InternshipBatchResponse, count),
		Pagination: types.NewPaginationResponse(count, filter.GetLimit(), filter.GetOffset()),
	}

	for i, batch := range batches {
		response.Items[i] = &dto.InternshipBatchResponse{InternshipBatch: *batch}
	}

	return response, nil
}
