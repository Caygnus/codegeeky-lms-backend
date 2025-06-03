package rbac

import (
	"sync"

	"github.com/omkar273/codegeeky/internal/domain/auth"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/types"
)

// Service provides Role-Based Access Control functionality
type Service interface {
	// HasPermission checks if a role has a specific permission
	HasPermission(role types.UserRole, permission auth.Permission) bool

	// GetUserPermissions returns all permissions for a user role
	GetUserPermissions(role types.UserRole) []auth.Permission

	// GetAllRoles returns all available roles
	GetAllRoles() []types.UserRole

	// GetRolePermissions returns the permission mapping for all roles
	GetRolePermissions() map[types.UserRole][]auth.Permission
}

type service struct {
	logger          *logger.Logger
	rolePermissions map[types.UserRole][]auth.Permission
	permissionCache map[string]bool // Cache for faster lookups
	mu              sync.RWMutex    // Protect cache for concurrent access
}

// NewService creates a new RBAC service with predefined role-permission mappings
func NewService(logger *logger.Logger) Service {
	s := &service{
		logger:          logger,
		rolePermissions: initializeRolePermissions(),
		permissionCache: make(map[string]bool),
	}

	// Pre-populate cache for better performance
	s.buildPermissionCache()

	return s
}

// HasPermission checks if a role has a specific permission (O(1) lookup)
func (s *service) HasPermission(role types.UserRole, permission auth.Permission) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Use cache for fast lookup
	cacheKey := string(role) + ":" + string(permission)
	hasPermission, exists := s.permissionCache[cacheKey]

	if exists {
		return hasPermission
	}

	// Fallback to direct lookup (shouldn't happen with proper cache)
	permissions, exists := s.rolePermissions[role]
	if !exists {
		return false
	}

	for _, p := range permissions {
		if p == permission {
			return true
		}
	}

	return false
}

// GetUserPermissions returns all permissions for a user role
func (s *service) GetUserPermissions(role types.UserRole) []auth.Permission {
	s.mu.RLock()
	defer s.mu.RUnlock()

	permissions, exists := s.rolePermissions[role]
	if !exists {
		return []auth.Permission{}
	}

	// Return a copy to prevent external modification
	result := make([]auth.Permission, len(permissions))
	copy(result, permissions)
	return result
}

// GetAllRoles returns all available roles
func (s *service) GetAllRoles() []types.UserRole {
	s.mu.RLock()
	defer s.mu.RUnlock()

	roles := make([]types.UserRole, 0, len(s.rolePermissions))
	for role := range s.rolePermissions {
		roles = append(roles, role)
	}

	return roles
}

// GetRolePermissions returns the complete role-permission mapping
func (s *service) GetRolePermissions() map[types.UserRole][]auth.Permission {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a deep copy to prevent external modification
	result := make(map[types.UserRole][]auth.Permission)
	for role, permissions := range s.rolePermissions {
		permissionsCopy := make([]auth.Permission, len(permissions))
		copy(permissionsCopy, permissions)
		result[role] = permissionsCopy
	}

	return result
}

// buildPermissionCache creates cache entries for all role-permission combinations
func (s *service) buildPermissionCache() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for role, permissions := range s.rolePermissions {
		// Create cache entries for each permission
		for _, permission := range permissions {
			cacheKey := string(role) + ":" + string(permission)
			s.permissionCache[cacheKey] = true
		}
	}

	s.logger.Infow("RBAC permission cache built",
		"total_entries", len(s.permissionCache),
		"roles_count", len(s.rolePermissions))
}

// initializeRolePermissions defines the static role-permission mapping
// This is the core of RBAC - simple, fast, and predictable
func initializeRolePermissions() map[types.UserRole][]auth.Permission {
	return map[types.UserRole][]auth.Permission{
		// ADMIN: Full system access
		types.UserRoleAdmin: {
			// Internship management
			auth.PermissionCreateInternship,
			auth.PermissionUpdateInternship,
			auth.PermissionDeleteInternship,
			auth.PermissionViewInternship,
			auth.PermissionPublishInternship,

			// Content access
			auth.PermissionViewLectures,
			auth.PermissionViewAssignments,
			auth.PermissionViewResources,
			auth.PermissionDownloadContent,

			// User management
			auth.PermissionManageUsers,
			auth.PermissionViewAllUsers,

			// System operations
			auth.PermissionSystemConfig,
			auth.PermissionViewAnalytics,
		},

		// INSTRUCTOR: Content management + limited admin
		types.UserRoleInstructor: {
			// Internship management (limited by ABAC to own content)
			auth.PermissionCreateInternship,
			auth.PermissionUpdateInternship,
			auth.PermissionDeleteInternship,
			auth.PermissionViewInternship,
			auth.PermissionPublishInternship,

			// Content access
			auth.PermissionViewLectures,
			auth.PermissionViewAssignments,
			auth.PermissionViewResources,
			auth.PermissionDownloadContent,

			// Limited analytics (own content only)
			auth.PermissionViewAnalytics,
		},

		// STUDENT: Read-only access to enrolled content
		types.UserRoleStudent: {
			// Limited internship access (filtered by ABAC)
			auth.PermissionViewInternship,

			// Content access (filtered by ABAC for enrolled content only)
			auth.PermissionViewLectures,
			auth.PermissionViewAssignments,
			auth.PermissionViewResources,
			auth.PermissionDownloadContent,
		},
	}
}

// ValidateRole checks if a role is valid
func ValidateRole(role types.UserRole) bool {
	validRoles := []types.UserRole{
		types.UserRoleAdmin,
		types.UserRoleInstructor,
		types.UserRoleStudent,
	}

	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}

	return false
}

// GetRoleHierarchy returns role hierarchy (for future use)
// Admin > Instructor > Student
func GetRoleHierarchy() map[types.UserRole]int {
	return map[types.UserRole]int{
		types.UserRoleAdmin:      3, // Highest level
		types.UserRoleInstructor: 2, // Middle level
		types.UserRoleStudent:    1, // Basic level
	}
}

// IsHigherRole checks if role1 has higher privileges than role2
func IsHigherRole(role1, role2 types.UserRole) bool {
	hierarchy := GetRoleHierarchy()
	return hierarchy[role1] > hierarchy[role2]
}
