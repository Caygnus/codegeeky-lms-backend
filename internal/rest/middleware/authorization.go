package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/omkar273/codegeeky/internal/auth"
	"github.com/omkar273/codegeeky/internal/auth/rbac"
	"github.com/omkar273/codegeeky/internal/config"
	domainAuth "github.com/omkar273/codegeeky/internal/domain/auth"
	domainUser "github.com/omkar273/codegeeky/internal/domain/user"
	ierr "github.com/omkar273/codegeeky/internal/errors"
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
		if c.GetBool(string(types.CtxIsGuest)) {
			c.Next()
			return
		}

		// Get user context from previous auth middleware
		userID := types.GetUserID(c.Request.Context())
		userEmail := types.GetUserEmail(c.Request.Context())
		userRole := types.GetUserRole(c.Request.Context())

		if userRole == "" {
			logger.Warnw("Authorization failed: no user role in context")
			c.Error(ierr.NewErrorf("no user role in context").
				WithHint("Please login to continue").
				Mark(ierr.ErrUnauthorized))
			return
		}

		if userID == "" {
			logger.Warnw("Authorization failed: no user ID in context")
			c.Error(ierr.NewErrorf("no user ID in context").
				WithHint("Please login to continue").
				Mark(ierr.ErrUnauthorized))
			return
		}

		// Get user details including role
		user, err := userRepo.Get(c.Request.Context(), userID)
		if err != nil {
			logger.Errorw("Failed to get user for authorization", "error", err, "user_id", userID)
			c.Error(ierr.NewErrorf("failed to get user for authorization").
				WithHint("Please try again later").
				Mark(err))
			return
		}

		if user == nil {
			logger.Warnw("User not found for authorization", "user_id", userID)
			c.Error(ierr.NewErrorf("user not found for authorization").
				WithHint("Please login to continue").
				Mark(ierr.ErrUnauthorized))
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
		c.Set(string(types.CtxAuthContext), authContext)
		c.Set(string(types.CtxUser), user)

		c.Next()
	}
}

// RequirePermission creates middleware that requires specific permission
func RequirePermission(permission domainAuth.Permission, resourceType domainAuth.ResourceType) gin.HandlerFunc {
	return func(c *gin.Context) {
		authCtx, exists := GetAuthContext(c)
		if !exists {
			c.Error(ierr.NewErrorf("no auth context in context").
				WithHint("Please login to continue").
				Mark(ierr.ErrUnauthorized))
			return
		}

		// Get authorization service from context (you'll need to set this up in your app)
		authzServiceInterface, exists := c.Get("authz_service")
		if !exists {
			c.Error(ierr.NewErrorf("authorization service not available").
				WithHint("Please try again later").
				Mark(ierr.ErrUnauthorized))
			return
		}

		authzService, ok := authzServiceInterface.(auth.AuthorizationService)
		if !ok {
			c.Error(ierr.NewErrorf("invalid authorization service").
				WithHint("Please try again later").
				Mark(ierr.ErrUnauthorized))
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
			c.Error(ierr.NewErrorf("authorization check failed").
				WithHint("Please try again later").
				Mark(err))
			return
		}

		if !allowed {
			c.Error(ierr.NewErrorf("access denied").
				WithHint("You are not authorized to access this resource").
				Mark(ierr.ErrUnauthorized))
			return
		}

		c.Next()
	}
}

// RequireRole creates middleware that requires specific role(s)
func RequireRole(allowedRoles ...types.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := types.GetUserRole(c.Request.Context())
		if userRole == "" {
			c.Error(ierr.NewErrorf("no user role in context").
				WithHint("Please login to continue").
				Mark(ierr.ErrUnauthorized))
			return
		}

		if !rbac.ValidateRole(userRole) {
			c.Error(ierr.NewErrorf("invalid user role").
				WithHint("Please login to continue").
				Mark(ierr.ErrUnauthorized))
			return
		}

		// Check if user has one of the allowed roles
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				c.Next()
				return
			}
		}

		c.Error(ierr.NewErrorf("insufficient permissions").
			WithHint("You are not authorized to access this resource").
			Mark(ierr.ErrPermissionDenied))
		// NO need to return here, as the error will be handled by the error handler
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

// GetAuthContext is a helper function to get auth context from gin context
func GetAuthContext(c *gin.Context) (*domainAuth.AuthContext, bool) {
	authContext, exists := c.Get(string(types.CtxAuthContext))
	if !exists {
		return nil, false
	}

	authCtx, ok := authContext.(*domainAuth.AuthContext)
	return authCtx, ok
}
