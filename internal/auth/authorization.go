package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/omkar273/codegeeky/internal/domain/auth"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/types"
)

// AuthorizationService provides authorization capabilities
type AuthorizationService interface {
	// IsAuthorized checks if a user is authorized to perform an action on a resource
	IsAuthorized(ctx context.Context, request *auth.AccessRequest) (bool, error)

	// GetUserPermissions returns all permissions for a user
	GetUserPermissions(ctx context.Context, userRole types.UserRole) []auth.Permission

	// CheckRolePermission checks if a role has a specific permission (RBAC)
	CheckRolePermission(role types.UserRole, permission auth.Permission) bool

	// CheckAttributeBasedAccess checks attribute-based access (ABAC)
	CheckAttributeBasedAccess(ctx context.Context, request *auth.AccessRequest) (bool, error)
}

type authorizationService struct {
	logger          *logger.Logger
	rolePermissions map[types.UserRole][]auth.Permission
	abacPolicies    []ABACPolicy
}

// ABACPolicy represents an attribute-based access control policy
type ABACPolicy interface {
	Evaluate(ctx context.Context, request *auth.AccessRequest) (bool, error)
	GetName() string
}

// NewAuthorizationService creates a new authorization service
func NewAuthorizationService(logger *logger.Logger) AuthorizationService {
	service := &authorizationService{
		logger:          logger,
		rolePermissions: initializeRolePermissions(),
		abacPolicies:    initializeABACPolicies(),
	}

	return service
}

// IsAuthorized is the main authorization method that combines RBAC and ABAC
func (s *authorizationService) IsAuthorized(ctx context.Context, request *auth.AccessRequest) (bool, error) {
	// First check RBAC - role-based permissions
	if s.CheckRolePermission(request.Subject.Role, request.Action) {
		s.logger.Infow("Access granted via RBAC",
			"user_id", request.Subject.UserID,
			"role", request.Subject.Role,
			"permission", request.Action,
			"resource_type", request.Resource.Type,
			"resource_id", request.Resource.ID)

		// Even if RBAC allows, check ABAC for fine-grained control
		abacResult, err := s.CheckAttributeBasedAccess(ctx, request)
		if err != nil {
			s.logger.Errorw("ABAC check failed", "error", err, "user_id", request.Subject.UserID)
			return false, ierr.NewErrorf("ABAC check failed").Mark(err)
		}

		return abacResult, nil
	}

	// If RBAC denies, still check ABAC (for special cases)
	abacResult, err := s.CheckAttributeBasedAccess(ctx, request)
	if err != nil {
		s.logger.Errorw("ABAC check failed", "error", err, "user_id", request.Subject.UserID)
		return false, ierr.NewErrorf("ABAC check failed").Mark(err)
	}

	if !abacResult {
		s.logger.Warnw("Access denied",
			"user_id", request.Subject.UserID,
			"role", request.Subject.Role,
			"permission", request.Action,
			"resource_type", request.Resource.Type,
			"resource_id", request.Resource.ID)
	}

	return abacResult, nil
}

// GetUserPermissions returns all permissions for a user role
func (s *authorizationService) GetUserPermissions(ctx context.Context, userRole types.UserRole) []auth.Permission {
	return s.rolePermissions[userRole]
}

// CheckRolePermission checks if a role has a specific permission (RBAC)
func (s *authorizationService) CheckRolePermission(role types.UserRole, permission auth.Permission) bool {
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

// CheckAttributeBasedAccess checks all ABAC policies
func (s *authorizationService) CheckAttributeBasedAccess(ctx context.Context, request *auth.AccessRequest) (bool, error) {
	// If no ABAC policies, allow (RBAC has already been checked)
	if len(s.abacPolicies) == 0 {
		return true, nil
	}

	// Evaluate all ABAC policies
	for _, policy := range s.abacPolicies {
		allowed, err := policy.Evaluate(ctx, request)
		if err != nil {
			s.logger.Errorw("ABAC policy evaluation failed",
				"policy", policy.GetName(),
				"error", err,
				"user_id", request.Subject.UserID)
			continue // Continue with other policies instead of failing completely
		}

		if !allowed {
			s.logger.Debugw("ABAC policy denied access",
				"policy", policy.GetName(),
				"user_id", request.Subject.UserID,
				"resource_type", request.Resource.Type,
				"resource_id", request.Resource.ID)
			return false, nil
		}
	}

	return true, nil
}

// initializeRolePermissions sets up the role-based permissions (RBAC)
func initializeRolePermissions() map[types.UserRole][]auth.Permission {
	return map[types.UserRole][]auth.Permission{
		types.UserRoleAdmin: {
			// All permissions for admin
			auth.PermissionCreateInternship,
			auth.PermissionUpdateInternship,
			auth.PermissionDeleteInternship,
			auth.PermissionViewInternship,
			auth.PermissionPublishInternship,
			auth.PermissionViewLectures,
			auth.PermissionViewAssignments,
			auth.PermissionViewResources,
			auth.PermissionDownloadContent,
			auth.PermissionManageUsers,
			auth.PermissionViewAllUsers,
			auth.PermissionSystemConfig,
			auth.PermissionViewAnalytics,
		},
		types.UserRoleInstructor: {
			// Instructor permissions
			auth.PermissionCreateInternship,
			auth.PermissionUpdateInternship,
			auth.PermissionViewInternship,
			auth.PermissionPublishInternship,
			auth.PermissionViewLectures,
			auth.PermissionViewAssignments,
			auth.PermissionViewResources,
			auth.PermissionDownloadContent,
			auth.PermissionViewAnalytics, // Can view their own analytics
		},
		types.UserRoleStudent: {
			// Student permissions - very limited
			auth.PermissionViewInternship, // Can view internships they're enrolled in
		},
	}
}

// initializeABACPolicies sets up attribute-based access control policies
func initializeABACPolicies() []ABACPolicy {
	return []ABACPolicy{
		&EnrollmentBasedAccessPolicy{},
		&ContentAccessPolicy{},
		&TimeBasedAccessPolicy{},
		&OwnershipPolicy{},
	}
}

// EnrollmentBasedAccessPolicy - Students can only access content of internships they're enrolled in
type EnrollmentBasedAccessPolicy struct{}

func (p *EnrollmentBasedAccessPolicy) GetName() string {
	return "EnrollmentBasedAccess"
}

func (p *EnrollmentBasedAccessPolicy) Evaluate(ctx context.Context, request *auth.AccessRequest) (bool, error) {
	// This policy applies to students accessing content
	if request.Subject.Role != types.UserRoleStudent {
		return true, nil // Policy doesn't apply to non-students
	}

	// For content access permissions
	contentPermissions := []auth.Permission{
		auth.PermissionViewLectures,
		auth.PermissionViewAssignments,
		auth.PermissionViewResources,
		auth.PermissionDownloadContent,
	}

	isContentAccess := false
	for _, perm := range contentPermissions {
		if request.Action == perm {
			isContentAccess = true
			break
		}
	}

	if !isContentAccess {
		return true, nil // Policy doesn't apply
	}

	// Check if student is enrolled in the internship
	internshipID, exists := request.Resource.Attributes["internship_id"]
	if !exists {
		return false, fmt.Errorf("internship_id not found in resource attributes")
	}

	// Get user's enrolled internships from context or attributes
	enrolledInternships, exists := request.Subject.Attributes["enrolled_internships"]
	if !exists {
		return false, nil // Student not enrolled in any internships
	}

	enrolledList, ok := enrolledInternships.([]string)
	if !ok {
		return false, fmt.Errorf("invalid enrolled_internships format")
	}

	internshipIDStr, ok := internshipID.(string)
	if !ok {
		return false, fmt.Errorf("invalid internship_id format")
	}

	// Check if student is enrolled in this internship
	for _, enrolledID := range enrolledList {
		if enrolledID == internshipIDStr {
			return true, nil
		}
	}

	return false, nil // Student not enrolled in this internship
}

// ContentAccessPolicy - Controls access to different types of content based on progress
type ContentAccessPolicy struct{}

func (p *ContentAccessPolicy) GetName() string {
	return "ContentAccess"
}

func (p *ContentAccessPolicy) Evaluate(ctx context.Context, request *auth.AccessRequest) (bool, error) {
	// This policy applies to content access
	if request.Resource.Type != "content" {
		return true, nil
	}

	// Admins and instructors have full access
	if request.Subject.Role == types.UserRoleAdmin || request.Subject.Role == types.UserRoleInstructor {
		return true, nil
	}

	// For students, check if they have the prerequisite progress
	requiredProgress, exists := request.Resource.Attributes["required_progress"]
	if !exists {
		return true, nil // No progress requirement
	}

	userProgress, exists := request.Subject.Attributes["progress"]
	if !exists {
		return false, nil // User has no progress recorded
	}

	// Compare progress (this is a simplified example)
	userProgressFloat, ok1 := userProgress.(float64)
	requiredProgressFloat, ok2 := requiredProgress.(float64)

	if !ok1 || !ok2 {
		return true, nil // Invalid progress format, allow
	}

	return userProgressFloat >= requiredProgressFloat, nil
}

// TimeBasedAccessPolicy - Controls access based on time windows
type TimeBasedAccessPolicy struct{}

func (p *TimeBasedAccessPolicy) GetName() string {
	return "TimeBasedAccess"
}

func (p *TimeBasedAccessPolicy) Evaluate(ctx context.Context, request *auth.AccessRequest) (bool, error) {
	// Check for time-based restrictions
	startTime, hasStart := request.Resource.Attributes["access_start_time"]
	endTime, hasEnd := request.Resource.Attributes["access_end_time"]

	if !hasStart && !hasEnd {
		return true, nil // No time restrictions
	}

	now := time.Now()

	if hasStart {
		startTimeVal, ok := startTime.(time.Time)
		if ok && now.Before(startTimeVal) {
			return false, nil // Too early to access
		}
	}

	if hasEnd {
		endTimeVal, ok := endTime.(time.Time)
		if ok && now.After(endTimeVal) {
			return false, nil // Too late to access
		}
	}

	return true, nil
}

// OwnershipPolicy - Users can only modify resources they own
type OwnershipPolicy struct{}

func (p *OwnershipPolicy) GetName() string {
	return "Ownership"
}

func (p *OwnershipPolicy) Evaluate(ctx context.Context, request *auth.AccessRequest) (bool, error) {
	// This policy applies to update and delete actions
	modifyActions := []auth.Permission{
		auth.PermissionUpdateInternship,
		auth.PermissionDeleteInternship,
	}

	isModifyAction := false
	for _, action := range modifyActions {
		if request.Action == action {
			isModifyAction = true
			break
		}
	}

	if !isModifyAction {
		return true, nil // Policy doesn't apply to non-modify actions
	}

	// Admins can modify anything
	if request.Subject.Role == types.UserRoleAdmin {
		return true, nil
	}

	// Check ownership
	resourceOwner, exists := request.Resource.Attributes["created_by"]
	if !exists {
		return true, nil // No ownership info, allow (other policies will handle)
	}

	ownerID, ok := resourceOwner.(string)
	if !ok {
		return true, nil // Invalid owner format
	}

	// User can only modify resources they created
	return request.Subject.UserID == ownerID, nil
}
