package auth

import (
	"context"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nedpals/supabase-go"
	"github.com/omkar273/codegeeky/internal/api/dto"
	"github.com/omkar273/codegeeky/internal/config"
	"github.com/omkar273/codegeeky/internal/domain/auth"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/security"
	"github.com/omkar273/codegeeky/internal/types"
)

type supabaseProvider struct {
	cfg               *config.Configuration
	supabase          *supabase.Client
	logger            *logger.Logger
	encryptionService security.EncryptionService
}

func NewSupabaseProvider(cfg *config.Configuration, logger *logger.Logger, encryptionService security.EncryptionService) Provider {

	supabaseUrl := cfg.Supabase.URL
	adminApiKey := cfg.Supabase.ServiceKey

	client := supabase.CreateClient(supabaseUrl, adminApiKey)

	if client == nil {
		log.Fatal("failed to create supabase client")
	}

	return &supabaseProvider{
		cfg:               cfg,
		supabase:          client,
		logger:            logger,
		encryptionService: encryptionService,
	}
}

func (p *supabaseProvider) GetProvider() types.AuthProvider {
	return types.AuthProviderSupabase
}

func (p *supabaseProvider) ValidateToken(ctx context.Context, token string) (*auth.Claims, error) {
	// Parse and validate the JWT token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ierr.NewErrorf("unexpected signing method: %v", token.Header["alg"]).
				WithHint("Please use the correct signing method").
				Mark(ierr.ErrValidation)
		}
		// Return the JWT secret as the key for validation
		return []byte(p.cfg.Supabase.JWTSecret), nil
	})

	if err != nil {
		p.logger.Error("Failed to parse JWT token", "error", err)
		return nil, ierr.NewErrorf("invalid token: %w", err).
			WithHint("Please use the correct token").
			Mark(ierr.ErrValidation)
	}

	// Check if token is valid
	if !parsedToken.Valid {
		p.logger.Error("JWT token is invalid")
		return nil, ierr.NewError("token is invalid").
			WithHint("Please use the correct token").
			Mark(ierr.ErrValidation)
	}

	// Extract claims from the token
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		p.logger.Error("Failed to extract claims from JWT token")
		return nil, ierr.NewError("failed to extract claims from token").
			WithHint("Please use the correct token").
			Mark(ierr.ErrValidation)
	}

	// Validate audience - should be "authenticated" for user tokens
	if aud, ok := claims["aud"].(string); ok && aud != "authenticated" {
		p.logger.Error("Invalid audience in JWT token", "audience", aud)
		return nil, ierr.NewError("invalid token audience").
			WithHint("Please use the correct token").
			Mark(ierr.ErrValidation)
	}

	// Validate role - should be "authenticated" for user tokens
	if role, ok := claims["role"].(string); ok && role != "authenticated" {
		p.logger.Error("Invalid role in JWT token", "role", role)
		return nil, ierr.NewError("invalid token role").
			WithHint("Please use the correct token").
			Mark(ierr.ErrValidation)
	}

	// Extract user information from JWT claims
	userID, _ := claims["sub"].(string)
	email, _ := claims["email"].(string)
	phone, _ := claims["phone"].(string)

	// get role from app metadata
	appMetadata, _ := claims["app_metadata"].(map[string]interface{})
	role, _ := appMetadata["role"].(string)

	// Validate that we have at least a user ID
	if userID == "" {
		p.logger.Error("JWT token missing user ID (sub claim)")
		return nil, ierr.NewError("token missing user ID").
			WithHint("Please use the correct token").
			Mark(ierr.ErrValidation)
	}

	p.logger.Debug("Successfully validated JWT token", "user_id", userID, "email", email, "role", role)

	return &auth.Claims{
		UserID: userID,
		Email:  email,
		Phone:  phone,
		Role:   types.UserRole(role),
	}, nil
}

// SignUp is not used directly for Supabase as users sign up through the Supabase UI
// This method is kept for compatibility with the Provider interface
func (p *supabaseProvider) SignUp(ctx context.Context, req *dto.SignupRequest) (*dto.SignupResponse, error) {

	// For Supabase, we don't directly sign up users through this method
	// Instead, we validate the token and get user info
	// For Supabase, we validate the token and extract user info
	if req.AccessToken == "" {
		return nil, ierr.NewError("token is required").
			Mark(ierr.ErrPermissionDenied)
	}

	// Validate the token
	claims, err := p.ValidateToken(ctx, req.AccessToken)
	if err != nil {
		return nil, err
	}

	if claims.Email != req.Email {
		return nil, ierr.NewError("email mismatch").
			Mark(ierr.ErrPermissionDenied)
	}

	if claims.Phone != req.Phone {
		return nil, ierr.NewError("phone mismatch").
			Mark(ierr.ErrPermissionDenied)
	}

	// inject users role in app meta
	err = p.UpdateAppMetadata(ctx, claims.UserID, map[string]interface{}{
		"role": req.Role,
	})
	if err != nil {
		return nil, err
	}

	// Create response
	resp := &dto.SignupResponse{
		ID:          claims.UserID,
		AccessToken: req.AccessToken,
	}

	return resp, nil
}

func (p *supabaseProvider) UpdateAppMetadata(ctx context.Context, userID string, appMetadata map[string]interface{}) error {

	params := supabase.AdminUserParams{
		AppMetadata: appMetadata,
	}

	_, err := p.supabase.Admin.UpdateUser(ctx, userID, params)
	if err != nil {
		return err
	}

	return nil
}
