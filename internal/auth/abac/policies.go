package abac

import (
	"context"
	"fmt"
	"time"

	"github.com/omkar273/codegeeky/internal/domain/auth"
	"github.com/omkar273/codegeeky/internal/types"
)

// EnrollmentBasedAccessPolicy - Students can only access content of internships they're enrolled in
type EnrollmentBasedAccessPolicy struct{}

// NewEnrollmentBasedAccessPolicy creates a new enrollment-based access policy
func NewEnrollmentBasedAccessPolicy() Policy {
	return &EnrollmentBasedAccessPolicy{}
}

// GetName returns the policy name
func (p *EnrollmentBasedAccessPolicy) GetName() string {
	return "EnrollmentBasedAccess"
}

// GetPriority returns the policy priority
func (p *EnrollmentBasedAccessPolicy) GetPriority() int {
	return 100 // High priority for enrollment checks
}

// Applies checks if this policy is applicable to the request
func (p *EnrollmentBasedAccessPolicy) Applies(request *auth.AccessRequest) bool {
	// This policy applies to students accessing content
	if request.Subject.Role != types.UserRoleStudent {
		return false
	}

	// For content access permissions
	contentPermissions := []auth.Permission{
		auth.PermissionViewLectures,
		auth.PermissionViewAssignments,
		auth.PermissionViewResources,
		auth.PermissionDownloadContent,
		auth.PermissionViewInternship,
	}

	for _, perm := range contentPermissions {
		if request.Action == perm {
			return true
		}
	}

	return false
}

// Evaluate evaluates the enrollment-based access policy
func (p *EnrollmentBasedAccessPolicy) Evaluate(ctx context.Context, request *auth.AccessRequest) (Decision, error) {
	// Check if student is enrolled in the internship
	internshipID, exists := request.Resource.Attributes["internship_id"]
	if !exists {
		return Decision{
			Allow:  false,
			Reason: "internship_id not found in resource attributes",
		}, nil
	}

	// Get user's enrolled internships from context or attributes
	enrolledInternships, exists := request.Subject.Attributes["enrolled_internships"]
	if !exists {
		return Decision{
			Allow:  false,
			Reason: "Student not enrolled in any internships",
		}, nil
	}

	enrolledList, ok := enrolledInternships.([]string)
	if !ok {
		return Decision{
			Allow:  false,
			Reason: "Invalid enrolled_internships format",
		}, nil
	}

	internshipIDStr, ok := internshipID.(string)
	if !ok {
		return Decision{
			Allow:  false,
			Reason: "Invalid internship_id format",
		}, nil
	}

	// Check if student is enrolled in this internship
	for _, enrolledID := range enrolledList {
		if enrolledID == internshipIDStr {
			return Decision{
				Allow:  true,
				Reason: "Student enrolled in internship",
			}, nil
		}
	}

	return Decision{
		Allow:  false,
		Reason: "Student not enrolled in this internship",
	}, nil
}

// OwnershipPolicy - Users can only modify resources they own
type OwnershipPolicy struct{}

// NewOwnershipPolicy creates a new ownership policy
func NewOwnershipPolicy() Policy {
	return &OwnershipPolicy{}
}

// GetName returns the policy name
func (p *OwnershipPolicy) GetName() string {
	return "Ownership"
}

// GetPriority returns the policy priority
func (p *OwnershipPolicy) GetPriority() int {
	return 200 // Very high priority for ownership checks
}

// Applies checks if this policy is applicable to the request
func (p *OwnershipPolicy) Applies(request *auth.AccessRequest) bool {
	// This policy applies to update and delete actions
	modifyActions := []auth.Permission{
		auth.PermissionUpdateInternship,
		auth.PermissionDeleteInternship,
	}

	for _, action := range modifyActions {
		if request.Action == action {
			return true
		}
	}

	return false
}

// Evaluate evaluates the ownership policy
func (p *OwnershipPolicy) Evaluate(ctx context.Context, request *auth.AccessRequest) (Decision, error) {
	// Admins can modify anything
	if request.Subject.Role == types.UserRoleAdmin {
		return Decision{
			Allow:  true,
			Reason: "Admin can modify any resource",
		}, nil
	}

	// Check ownership
	resourceOwner, exists := request.Resource.Attributes["created_by"]
	if !exists {
		return Decision{
			Allow:  true,
			Reason: "No ownership info available",
		}, nil
	}

	ownerID, ok := resourceOwner.(string)
	if !ok {
		return Decision{
			Allow:  true,
			Reason: "Invalid owner format",
		}, nil
	}

	// User can only modify resources they created
	if request.Subject.UserID == ownerID {
		return Decision{
			Allow:  true,
			Reason: "User owns the resource",
		}, nil
	}

	return Decision{
		Allow:  false,
		Reason: "User does not own the resource",
	}, nil
}

// TimeBasedAccessPolicy - Controls access based on time windows
type TimeBasedAccessPolicy struct{}

// NewTimeBasedAccessPolicy creates a new time-based access policy
func NewTimeBasedAccessPolicy() Policy {
	return &TimeBasedAccessPolicy{}
}

// GetName returns the policy name
func (p *TimeBasedAccessPolicy) GetName() string {
	return "TimeBasedAccess"
}

// GetPriority returns the policy priority
func (p *TimeBasedAccessPolicy) GetPriority() int {
	return 150 // High priority for time restrictions
}

// Applies checks if this policy is applicable to the request
func (p *TimeBasedAccessPolicy) Applies(request *auth.AccessRequest) bool {
	// Check if resource has time-based restrictions
	_, hasStart := request.Resource.Attributes["access_start_time"]
	_, hasEnd := request.Resource.Attributes["access_end_time"]

	return hasStart || hasEnd
}

// Evaluate evaluates the time-based access policy
func (p *TimeBasedAccessPolicy) Evaluate(ctx context.Context, request *auth.AccessRequest) (Decision, error) {
	now := time.Now()

	// Check start time restriction
	startTime, hasStart := request.Resource.Attributes["access_start_time"]
	if hasStart {
		startTimeVal, ok := startTime.(time.Time)
		if ok && now.Before(startTimeVal) {
			return Decision{
				Allow:  false,
				Reason: "Access not yet available",
				Metadata: map[string]interface{}{
					"available_at": startTimeVal,
				},
			}, nil
		}
	}

	// Check end time restriction
	endTime, hasEnd := request.Resource.Attributes["access_end_time"]
	if hasEnd {
		endTimeVal, ok := endTime.(time.Time)
		if ok && now.After(endTimeVal) {
			return Decision{
				Allow:  false,
				Reason: "Access period has expired",
				Metadata: map[string]interface{}{
					"expired_at": endTimeVal,
				},
			}, nil
		}
	}

	return Decision{
		Allow:  true,
		Reason: "Within allowed time window",
	}, nil
}

// ProgressBasedPolicy - Controls access to different types of content based on progress
type ProgressBasedPolicy struct{}

// NewProgressBasedPolicy creates a new progress-based policy
func NewProgressBasedPolicy() Policy {
	return &ProgressBasedPolicy{}
}

// GetName returns the policy name
func (p *ProgressBasedPolicy) GetName() string {
	return "ProgressBased"
}

// GetPriority returns the policy priority
func (p *ProgressBasedPolicy) GetPriority() int {
	return 75 // Medium priority for progress checks
}

// Applies checks if this policy is applicable to the request
func (p *ProgressBasedPolicy) Applies(request *auth.AccessRequest) bool {
	// This policy applies to content access for students
	if request.Subject.Role != types.UserRoleStudent {
		return false
	}

	// Check if resource has progress requirements
	_, hasRequirement := request.Resource.Attributes["required_progress"]
	return hasRequirement
}

// Evaluate evaluates the progress-based policy
func (p *ProgressBasedPolicy) Evaluate(ctx context.Context, request *auth.AccessRequest) (Decision, error) {
	// For students, check if they have the prerequisite progress
	requiredProgress, exists := request.Resource.Attributes["required_progress"]
	if !exists {
		return Decision{
			Allow:  true,
			Reason: "No progress requirement",
		}, nil
	}

	userProgress, exists := request.Subject.Attributes["progress"]
	if !exists {
		return Decision{
			Allow:  false,
			Reason: "User has no progress recorded",
		}, nil
	}

	// Compare progress
	userProgressFloat, ok1 := userProgress.(float64)
	requiredProgressFloat, ok2 := requiredProgress.(float64)

	if !ok1 || !ok2 {
		return Decision{
			Allow:  true,
			Reason: "Invalid progress format, allowing access",
		}, nil
	}

	if userProgressFloat >= requiredProgressFloat {
		return Decision{
			Allow:  true,
			Reason: fmt.Sprintf("Sufficient progress: %.1f%% >= %.1f%%", userProgressFloat, requiredProgressFloat),
			Metadata: map[string]interface{}{
				"user_progress":     userProgressFloat,
				"required_progress": requiredProgressFloat,
			},
		}, nil
	}

	return Decision{
		Allow:  false,
		Reason: fmt.Sprintf("Insufficient progress: %.1f%% < %.1f%%", userProgressFloat, requiredProgressFloat),
		Metadata: map[string]interface{}{
			"user_progress":     userProgressFloat,
			"required_progress": requiredProgressFloat,
		},
	}, nil
}

// ContentAccessPolicy - Controls access to different types of content
type ContentAccessPolicy struct{}

// NewContentAccessPolicy creates a new content access policy
func NewContentAccessPolicy() Policy {
	return &ContentAccessPolicy{}
}

// GetName returns the policy name
func (p *ContentAccessPolicy) GetName() string {
	return "ContentAccess"
}

// GetPriority returns the policy priority
func (p *ContentAccessPolicy) GetPriority() int {
	return 50 // Lower priority, general content access rules
}

// Applies checks if this policy is applicable to the request
func (p *ContentAccessPolicy) Applies(request *auth.AccessRequest) bool {
	// This policy applies to content access
	return request.Resource.Type == "content"
}

// Evaluate evaluates the content access policy
func (p *ContentAccessPolicy) Evaluate(ctx context.Context, request *auth.AccessRequest) (Decision, error) {
	// Admins and instructors have full access
	if request.Subject.Role == types.UserRoleAdmin || request.Subject.Role == types.UserRoleInstructor {
		return Decision{
			Allow:  true,
			Reason: "Admin/Instructor has full content access",
		}, nil
	}

	// For students, this is a general content access policy
	// Other more specific policies (enrollment, progress) will handle detailed checks
	return Decision{
		Allow:  true,
		Reason: "General content access allowed (subject to other policies)",
	}, nil
}
