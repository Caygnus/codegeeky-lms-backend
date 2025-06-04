package auth

import (
	"context"
	"time"

	"github.com/omkar273/codegeeky/internal/auth/abac"
	"github.com/omkar273/codegeeky/internal/auth/rbac"
	"github.com/omkar273/codegeeky/internal/domain/auth"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/types"
)

// UnifiedAuthorizationService provides comprehensive authorization capabilities
// combining both RBAC (Role-Based Access Control) and ABAC (Attribute-Based Access Control)
type UnifiedAuthorizationService interface {
	// Main authorization methods (matching existing interface)
	IsAuthorized(ctx context.Context, request *auth.AccessRequest) (bool, error)
	GetUserPermissions(ctx context.Context, userRole types.UserRole) []auth.Permission
	CheckRolePermission(role types.UserRole, permission auth.Permission) bool
	CheckAttributeBasedAccess(ctx context.Context, request *auth.AccessRequest) (bool, error)

	// Additional ABAC methods
	RegisterABACPolicy(policy abac.Policy) error
	RegisterAttributeProvider(provider abac.AttributeProvider) error

	// Access to underlying services
	GetRBACService() rbac.Service
	GetABACService() abac.Service
}

type unifiedService struct {
	logger      *logger.Logger
	rbacService rbac.Service
	abacService abac.Service
}

// NewUnifiedAuthorizationService creates a new unified authorization service
func NewUnifiedAuthorizationService(logger *logger.Logger) UnifiedAuthorizationService {
	return &unifiedService{
		logger:      logger,
		rbacService: rbac.NewService(logger),
		abacService: abac.NewService(logger),
	}
}

// IsAuthorized is the main authorization method that combines RBAC and ABAC
// Authorization Flow:
// 1. Check RBAC first (fast role-based check)
// 2. If RBAC allows, check ABAC for fine-grained control
// 3. Both must allow for access to be granted
func (s *unifiedService) IsAuthorized(ctx context.Context, request *auth.AccessRequest) (bool, error) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		s.logger.Debugw("Authorization check completed",
			"user_id", request.Subject.UserID,
			"action", request.Action,
			"resource_type", request.Resource.Type,
			"resource_id", request.Resource.ID,
			"duration_ms", duration.Milliseconds())
	}()

	// Step 1: Check RBAC (Role-Based Access Control)
	rbacAllowed := s.rbacService.HasPermission(request.Subject.Role, request.Action)

	if !rbacAllowed {
		s.logger.Infow("Access denied by RBAC",
			"user_id", request.Subject.UserID,
			"role", request.Subject.Role,
			"permission", request.Action,
			"resource_type", request.Resource.Type,
			"resource_id", request.Resource.ID)
		return false, nil
	}

	// Step 2: Check ABAC (Attribute-Based Access Control)
	// Even if RBAC allows, ABAC provides fine-grained control
	abacAllowed, err := s.abacService.Evaluate(ctx, request)
	if err != nil {
		s.logger.Errorw("ABAC evaluation failed",
			"error", err,
			"user_id", request.Subject.UserID,
			"action", request.Action)
		return false, err
	}

	if !abacAllowed {
		s.logger.Infow("Access denied by ABAC",
			"user_id", request.Subject.UserID,
			"role", request.Subject.Role,
			"permission", request.Action,
			"resource_type", request.Resource.Type,
			"resource_id", request.Resource.ID)
		return false, nil
	}

	// Both RBAC and ABAC allow access
	s.logger.Infow("Access granted",
		"user_id", request.Subject.UserID,
		"role", request.Subject.Role,
		"permission", request.Action,
		"resource_type", request.Resource.Type,
		"resource_id", request.Resource.ID)

	return true, nil
}

// GetUserPermissions returns all permissions for a user role (matches existing interface)
func (s *unifiedService) GetUserPermissions(ctx context.Context, userRole types.UserRole) []auth.Permission {
	return s.rbacService.GetUserPermissions(userRole)
}

// CheckRolePermission checks if a role has a specific permission (RBAC only)
func (s *unifiedService) CheckRolePermission(role types.UserRole, permission auth.Permission) bool {
	return s.rbacService.HasPermission(role, permission)
}

// CheckAttributeBasedAccess checks attribute-based access (matches existing interface)
func (s *unifiedService) CheckAttributeBasedAccess(ctx context.Context, request *auth.AccessRequest) (bool, error) {
	return s.abacService.Evaluate(ctx, request)
}

// RegisterABACPolicy adds a new ABAC policy
func (s *unifiedService) RegisterABACPolicy(policy abac.Policy) error {
	return s.abacService.RegisterPolicy(policy)
}

// RegisterAttributeProvider adds an attribute provider for ABAC
func (s *unifiedService) RegisterAttributeProvider(provider abac.AttributeProvider) error {
	return s.abacService.RegisterAttributeProvider(provider)
}

// GetRBACService returns the underlying RBAC service for direct access
func (s *unifiedService) GetRBACService() rbac.Service {
	return s.rbacService
}

// GetABACService returns the underlying ABAC service for direct access
func (s *unifiedService) GetABACService() abac.Service {
	return s.abacService
}

// AuthorizationBuilder provides a fluent interface for complex authorization scenarios
type AuthorizationBuilder struct {
	service *unifiedService
	request *auth.AccessRequest
}

// NewAuthorizationBuilder creates a new authorization builder
func (s *unifiedService) NewAuthorizationBuilder() *AuthorizationBuilder {
	return &AuthorizationBuilder{
		service: s,
		request: &auth.AccessRequest{
			Subject:  &auth.AuthContext{},
			Resource: &auth.Resource{Attributes: make(map[string]interface{})},
			Context:  make(map[string]interface{}),
		},
	}
}

// ForUser sets the user context
func (b *AuthorizationBuilder) ForUser(userID string, role types.UserRole) *AuthorizationBuilder {
	b.request.Subject.UserID = userID
	b.request.Subject.Role = role
	if b.request.Subject.Attributes == nil {
		b.request.Subject.Attributes = make(map[string]interface{})
	}
	return b
}

// OnResource sets the resource being accessed
func (b *AuthorizationBuilder) OnResource(resourceType auth.ResourceType, resourceID string) *AuthorizationBuilder {
	b.request.Resource.Type = resourceType
	b.request.Resource.ID = resourceID
	return b
}

// WithAction sets the action being performed
func (b *AuthorizationBuilder) WithAction(action auth.Permission) *AuthorizationBuilder {
	b.request.Action = action
	return b
}

// WithUserAttribute adds a user attribute
func (b *AuthorizationBuilder) WithUserAttribute(key string, value interface{}) *AuthorizationBuilder {
	b.request.Subject.Attributes[key] = value
	return b
}

// WithResourceAttribute adds a resource attribute
func (b *AuthorizationBuilder) WithResourceAttribute(key string, value interface{}) *AuthorizationBuilder {
	b.request.Resource.Attributes[key] = value
	return b
}

// WithContext adds context information
func (b *AuthorizationBuilder) WithContext(key string, value interface{}) *AuthorizationBuilder {
	b.request.Context[key] = value
	return b
}

// Check performs the authorization check
func (b *AuthorizationBuilder) Check(ctx context.Context) (bool, error) {
	return b.service.IsAuthorized(ctx, b.request)
}

// CheckRBACOnly performs only RBAC check
func (b *AuthorizationBuilder) CheckRBACOnly() bool {
	return b.service.CheckRolePermission(b.request.Subject.Role, b.request.Action)
}

// CheckABACOnly performs only ABAC check
func (b *AuthorizationBuilder) CheckABACOnly(ctx context.Context) (bool, error) {
	return b.service.CheckAttributeBasedAccess(ctx, b.request)
}
