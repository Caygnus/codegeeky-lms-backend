package service

import (
	"context"

	"github.com/omkar273/codegeeky/internal/api/dto"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
)

type UserService interface {
	Me(ctx context.Context) (*dto.MeResponse, error)
	Update(ctx context.Context, req *dto.UpdateUserRequest) (*dto.MeResponse, error)
}

type userService struct {
	ServiceParams
}

// NewUserService creates a new user service
func NewUserService(params ServiceParams) UserService {
	return &userService{ServiceParams: params}
}

// Me returns the current user
func (s *userService) Me(ctx context.Context) (*dto.MeResponse, error) {
	userID := types.GetUserID(ctx)

	if userID == "" {
		return nil, ierr.WithError(ierr.ErrPermissionDenied).
			WithHint("User not authenticated").
			Mark(ierr.ErrPermissionDenied)
	}

	user, err := s.UserRepo.Get(ctx, userID)
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to get user").
			Mark(ierr.ErrDatabase)
	}
	return &dto.MeResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     string(user.Role),
		Phone:    user.Phone,
	}, nil
}

// Update updates the current user
func (s *userService) Update(ctx context.Context, req *dto.UpdateUserRequest) (*dto.MeResponse, error) {
	userID := types.GetUserID(ctx)

	if userID == "" {
		return nil, ierr.WithError(ierr.ErrPermissionDenied).
			WithHint("User not authenticated").
			Mark(ierr.ErrPermissionDenied)
	}

	user, err := s.UserRepo.Get(ctx, userID)
	if err != nil {
		return nil, ierr.WithError(err).
			WithHint("Failed to get user").
			Mark(ierr.ErrDatabase)
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	err = s.UserRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return &dto.MeResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     string(user.Role),
		Phone:    user.Phone,
	}, nil
}
