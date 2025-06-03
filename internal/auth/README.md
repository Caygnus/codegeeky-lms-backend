# Authorization System

## Overview

A **comprehensive, modular authorization system** that combines **RBAC** (Role-Based Access Control) and **ABAC** (Attribute-Based Access Control) to provide both **fast role-based permissions** and **fine-grained contextual access control**.

## 🏗️ Architecture

```
📁 internal/auth/
├── 📁 rbac/                    # Role-Based Access Control
│   ├── 📄 README.md           # RBAC documentation
│   └── 📄 service.go          # RBAC implementation
├── 📁 abac/                    # Attribute-Based Access Control
│   ├── 📄 README.md           # ABAC documentation
│   ├── 📄 service.go          # ABAC engine
│   ├── 📄 policies.go         # Built-in policies
│   └── 📄 combiners.go        # Policy combination strategies
├── 📁 integration/             # Integration guides
│   └── 📄 README.md           # How RBAC + ABAC work together
├── 📄 README.md               # This file
├── 📄 service.go              # Unified authorization service
├── 📄 authorization.go        # Legacy implementation (to be replaced)
└── 📄 provider.go             # Service provider
```

## 🚀 Quick Start

### 1. **Basic Setup**

```go
import (
    "github.com/omkar273/codegeeky/internal/auth"
    "github.com/omkar273/codegeeky/internal/domain/auth"
)

// Create unified authorization service
authService := auth.NewUnifiedAuthorizationService(logger)

// Basic authorization check
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
if err != nil {
    // Handle error
}
if !allowed {
    // Access denied
}
```

### 2. **Middleware Integration**

```go
// In your router setup
router.GET("/internships/:id",
    middleware.AuthorizationMiddleware(cfg, logger, authService, userRepo),
    middleware.RequirePermission(auth.PermissionViewInternship, "internship"),
    handler.GetInternship,
)
```

### 3. **Service Layer Integration**

```go
func (s *InternshipService) UpdateInternship(ctx context.Context, id string, updates *UpdateRequest, authContext *auth.AuthContext) error {
    // Check authorization before proceeding
    allowed, err := s.authService.IsAuthorized(ctx, &auth.AccessRequest{
        Subject: authContext,
        Resource: &auth.Resource{
            Type: "internship",
            ID:   id,
            Attributes: map[string]interface{}{
                "created_by": existingInternship.CreatedBy,
            },
        },
        Action: auth.PermissionUpdateInternship,
    })

    if err != nil || !allowed {
        return ierr.ErrPermissionDenied
    }

    // Proceed with business logic
    return s.repository.Update(ctx, id, updates)
}
```

## 🔐 Permission System

### **Role Hierarchy**

```
👑 Admin
├── Full system access
├── All permissions
└── No restrictions

👨‍🏫 Instructor
├── Content management
├── Own internship operations
└── Limited analytics

👨‍🎓 Student
├── View enrolled content
├── Read-only access
└── Progress-based restrictions
```

### **Permission Categories**

| Category        | Permissions                           | Admin | Instructor | Student |
| --------------- | ------------------------------------- | ----- | ---------- | ------- |
| **Internships** | Create, Update, Delete, View, Publish | ✅    | ✅\*       | ✅\*\*  |
| **Content**     | View Lectures, Assignments, Resources | ✅    | ✅         | ✅\*\*  |
| **Users**       | Manage, View All                      | ✅    | ❌         | ❌      |
| **System**      | Config, Analytics                     | ✅    | ✅\*\*\*   | ❌      |

_\* Limited by ownership (ABAC)_  
_\*\* Limited by enrollment (ABAC)_  
_\*\*\* Own analytics only (ABAC)_

## 🎯 Key Features

### **1. Two-Layer Security**

- **RBAC Layer**: Fast O(1) role-based permission checks
- **ABAC Layer**: Context-aware policy evaluation
- **Combined**: Both layers must approve for access

### **2. Built-in Policies**

- **🎓 Enrollment Policy**: Students access only enrolled content
- **👤 Ownership Policy**: Users modify only their resources
- **⏰ Time-based Policy**: Content available within time windows
- **📊 Progress Policy**: Access gated by completion requirements

### **3. Flexible Architecture**

- **🔌 Plugin-based**: Easy to add custom policies
- **⚡ High Performance**: Cached decisions and fast lookups
- **📈 Scalable**: Handles complex authorization scenarios
- **🧪 Testable**: Comprehensive testing support

## 📊 Performance

| Operation       | Complexity    | Cache     | Typical Latency |
| --------------- | ------------- | --------- | --------------- |
| RBAC Check      | O(1)          | In-memory | < 1ms           |
| ABAC Evaluation | O(n) policies | 5min TTL  | 5-20ms          |
| Combined Check  | O(1) + O(n)   | Cached    | 5-25ms          |

## 📖 Documentation

### **Core Components**

- **[RBAC Documentation](rbac/README.md)** - Role-based permissions
- **[ABAC Documentation](abac/README.md)** - Attribute-based policies
- **[Integration Guide](integration/README.md)** - How components work together

### **Examples & Use Cases**

#### **Student Accessing Course Content**

```go
// Student can view lecture if:
// 1. RBAC: Student role has view permission ✅
// 2. ABAC: Student enrolled in course ✅
// 3. ABAC: Sufficient progress completed ✅
// 4. ABAC: Within time window ✅

authService.NewAuthorizationBuilder().
    ForUser("student_123", types.UserRoleStudent).
    OnResource("content", "lecture_456").
    WithAction(auth.PermissionViewLectures).
    WithUserAttribute("enrolled_internships", []string{"internship_789"}).
    WithResourceAttribute("internship_id", "internship_789").
    WithResourceAttribute("required_progress", 50.0).
    Check(ctx)
```

#### **Instructor Managing Content**

```go
// Instructor can update internship if:
// 1. RBAC: Instructor role has update permission ✅
// 2. ABAC: Instructor owns the internship ✅

authService.IsAuthorized(ctx, &auth.AccessRequest{
    Subject: instructorAuthContext,
    Resource: &auth.Resource{
        Type: "internship",
        ID:   "internship_456",
        Attributes: map[string]interface{}{
            "created_by": "instructor_123", // Same as subject.UserID
        },
    },
    Action: auth.PermissionUpdateInternship,
})
```

## 🔧 Custom Extensions

### **Adding New Permissions**

```go
// 1. Define permission constant
const PermissionCreateCertificate auth.Permission = "certificate:create"

// 2. Update role mappings in rbac/service.go
types.UserRoleInstructor: {
    // ... existing permissions
    auth.PermissionCreateCertificate,
}

// 3. Use in middleware
middleware.RequirePermission(auth.PermissionCreateCertificate, "certificate")
```

### **Custom ABAC Policy**

```go
// 1. Implement Policy interface
type GeolocationPolicy struct{}

func (p *GeolocationPolicy) Evaluate(ctx context.Context, request *auth.AccessRequest) (abac.Decision, error) {
    userCountry := request.Subject.Attributes["country"]
    allowedCountries := request.Resource.Attributes["allowed_countries"]

    // Implementation logic...
    return abac.Decision{Allow: true, Reason: "Geolocation check passed"}, nil
}

// 2. Register policy
authService.RegisterABACPolicy(&GeolocationPolicy{})
```

### **Custom Attribute Provider**

```go
// 1. Implement AttributeProvider interface
type EnrollmentAttributeProvider struct {
    enrollmentRepo EnrollmentRepository
}

func (p *EnrollmentAttributeProvider) LoadUserAttributes(ctx context.Context, userID string) (map[string]interface{}, error) {
    enrollments, err := p.enrollmentRepo.GetActiveByUserID(ctx, userID)
    if err != nil {
        return nil, err
    }

    return map[string]interface{}{
        "enrolled_internships": extractInternshipIDs(enrollments),
        "progress": calculateProgress(enrollments),
    }, nil
}

// 2. Register provider
authService.RegisterAttributeProvider(&EnrollmentAttributeProvider{})
```

## 🧪 Testing

### **Unit Testing**

```go
func TestRBACPermissions(t *testing.T) {
    rbacService := rbac.NewService(logger)

    // Test admin permissions
    assert.True(t, rbacService.HasPermission(
        types.UserRoleAdmin,
        auth.PermissionSystemConfig,
    ))

    // Test student restrictions
    assert.False(t, rbacService.HasPermission(
        types.UserRoleStudent,
        auth.PermissionDeleteInternship,
    ))
}

func TestABACPolicies(t *testing.T) {
    abacService := abac.NewService(logger)

    request := createEnrolledStudentRequest()
    allowed, err := abacService.Evaluate(ctx, request)

    assert.NoError(t, err)
    assert.True(t, allowed)
}
```

### **Integration Testing**

```go
func TestUnifiedAuthorization(t *testing.T) {
    authService := auth.NewUnifiedAuthorizationService(logger)

    // Test complex scenario
    request := &auth.AccessRequest{
        Subject: createStudentContext(),
        Resource: createContentResource(),
        Action: auth.PermissionViewLectures,
    }

    allowed, err := authService.IsAuthorized(ctx, request)
    assert.NoError(t, err)
    assert.True(t, allowed)
}
```

## 📈 Monitoring

### **Metrics**

```go
// Authorization decisions
authorizationDecisions.WithLabelValues(
    string(request.Subject.Role),
    string(request.Action),
    request.Resource.Type,
    strconv.FormatBool(allowed),
).Inc()

// Authorization latency
authorizationDuration.WithLabelValues("unified").Observe(duration.Seconds())
```

### **Audit Logging**

```go
logger.Infow("Authorization decision",
    "user_id", request.Subject.UserID,
    "user_role", request.Subject.Role,
    "action", request.Action,
    "resource_type", request.Resource.Type,
    "resource_id", request.Resource.ID,
    "decision", allowed,
    "duration_ms", duration.Milliseconds(),
)
```

## 🚀 Migration Guide

### **From Existing System**

1. **Phase 1**: Deploy new authorization service alongside existing
2. **Phase 2**: Migrate middleware to use unified service
3. **Phase 3**: Update service layer authorization calls
4. **Phase 4**: Remove legacy authorization code

### **Example Migration**

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

## 🔮 Future Enhancements

### **Planned Features**

- **🤖 ML-based anomaly detection** for access patterns
- **🌐 External policy engines** (OPA, XACML) integration
- **📍 Geolocation-based** access control
- **🔒 Device trust** and network-based policies
- **📊 Advanced analytics** and reporting
- **🎛️ Admin dashboard** for policy management

### **Scalability Improvements**

- **⚡ Redis caching** for policy decisions
- **🔄 Policy hot-reloading** from database/config
- **📊 Distributed policy evaluation** for high-load scenarios
- **🎯 Smart caching** based on access patterns

## 🤝 Contributing

### **Adding New Features**

1. Create feature branch from `main`
2. Implement with comprehensive tests
3. Update documentation
4. Submit PR with clear description

### **Guidelines**

- Follow existing code patterns
- Add tests for all new functionality
- Update documentation
- Consider performance impact
- Ensure backward compatibility

---

## 🎯 Key Benefits

✅ **Performance** - Fast RBAC with cached ABAC decisions  
✅ **Security** - Multi-layer authorization with fine-grained control  
✅ **Scalability** - Modular design supports complex requirements  
✅ **Maintainability** - Clear separation of concerns and documentation  
✅ **Testability** - Comprehensive testing support and examples  
✅ **Flexibility** - Easy to extend with custom policies and providers

This authorization system provides a **robust foundation** for building **secure, scalable applications** with **complex access control requirements**.
