package service

import (
	"context"

	"github.com/omkar273/codegeeky/internal/api/dto"
)

type OnboardingService interface {
	Onboard(ctx context.Context, req *dto.OnboardingRequest) error
}

type onboardingService struct {
	ServiceParams
}

func NewOnboardingService(params ServiceParams) OnboardingService {
	return &onboardingService{
		ServiceParams: params,
	}
}

func (s *onboardingService) Onboard(ctx context.Context, req *dto.OnboardingRequest) error {

	// validate request
	if err := req.Validate(); err != nil {
		return err
	}

	user := req.ToUser(ctx)

	// create user
	err := s.ServiceParams.UserRepo.Create(ctx, user)
	if err != nil {
		return err
	}

	// NOTE: Currently we are only craeting user during onboarding
	// We will also have functionality to give user credits, goodies, etc.
	return nil
}
