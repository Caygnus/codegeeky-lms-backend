package service

import (
	"context"

	"github.com/omkar273/codegeeky/internal/api/dto"
	domainInternship "github.com/omkar273/codegeeky/internal/domain/internship"
)

type InternshipService interface {
	Create(ctx context.Context, req *dto.CreateInternshipRequest) (*domainInternship.Internship, error)
}

type internshipService struct {
	internshipRepo domainInternship.InternshipRepository
}
