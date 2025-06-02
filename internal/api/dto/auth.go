package dto

import (
	"context"

	"github.com/omkar273/codegeeky/internal/domain/user"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/omkar273/codegeeky/internal/validator"
	"github.com/samber/lo"
)

type SignupRequest struct {
	// basic info
	Email    string `json:"email" validate:"email"`
	Phone    string `json:"phone"`
	FullName string `json:"full_name"`

	// role
	Role types.UserRole `json:"role" validate:"required"`

	// access token
	AccessToken string `json:"access_token" validate:"required"`
}

func (r *SignupRequest) ToUser(ctx context.Context) *user.User {
	return &user.User{
		ID:        types.GenerateUUIDWithPrefix(types.UUID_PREFIX_USER),
		Email:     r.Email,
		Phone:     r.Phone,
		FullName:  r.FullName,
		Role:      r.Role,	
		BaseModel: types.GetDefaultBaseModel(ctx),
	}
}

func (r *SignupRequest) Validate() error {
	if err := validator.ValidateRequest(r); err != nil {
		return err
	}

	// validate role
	if !lo.Contains(types.UserRoles, string(r.Role)) {
		return ierr.NewError("invalid role").
			WithHint("Invalid role").
			Mark(ierr.ErrValidation)
	}

	return nil
}

type SignupResponse struct {
	ID          string `json:"id"`
	AccessToken string `json:"access_token"`
}
