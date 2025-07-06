package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/omkar273/codegeeky/internal/auth"
	"github.com/omkar273/codegeeky/internal/config"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/security"
	"github.com/omkar273/codegeeky/internal/types"
)

// setContextValues sets the user ID and user email in the context
func setContextValues(c *gin.Context, userID, userEmail string, userRole types.UserRole) {
	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, types.CtxUserID, userID)
	ctx = context.WithValue(ctx, types.CtxUserEmail, userEmail)
	ctx = context.WithValue(ctx, types.CtxUserRole, userRole)
	c.Request = c.Request.WithContext(ctx)
}

// GuestAuthenticateMiddleware is a middleware that allows requests without authentication
// For now it sets a default user ID and user email in the request context
func GuestAuthenticateMiddleware(c *gin.Context) {
	// TODO: This is a temporary solution to allow requests without authentication.
	// We need to find a better way to handle this.
	c.Set(string(types.CtxIsGuest), true)
	c.Next()
}

// AuthenticateMiddleware is a middleware that authenticates requests based on either:
// 1. JWT token in the Authorization header as a Bearer token
// 2. Development API key in the Authorization header (only in development mode)
func AuthenticateMiddleware(cfg *config.Configuration, logger *logger.Logger) gin.HandlerFunc {

	return func(c *gin.Context) {

		encryptionService, err := security.NewEncryptionService(cfg, logger)
		if err != nil {
			logger.Errorw("failed to create encryption service", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create encryption service"})
			c.Abort()
			return
		}

		authProvider := auth.NewSupabaseProvider(cfg, logger, encryptionService)

		// Get authorization header
		authHeader := c.GetHeader(types.HeaderAuthorization)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Check if it's a development API key (only in development mode)
		if cfg.Server.Env == config.EnvLocal || cfg.Server.Env == config.EnvDev {
			if cfg.Supabase.DevAPIKey != "" && authHeader == "Bearer "+cfg.Supabase.DevAPIKey {
				// Use default development user
				setContextValues(c, types.DefaultUserID, types.DefaultUserEmail, types.DefaultUserRole)
				logger.Debugw("Using development API key for authentication", "user_id", types.DefaultUserID)
				c.Next()
				return
			}
		}

		// Check if the authorization header is in the correct format for JWT
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := authProvider.ValidateToken(c.Request.Context(), tokenString)
		if err != nil {
			logger.Errorw("failed to validate token", "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims == nil || claims.UserID == "" || claims.Email == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		setContextValues(c, claims.UserID, claims.Email, claims.Role)
		c.Next()
	}
}
