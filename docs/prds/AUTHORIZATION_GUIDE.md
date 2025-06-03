# Authorization Guide: RBAC and ABAC Implementation

This guide explains how to implement and use Role-Based Access Control (RBAC) and Attribute-Based Access Control (ABAC) in your Go backend system.

## Overview

Our authorization system combines two approaches:

1. **RBAC (Role-Based Access Control)**: Controls access based on user roles (Student, Instructor, Admin)
2. **ABAC (Attribute-Based Access Control)**: Fine-grained access control based on attributes like ownership, enrollment, progress, etc.

## Architecture

```
Request → Authentication → Authorization (RBAC + ABAC) → Business Logic
```

### Components

1. **Authorization Service** (`internal/auth/authorization.go`)
2. **Authorization Middleware** (`internal/rest/middleware/authorization.go`)
3. **Domain Models** (`internal/domain/auth/model.go`)
4. **ABAC Policies** (Embedded in authorization service)

## User Roles and Permissions

### Role Hierarchy

```
Admin (Full Access)
├── All Instructor permissions
├── System configuration
├── User management
└── Analytics

Instructor (Content Creator)
├── Create/Update/Delete own internships
├── View all published internships
├── Access all content
└── View analytics for own content

Student (Consumer)
├── View enrolled internships
├── Access content based on enrollment and progress
└── Limited content access
```

### Permission Matrix

| Permission                 | Admin | Instructor | Student  |
| -------------------------- | ----- | ---------- | -------- |
| `internship:create`        | ✅    | ✅         | ❌       |
| `internship:update`        | ✅    | ✅\*       | ❌       |
| `internship:delete`        | ✅    | ✅\*       | ❌       |
| `internship:view`          | ✅    | ✅         | ✅\*\*   |
| `content:lectures:view`    | ✅    | ✅         | ✅\*\*\* |
| `content:assignments:view` | ✅    | ✅         | ✅\*\*\* |
| `content:resources:view`   | ✅    | ✅         | ✅\*\*\* |
| `users:manage`             | ✅    | ❌         | ❌       |
| `system:config`            | ✅    | ❌         | ❌       |

_\* Only own content_  
_\*\* Only enrolled internships_  
_\*\*\* Only enrolled internships with sufficient progress_

## RBAC Implementation

### 1. Role Definition

```go
// internal/types/user.go
type UserRole string

const (
    UserRoleStudent    UserRole = "STUDENT"
    UserRoleInstructor UserRole = "INSTRUCTOR"
    UserRoleAdmin      UserRole = "ADMIN"
)
```

### 2. Permission Definition

```go
// internal/domain/auth/model.go
type Permission string

const (
    PermissionCreateInternship Permission = "internship:create"
    PermissionViewInternship   Permission = "internship:view"
    // ... more permissions
)
```

### 3. Role-Permission Mapping

```go
// internal/auth/authorization.go
func initializeRolePermissions() map[types.UserRole][]auth.Permission {
    return map[types.UserRole][]auth.Permission{
        types.UserRoleAdmin: {
            // All permissions
        },
        types.UserRoleInstructor: {
            auth.PermissionCreateInternship,
            auth.PermissionUpdateInternship,
            // ... instructor permissions
        },
        types.UserRoleStudent: {
            auth.PermissionViewInternship,
        },
    }
}
```

## ABAC Implementation

### 1. Policy Interface

```go
type ABACPolicy interface {
    Evaluate(ctx context.Context, request *auth.AccessRequest) (bool, error)
    GetName() string
}
```

### 2. Built-in Policies

#### Enrollment Policy

Students can only access content of internships they're enrolled in.

```go
type EnrollmentBasedAccessPolicy struct{}

func (p *EnrollmentBasedAccessPolicy) Evaluate(ctx context.Context, request *auth.AccessRequest) (bool, error) {
    // Check if student is enrolled in the internship
    if request.Subject.Role != types.UserRoleStudent {
        return true, nil // Policy doesn't apply to non-students
    }

    // Get enrolled internships from user attributes
    enrolledInternships := request.Subject.Attributes["enrolled_internships"]
    // ... validation logic
}
```

#### Ownership Policy

Users can only modify resources they created.

```go
type OwnershipPolicy struct{}

func (p *OwnershipPolicy) Evaluate(ctx context.Context, request *auth.AccessRequest) (bool, error) {
    // Check if user owns the resource
    resourceOwner := request.Resource.Attributes["created_by"]
    return request.Subject.UserID == resourceOwner, nil
}
```

#### Progress Policy

Students need sufficient progress to access advanced content.

```go
type ContentAccessPolicy struct{}

func (p *ContentAccessPolicy) Evaluate(ctx context.Context, request *auth.AccessRequest) (bool, error) {
    // Check user progress against required progress
    userProgress := request.Subject.Attributes["progress"]
    requiredProgress := request.Resource.Attributes["required_progress"]
    // ... comparison logic
}
```

#### Time-based Policy

Content access based on time windows.

```go
type TimeBasedAccessPolicy struct{}

func (p *TimeBasedAccessPolicy) Evaluate(ctx context.Context, request *auth.AccessRequest) (bool, error) {
    // Check access time windows
    startTime := request.Resource.Attributes["access_start_time"]
    endTime := request.Resource.Attributes["access_end_time"]
    // ... time validation logic
}
```

## Usage Examples

### 1. Service Layer Authorization

```go
func (s *internshipService) Create(ctx context.Context, req *dto.CreateInternshipRequest, userID string, userRole types.UserRole) (*domainInternship.Internship, error) {
    // Create authorization request
    authRequest := &domainAuth.AccessRequest{
        Subject: &domainAuth.AuthContext{
            UserID: userID,
            Role:   userRole,
        },
        Resource: &domainAuth.Resource{
            Type: "internship",
        },
        Action: domainAuth.PermissionCreateInternship,
    }

    // Check authorization
    allowed, err := s.authzService.IsAuthorized(ctx, authRequest)
    if err != nil {
        return nil, err
    }

    if !allowed {
        return nil, ierr.ErrPermissionDenied
    }

    // Business logic continues...
}
```

### 2. Middleware-based Authorization

```go
// Role-based middleware
internshipGroup.POST("",
    middleware.RequireInstructorOrAdmin(),
    handler.CreateInternship,
)

// Permission-based middleware
internshipGroup.POST("",
    middleware.RequirePermission(domainAuth.PermissionCreateInternship, "internship"),
    handler.CreateInternship,
)
```

### 3. Content Access with ABAC

```go
func (h *InternshipHandler) AccessInternshipContent(c *gin.Context) {
    authContext := middleware.GetAuthContext(c)

    // Create authorization request with attributes
    authRequest := &domainAuth.AccessRequest{
        Subject: authContext,
        Resource: &domainAuth.Resource{
            Type: "content",
            ID:   contentID,
            Attributes: map[string]interface{}{
                "internship_id":     internshipID,
                "content_type":      contentType,
                "required_progress": 0.5, // 50% progress required
                "access_start_time": time.Now().Add(-24*time.Hour),
                "access_end_time":   time.Now().Add(24*time.Hour),
            },
        },
        Action: domainAuth.PermissionViewLectures,
    }

    allowed, err := h.authzService.IsAuthorized(c.Request.Context(), authRequest)
    // ... handle result
}
```

## Setting Up User Attributes

### 1. In Authentication Middleware

```go
func AuthorizationMiddleware(...) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get user from database
        user, err := userRepo.Get(c.Request.Context(), userID)

        // Create auth context with attributes
        authContext := &domainAuth.AuthContext{
            UserID: userID,
            Email:  userEmail,
            Role:   user.Role,
            Attributes: map[string]interface{}{
                "enrolled_internships": getUserEnrollments(userID),
                "progress": getUserProgress(userID),
                "subscription_level": user.SubscriptionLevel,
            },
        }

        c.Set("auth_context", authContext)
        c.Next()
    }
}
```

### 2. Dynamic Attribute Loading

```go
func (s *authorizationService) IsAuthorized(ctx context.Context, request *auth.AccessRequest) (bool, error) {
    // Load additional attributes if needed
    if request.Action == auth.PermissionViewLectures {
        enrollments, err := s.enrollmentService.GetUserEnrollments(ctx, request.Subject.UserID)
        if err == nil {
            request.Subject.Attributes["enrolled_internships"] = enrollments
        }
    }

    // Continue with authorization...
}
```

## Best Practices

### 1. Security

- Always validate user input before authorization checks
- Use least privilege principle
- Log authorization decisions for audit trails
- Implement rate limiting for authorization checks

### 2. Performance

- Cache role-permission mappings
- Minimize database calls in ABAC policies
- Use efficient data structures for attribute lookups
- Consider async attribute loading for non-critical attributes

### 3. Maintainability

- Keep policies simple and focused
- Document policy logic clearly
- Use descriptive permission names
- Separate authorization logic from business logic

### 4. Testing

```go
func TestInternshipCreation(t *testing.T) {
    // Test RBAC
    t.Run("instructor can create internship", func(t *testing.T) {
        authRequest := &auth.AccessRequest{
            Subject: &auth.AuthContext{Role: types.UserRoleInstructor},
            Resource: &auth.Resource{Type: "internship"},
            Action: auth.PermissionCreateInternship,
        }

        allowed, err := authzService.IsAuthorized(ctx, authRequest)
        assert.NoError(t, err)
        assert.True(t, allowed)
    })

    t.Run("student cannot create internship", func(t *testing.T) {
        authRequest := &auth.AccessRequest{
            Subject: &auth.AuthContext{Role: types.UserRoleStudent},
            Resource: &auth.Resource{Type: "internship"},
            Action: auth.PermissionCreateInternship,
        }

        allowed, err := authzService.IsAuthorized(ctx, authRequest)
        assert.NoError(t, err)
        assert.False(t, allowed)
    })
}

func TestContentAccess(t *testing.T) {
    // Test ABAC
    t.Run("enrolled student can access content", func(t *testing.T) {
        authRequest := &auth.AccessRequest{
            Subject: &auth.AuthContext{
                Role: types.UserRoleStudent,
                Attributes: map[string]interface{}{
                    "enrolled_internships": []string{"internship-123"},
                },
            },
            Resource: &auth.Resource{
                Type: "content",
                Attributes: map[string]interface{}{
                    "internship_id": "internship-123",
                },
            },
            Action: auth.PermissionViewLectures,
        }

        allowed, err := authzService.IsAuthorized(ctx, authRequest)
        assert.NoError(t, err)
        assert.True(t, allowed)
    })
}
```

## Common Use Cases

### 1. Creating an Internship

- **RBAC**: Only instructors and admins can create
- **ABAC**: No additional restrictions

### 2. Viewing an Internship

- **RBAC**: All authenticated users can view
- **ABAC**: Students can only view enrolled internships

### 3. Updating an Internship

- **RBAC**: Instructors and admins can update
- **ABAC**: Instructors can only update their own internships

### 4. Accessing Content

- **RBAC**: Students can access content
- **ABAC**: Must be enrolled + sufficient progress + within time window

### 5. Admin Operations

- **RBAC**: Only admins
- **ABAC**: No additional restrictions

## Extending the System

### Adding New Roles

1. Add role constant to `internal/types/user.go`
2. Update role-permission mapping in `initializeRolePermissions()`
3. Update middleware functions if needed

### Adding New Permissions

1. Add permission constant to `internal/domain/auth/model.go`
2. Update role-permission mappings
3. Use in service layer or middleware

### Adding New ABAC Policies

1. Implement `ABACPolicy` interface
2. Add to `initializeABACPolicies()`
3. Test thoroughly

### Custom Attribute Providers

```go
type AttributeProvider interface {
    LoadAttributes(ctx context.Context, userID string) (map[string]interface{}, error)
}

type EnrollmentAttributeProvider struct {
    enrollmentRepo EnrollmentRepository
}

func (p *EnrollmentAttributeProvider) LoadAttributes(ctx context.Context, userID string) (map[string]interface{}, error) {
    enrollments, err := p.enrollmentRepo.GetByUserID(ctx, userID)
    if err != nil {
        return nil, err
    }

    return map[string]interface{}{
        "enrolled_internships": extractInternshipIDs(enrollments),
        "enrollment_dates": extractEnrollmentDates(enrollments),
    }, nil
}
```

This authorization system provides a flexible and scalable approach to access control, combining the simplicity of RBAC with the fine-grained control of ABAC.
