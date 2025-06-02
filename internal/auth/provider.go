package auth

import (
	"context"

	"github.com/omkar273/codegeeky/internal/api/dto"
	"github.com/omkar273/codegeeky/internal/domain/auth"
	"github.com/omkar273/codegeeky/internal/types"
)

type Provider interface {

	// User Management
	GetProvider() types.AuthProvider
	SignUp(ctx context.Context, req *dto.SignupRequest) (*dto.SignupResponse, error)
	// Login(ctx context.Context, req AuthRequest, userAuthInfo *auth.Auth) (*AuthResponse, error)
	ValidateToken(ctx context.Context, token string) (*auth.Claims, error)
}
