package auth

import (
	"context"

	"github.com/omkar273/police/internal/domain/auth"
	"github.com/omkar273/police/internal/types"
)

type Provider interface {

	// User Management
	GetProvider() types.AuthProvider
	// SignUp(ctx context.Context, req AuthRequest) (*AuthResponse, error)
	// Login(ctx context.Context, req AuthRequest, userAuthInfo *auth.Auth) (*AuthResponse, error)
	ValidateToken(ctx context.Context, token string) (*auth.Claims, error)
}
