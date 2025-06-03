# Authorization Integration Guide

## Overview

This guide explains how **RBAC** and **ABAC** work together in our authorization system to provide **comprehensive, scalable access control** that balances **performance** with **flexibility**.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                    Authorization Request                        │
└─────────────────────────┬───────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────┐
│                  Unified Authorization Service                  │
└─────────────────────────┬───────────────────────────────────────┘
                          │
           ┌──────────────┴──────────────┐
           │                             │
           ▼                             ▼
┌──────────────────┐            ┌──────────────────┐
│  RBAC Service    │            │  ABAC Service    │
│                  │            │                  │
│ • Fast O(1)      │            │ • Policy Engine  │
│ • Role-based     │            │ • Attribute-     │
│ • Static rules   │            │   based          │
│                  │            │ • Dynamic rules  │
└─────────┬────────┘            └─────────┬────────┘
          │                               │
          ▼                               ▼
┌──────────────────┐            ┌──────────────────┐
│ Permission Check │            │ Policy           │
│ Result           │            │ Evaluation       │
│                  │            │ Result           │
└─────────┬────────┘            └─────────┬────────┘
          │                               │
          └──────────────┬──────────────────┘
                         │
                         ▼
                ┌──────────────────┐
                │ Final Decision   │
                │ (Allow/Deny)     │
                └──────────────────┘
```

## Authorization Flow

### 1. **Two-Phase Authorization**

```go
func (s *unifiedService) IsAuthorized(ctx context.Context, request *auth.AccessRequest) (bool, error) {
    // Phase 1: RBAC Check (Fast)
    if !s.rbacService.HasPermission(request.Subject.Role, request.Action) {
        return false, nil // RBAC denies - immediate fail
    }

    // Phase 2: ABAC Check (Contextual)
    return s.abacService.Evaluate(ctx, request)
}
```

**Benefits:**

- **Performance**: RBAC fails fast for unauthorized roles
- **Security**: ABAC provides fine-grained control
- **Maintainability**: Clear separation of concerns

### 2. **Decision Matrix**

| RBAC Result | ABAC Result | Final Decision | Use Case Example                      |
| ----------- | ----------- | -------------- | ------------------------------------- |
| ❌ Deny     | N/A         | ❌ Deny        | Student trying to delete internship   |
| ✅ Allow    | ❌ Deny     | ❌ Deny        | Student accessing non-enrolled course |
| ✅ Allow    | ✅ Allow    | ✅ Allow       | Student accessing enrolled course     |
| ✅ Allow    | ⚠️ Error    | ❌ Deny        | ABAC evaluation failure               |

### 3. **Performance Characteristics**

| Component | Complexity | Cache     | Use Case               |
| --------- | ---------- | --------- | ---------------------- |
| **RBAC**  | O(1)       | In-memory | Role-based permissions |
| **ABAC**  | O(n)       | 5min TTL  | Context-aware policies |

## Integration Patterns

### 1. **Basic Integration**

```go
// Simple authorization check
authService := auth.NewUnifiedAuthorizationService(logger)

request := &auth.AccessRequest{
    Subject: &auth.AuthContext{
        UserID: "user_123",
        Role:   types.UserRoleStudent,
    },
    Resource: &auth.Resource{
        Type: "internship",
        ID:   "internship_456",
    },
    Action: auth.PermissionViewInternship,
}

allowed, err := authService.IsAuthorized(ctx, request)
```

### 2. **Middleware Integration**

```go
// HTTP middleware
func RequirePermission(permission auth.Permission, resourceType string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get auth context from previous middleware
        authContext := getAuthContext(c)

        // Build access request
        request := &auth.AccessRequest{
            Subject: authContext,
            Resource: &auth.Resource{
                Type: resourceType,
                Attributes: extractResourceAttributes(c),
            },
            Action: permission,
        }

        // Check authorization
        allowed, err := authService.IsAuthorized(c.Request.Context(), request)
        if err != nil || !allowed {
            c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
            c.Abort()
            return
        }

        c.Next()
    }
}
```

### 3. **Service Layer Integration**

```go
// Service method with authorization
func (s *InternshipService) UpdateInternship(ctx context.Context, id string, updates *UpdateRequest, authContext *auth.AuthContext) error {
    // Build authorization request
    request := &auth.AccessRequest{
        Subject: authContext,
        Resource: &auth.Resource{
            Type: "internship",
            ID:   id,
            Attributes: map[string]interface{}{
                "created_by": existingInternship.CreatedBy, // For ownership check
            },
        },
        Action: auth.PermissionUpdateInternship,
    }

    // Check authorization
    allowed, err := s.authService.IsAuthorized(ctx, request)
    if err != nil {
        return fmt.Errorf("authorization check failed: %w", err)
    }
    if !allowed {
        return ierr.ErrPermissionDenied
    }

    // Proceed with update
    return s.repository.Update(ctx, id, updates)
}
```

### 4. **Fluent Builder Pattern**

```go
// Complex authorization with builder
allowed, err := authService.NewAuthorizationBuilder().
    ForUser("user_123", types.UserRoleStudent).
    OnResource("content", "lecture_456").
    WithAction(auth.PermissionViewLectures).
    WithUserAttribute("enrolled_internships", []string{"internship_789"}).
    WithResourceAttribute("internship_id", "internship_789").
    WithResourceAttribute("required_progress", 50.0).
    Check(ctx)
```

## Real-World Scenarios

### 1. **Student Accessing Course Content**

```go
// Scenario: Student wants to view a lecture
request := &auth.AccessRequest{
    Subject: &auth.AuthContext{
        UserID: "student_123",
        Role:   types.UserRoleStudent,
        Attributes: map[string]interface{}{
            "enrolled_internships": []string{"internship_456", "internship_789"},
            "progress": 75.0,
        },
    },
    Resource: &auth.Resource{
        Type: "content",
        ID:   "lecture_001",
        Attributes: map[string]interface{}{
            "internship_id": "internship_456",
            "required_progress": 50.0,
            "access_start_time": time.Now().Add(-1 * time.Hour),
            "access_end_time": time.Now().Add(24 * time.Hour),
        },
    },
    Action: auth.PermissionViewLectures,
}

// Authorization flow:
// 1. RBAC: Student role has PermissionViewLectures ✅
// 2. ABAC Policies:
//    - EnrollmentPolicy: Student enrolled in internship_456 ✅
//    - ProgressPolicy: 75% >= 50% required ✅
//    - TimePolicy: Current time within access window ✅
// Result: ALLOW
```

### 2. **Instructor Modifying Internship**

```go
// Scenario: Instructor wants to update an internship
request := &auth.AccessRequest{
    Subject: &auth.AuthContext{
        UserID: "instructor_456",
        Role:   types.UserRoleInstructor,
    },
    Resource: &auth.Resource{
        Type: "internship",
        ID:   "internship_789",
        Attributes: map[string]interface{}{
            "created_by": "instructor_456", // Same instructor
        },
    },
    Action: auth.PermissionUpdateInternship,
}

// Authorization flow:
// 1. RBAC: Instructor role has PermissionUpdateInternship ✅
// 2. ABAC Policies:
//    - OwnershipPolicy: instructor_456 == instructor_456 ✅
// Result: ALLOW
```

### 3. **Admin System Configuration**

```go
// Scenario: Admin wants to modify system settings
request := &auth.AccessRequest{
    Subject: &auth.AuthContext{
        UserID: "admin_789",
        Role:   types.UserRoleAdmin,
    },
    Resource: &auth.Resource{
        Type: "system",
        ID:   "config",
    },
    Action: auth.PermissionSystemConfig,
}

// Authorization flow:
// 1. RBAC: Admin role has PermissionSystemConfig ✅
// 2. ABAC Policies: No applicable policies ✅ (default allow)
// Result: ALLOW
```

## Custom Extensions

### 1. **Adding New RBAC Permissions**

```go
// 1. Add permission constant
const PermissionCreateCertificate Permission = "certificate:create"

// 2. Update role mappings in rbac/service.go
types.UserRoleInstructor: {
    // ... existing permissions
    auth.PermissionCreateCertificate,
}

// 3. Use in middleware
middleware.RequirePermission(auth.PermissionCreateCertificate, "certificate")
```

### 2. **Creating Custom ABAC Policy**

```go
// 1. Implement Policy interface
type CompliancePolicy struct{}

func (p *CompliancePolicy) Evaluate(ctx context.Context, request *auth.AccessRequest) (abac.Decision, error) {
    // Check compliance requirements
    complianceLevel := request.Subject.Attributes["compliance_level"]
    requiredLevel := request.Resource.Attributes["required_compliance"]

    if complianceLevel.(int) >= requiredLevel.(int) {
        return abac.Decision{Allow: true, Reason: "Compliance requirements met"}, nil
    }

    return abac.Decision{Allow: false, Reason: "Insufficient compliance level"}, nil
}

// 2. Register policy
authService.RegisterABACPolicy(&CompliancePolicy{})
```

### 3. **Custom Attribute Provider**

```go
// 1. Implement AttributeProvider interface
type DatabaseAttributeProvider struct {
    db *sql.DB
}

func (p *DatabaseAttributeProvider) LoadUserAttributes(ctx context.Context, userID string) (map[string]interface{}, error) {
    // Load from database
    var attrs map[string]interface{}
    // ... database query
    return attrs, nil
}

// 2. Register provider
authService.RegisterAttributeProvider(&DatabaseAttributeProvider{db: db})
```

## Testing Strategies

### 1. **Unit Testing**

```go
func TestAuthorizationFlow(t *testing.T) {
    // Test RBAC only
    rbacService := rbac.NewService(logger)
    assert.True(t, rbacService.HasPermission(types.UserRoleAdmin, auth.PermissionSystemConfig))

    // Test ABAC only
    abacService := abac.NewService(logger)
    decision, err := abacService.Evaluate(ctx, request)
    assert.NoError(t, err)
    assert.True(t, decision)

    // Test integration
    unifiedService := auth.NewUnifiedAuthorizationService(logger)
    allowed, err := unifiedService.IsAuthorized(ctx, request)
    assert.NoError(t, err)
    assert.True(t, allowed)
}
```

### 2. **Integration Testing**

```go
func TestMiddlewareIntegration(t *testing.T) {
    // Setup test server with middleware
    router := gin.New()
    router.Use(middleware.AuthorizationMiddleware(cfg, logger, authService, userRepo))
    router.GET("/internships/:id",
        middleware.RequirePermission(auth.PermissionViewInternship, "internship"),
        handler.GetInternship,
    )

    // Test authorized request
    req := httptest.NewRequest("GET", "/internships/123", nil)
    req.Header.Set("Authorization", "Bearer " + validToken)

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
}
```

### 3. **Performance Testing**

```go
func BenchmarkAuthorization(b *testing.B) {
    authService := auth.NewUnifiedAuthorizationService(logger)
    request := createTestRequest()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := authService.IsAuthorized(context.Background(), request)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

## Monitoring and Observability

### 1. **Metrics**

```go
// Track authorization decisions
authorizationDecisions.WithLabelValues(
    string(request.Subject.Role),
    string(request.Action),
    request.Resource.Type,
    strconv.FormatBool(allowed),
).Inc()

// Track authorization latency
authorizationDuration.WithLabelValues(
    "rbac",
).Observe(rbacDuration.Seconds())

authorizationDuration.WithLabelValues(
    "abac",
).Observe(abacDuration.Seconds())
```

### 2. **Audit Logging**

```go
logger.Infow("Authorization decision",
    "user_id", request.Subject.UserID,
    "user_role", request.Subject.Role,
    "action", request.Action,
    "resource_type", request.Resource.Type,
    "resource_id", request.Resource.ID,
    "decision", allowed,
    "rbac_allowed", rbacAllowed,
    "abac_allowed", abacAllowed,
    "duration_ms", duration.Milliseconds(),
)
```

### 3. **Error Handling**

```go
// Graceful degradation
func (s *unifiedService) IsAuthorized(ctx context.Context, request *auth.AccessRequest) (bool, error) {
    // Always check RBAC first
    rbacAllowed := s.rbacService.HasPermission(request.Subject.Role, request.Action)
    if !rbacAllowed {
        return false, nil
    }

    // Try ABAC, but don't fail if there are issues
    abacAllowed, err := s.abacService.Evaluate(ctx, request)
    if err != nil {
        // Log error but don't fail the request
        s.logger.Errorw("ABAC evaluation failed, falling back to RBAC", "error", err)
        return rbacAllowed, nil // Fallback to RBAC decision
    }

    return abacAllowed, nil
}
```

## Migration Guide

### From Existing System

1. **Phase 1**: Add RBAC layer alongside existing authorization
2. **Phase 2**: Implement ABAC policies gradually
3. **Phase 3**: Migrate middleware to use unified service
4. **Phase 4**: Remove old authorization code

### Example Migration

```go
// Before: Direct role check
if user.Role != types.UserRoleAdmin {
    return ierr.ErrPermissionDenied
}

// After: Unified authorization
allowed, err := authService.IsAuthorized(ctx, &auth.AccessRequest{
    Subject: userAuthContext,
    Resource: &auth.Resource{Type: "system", ID: "config"},
    Action: auth.PermissionSystemConfig,
})
if err != nil || !allowed {
    return ierr.ErrPermissionDenied
}
```

This integration provides a **robust, scalable authorization system** that balances **performance** with **flexibility**, making it easy to implement complex access control scenarios while maintaining clean, maintainable code.
