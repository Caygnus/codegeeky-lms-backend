package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/omkar273/codegeeky/internal/auth"
	"github.com/omkar273/codegeeky/internal/config"
	domainAuth "github.com/omkar273/codegeeky/internal/domain/auth"
	domainUser "github.com/omkar273/codegeeky/internal/domain/user"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/types"
)

// AuthorizationMiddleware creates middleware that checks user permissions
func AuthorizationMiddleware(
	cfg *config.Configuration,
	logger *logger.Logger,
	authzService auth.AuthorizationService,
	userRepo domainUser.Repository,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authorization for guest endpoints
		if c.GetBool("is_guest") {
			c.Next()
			return
		}

		// Get user context from previous auth middleware
		userID := getUserIDFromContext(c.Request.Context())
		userEmail := getUserEmailFromContext(c.Request.Context())
		userRole := getUserRoleFromContext(c.Request.Context())

		if userRole == "" {
			logger.Warnw("Authorization failed: no user role in context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		if userID == "" {
			logger.Warnw("Authorization failed: no user ID in context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// Get user details including role
		user, err := userRepo.Get(c.Request.Context(), userID)
		if err != nil {
			logger.Errorw("Failed to get user for authorization", "error", err, "user_id", userID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Authorization check failed"})
			c.Abort()
			return
		}

		if user == nil {
			logger.Warnw("User not found for authorization", "user_id", userID)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Create auth context
		authContext := &domainAuth.AuthContext{
			UserID:     userID,
			Email:      userEmail,
			Phone:      user.Phone,
			Role:       user.Role,
			Attributes: map[string]interface{}{
				// Add any additional user attributes here
				// These could be loaded from database or other sources
			},
		}

		// Store auth context in gin context for use by handlers
		c.Set("auth_context", authContext)
		c.Set("user_role", user.Role)
		c.Set("user", user)

		c.Next()
	}
}

// RequirePermission creates middleware that requires specific permission
func RequirePermission(permission domainAuth.Permission, resourceType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authContext, exists := c.Get("auth_context")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		authCtx, ok := authContext.(*domainAuth.AuthContext)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid auth context"})
			c.Abort()
			return
		}

		// Get authorization service from context (you'll need to set this up in your app)
		authzServiceInterface, exists := c.Get("authz_service")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Authorization service not available"})
			c.Abort()
			return
		}

		authzService, ok := authzServiceInterface.(auth.AuthorizationService)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid authorization service"})
			c.Abort()
			return
		}

		// Create access request
		resource := &domainAuth.Resource{
			Type:       resourceType,
			Attributes: make(map[string]interface{}),
		}

		// Extract resource ID from URL params if available
		if resourceID := c.Param("id"); resourceID != "" {
			resource.ID = resourceID
			resource.Attributes["id"] = resourceID
		}

		// Extract internship ID from URL params if available
		if internshipID := c.Param("internship_id"); internshipID != "" {
			resource.Attributes["internship_id"] = internshipID
		}

		accessRequest := &domainAuth.AccessRequest{
			Subject:  authCtx,
			Resource: resource,
			Action:   permission,
			Context:  make(map[string]interface{}),
		}

		// Check authorization
		allowed, err := authzService.IsAuthorized(c.Request.Context(), accessRequest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Authorization check failed"})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole creates middleware that requires specific role(s)
func RequireRole(allowedRoles ...types.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		role, ok := userRole.(types.UserRole)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user role"})
			c.Abort()
			return
		}

		// Check if user has one of the allowed roles
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}

// RequireInstructorOrAdmin is a convenience middleware for instructor/admin only endpoints
func RequireInstructorOrAdmin() gin.HandlerFunc {
	return RequireRole(types.UserRoleInstructor, types.UserRoleAdmin)
}

// RequireAdmin is a convenience middleware for admin only endpoints
func RequireAdmin() gin.HandlerFunc {
	return RequireRole(types.UserRoleAdmin)
}

// Helper functions to extract values from context
func getUserIDFromContext(ctx context.Context) string {
	if userID := ctx.Value(types.CtxUserID); userID != nil {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

func getUserEmailFromContext(ctx context.Context) string {
	if userEmail := ctx.Value(types.CtxUserEmail); userEmail != nil {
		if email, ok := userEmail.(string); ok {
			return email
		}
	}
	return ""
}

func getUserRoleFromContext(ctx context.Context) types.UserRole {
	if userRole := ctx.Value(types.CtxUserRole); userRole != nil {
		if role, ok := userRole.(types.UserRole); ok {
			return role
		}
	}
	return ""
}

// GetAuthContext is a helper function to get auth context from gin context
func GetAuthContext(c *gin.Context) (*domainAuth.AuthContext, bool) {
	authContext, exists := c.Get("auth_context")
	if !exists {
		return nil, false
	}

	authCtx, ok := authContext.(*domainAuth.AuthContext)
	return authCtx, ok
}

// GetCurrentUser is a helper function to get current user from gin context
func GetCurrentUser(c *gin.Context) (*domainUser.User, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	u, ok := user.(*domainUser.User)
	return u, ok
}
