package service

import (
	"context"

	"github.com/omkar273/codegeeky/internal/api/dto"
	"github.com/omkar273/codegeeky/internal/auth"
	ierr "github.com/omkar273/codegeeky/internal/errors"
)

type AuthService interface {
	Signup(ctx context.Context, req *dto.SignupRequest) (*dto.SignupResponse, error)
}

type authService struct {
	ServiceParams ServiceParams
	AuthProvider  auth.Provider
}

func NewAuthService(params ServiceParams, authProvider auth.Provider) AuthService {
	return &authService{
		ServiceParams: params,
		AuthProvider:  authProvider,
	}
}

func (s *authService) Signup(ctx context.Context, req *dto.SignupRequest) (*dto.SignupResponse, error) {

	// validate access token
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// check if user already exists
	existingUser, err := s.ServiceParams.UserRepo.GetByEmail(ctx, req.Email)
	if err != nil && !ierr.IsNotFound(err) {
		return nil, err
	}

	// if user already exists, return error

	if existingUser != nil {
		return nil, ierr.NewError("user already exists").
			WithHint("User already exists").
			Mark(ierr.ErrAlreadyExists)
	}

	// check using supabase provider
	claims, err := s.AuthProvider.SignUp(ctx, req)
	if err != nil {
		return nil, ierr.NewError("Failed to signup with auth provider").
			WithHint("Failed to signup with auth provider").
			Mark(ierr.ErrDatabase)
	}

	// create user
	err = s.ServiceParams.DB.WithTx(ctx, func(ctx context.Context) error {

		// onboard user
		// it handles the creation of user if not exists
		onboardingService := NewOnboardingService(s.ServiceParams)
		err = onboardingService.Onboard(ctx, &dto.OnboardingRequest{
			SignupRequest: *req,
		})

		if err != nil {
			return err
		}
		s.ServiceParams.Logger.Debugw("onboarded user", "user", claims.ID)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}
