package auth

import "github.com/omkar273/codegeeky/internal/types"

type Claims struct {
	UserID string         `json:"user_id"`
	Email  string         `json:"email"`
	Phone  string         `json:"phone"`
	Role   types.UserRole `json:"role"`
}

// AuthContext represents the current authenticated user's context
type AuthContext struct {
	UserID     string                 `json:"user_id"`
	Email      string                 `json:"email"`
	Phone      string                 `json:"phone"`
	Role       types.UserRole         `json:"role"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// Permission represents a specific permission
type Permission string

const (
	// Internship permissions
	PermissionCreateInternship  Permission = "internship:create"
	PermissionUpdateInternship  Permission = "internship:update"
	PermissionDeleteInternship  Permission = "internship:delete"
	PermissionViewInternship    Permission = "internship:view"
	PermissionPublishInternship Permission = "internship:publish"

	// Content permissions
	PermissionViewLectures    Permission = "content:lectures:view"
	PermissionViewAssignments Permission = "content:assignments:view"
	PermissionViewResources   Permission = "content:resources:view"
	PermissionDownloadContent Permission = "content:download"

	// User management permissions
	PermissionManageUsers  Permission = "users:manage"
	PermissionViewAllUsers Permission = "users:view:all"

	// Admin permissions
	PermissionSystemConfig  Permission = "system:config"
	PermissionViewAnalytics Permission = "analytics:view"
)

// Resource represents a resource being accessed
type Resource struct {
	Type       string                 `json:"type"`
	ID         string                 `json:"id,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// AccessRequest represents an access control request
type AccessRequest struct {
	Subject  *AuthContext           `json:"subject"`
	Resource *Resource              `json:"resource"`
	Action   Permission             `json:"action"`
	Context  map[string]interface{} `json:"context,omitempty"`
}
