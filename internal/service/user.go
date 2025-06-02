package service

import (
	"context"

	"github.com/omkar273/codegeeky/internal/api/dto"
	"github.com/omkar273/codegeeky/internal/domain/user"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
)

type UserService interface {
	Me(ctx context.Context) (*dto.MeResponse, error)
	Update(ctx context.Context, req *dto.UpdateUserRequest) (*dto.MeResponse, error)
}

type userService struct {
	userRepository user.Repository
}

// NewUserService creates a new user service
func NewUserService(userRepository user.Repository) UserService {
	return &userService{userRepository: userRepository}
}

// Me returns the current user
func (s *userService) Me(ctx context.Context) (*dto.MeResponse, error) {
	userID := types.GetUserID(ctx)

	if userID == "" {
		return nil, ierr.WithError(ierr.ErrPermissionDenied).
			WithHint("User not authenticated").
			Mark(ierr.ErrPermissionDenied)
	}

	user, err := s.userRepository.Get(ctx, userID)
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

	user, err := s.userRepository.Get(ctx, userID)
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

	err = s.userRepository.Update(ctx, user)
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
