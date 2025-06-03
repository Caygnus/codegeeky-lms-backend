# RBAC (Role-Based Access Control) System

## Overview

The RBAC system provides **simple, efficient role-based authorization** by mapping user roles to specific permissions. It serves as the **primary authorization mechanism** with fast lookups and clear permission boundaries.

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   User Roles    │────│ Role-Permission  │────│   Permissions   │
│                 │    │    Mapping       │    │                 │
│ • Admin         │    │                  │    │ • Create        │
│ • Instructor    │    │   (RBAC Core)    │    │ • Read          │
│ • Student       │    │                  │    │ • Update        │
└─────────────────┘    └──────────────────┘    │ • Delete        │
                                               │ • Publish       │
                                               └─────────────────┘
```

## Key Features

### 1. **Static Role-Permission Mapping**

- **Fast O(1) lookups** for permission checks
- **Predefined roles** with fixed permission sets
- **Clear security boundaries** - no dynamic permission assignment

### 2. **Hierarchical Permissions**

- **Admin**: Full system access (all permissions)
- **Instructor**: Content management + limited admin functions
- **Student**: Read-only access to enrolled content

### 3. **Permission Categories**

- **Internship Management**: Create, update, delete, view, publish
- **Content Access**: View lectures, assignments, resources
- **User Management**: Manage users, view all users
- **System Operations**: Configuration, analytics

## Implementation

### Role Definition

```go
type UserRole string

const (
    UserRoleStudent    UserRole = "STUDENT"
    UserRoleInstructor UserRole = "INSTRUCTOR"
    UserRoleAdmin      UserRole = "ADMIN"
)
```

### Permission Definition

```go
type Permission string

const (
    // Internship permissions
    PermissionCreateInternship  Permission = "internship:create"
    PermissionUpdateInternship  Permission = "internship:update"
    PermissionDeleteInternship  Permission = "internship:delete"
    PermissionViewInternship    Permission = "internship:view"
    PermissionPublishInternship Permission = "internship:publish"

    // Content permissions
    PermissionViewLectures      Permission = "content:lectures:view"
    PermissionViewAssignments   Permission = "content:assignments:view"
    PermissionViewResources     Permission = "content:resources:view"
    PermissionDownloadContent   Permission = "content:download"

    // Admin permissions
    PermissionManageUsers       Permission = "users:manage"
    PermissionViewAllUsers      Permission = "users:view:all"
    PermissionSystemConfig      Permission = "system:config"
    PermissionViewAnalytics     Permission = "analytics:view"
)
```

## Usage Examples

### 1. **Basic Permission Check**

```go
rbacService := rbac.NewService()

// Check if instructor can create internships
hasPermission := rbacService.HasPermission(
    types.UserRoleInstructor,
    auth.PermissionCreateInternship,
)
// Returns: true
```

### 2. **Get All User Permissions**

```go
permissions := rbacService.GetUserPermissions(types.UserRoleStudent)
// Returns: [PermissionViewInternship]
```

### 3. **Middleware Integration**

```go
// Role-based middleware
middleware.RequireRole(types.UserRoleInstructor, types.UserRoleAdmin)

// Permission-based middleware
middleware.RequirePermission(auth.PermissionCreateInternship)
```

## Permission Matrix

| Permission                 | Admin | Instructor | Student |
| -------------------------- | ----- | ---------- | ------- |
| **Internship Management**  |
| `internship:create`        | ✅    | ✅         | ❌      |
| `internship:update`        | ✅    | ✅\*       | ❌      |
| `internship:delete`        | ✅    | ✅\*       | ❌      |
| `internship:view`          | ✅    | ✅         | ✅\*\*  |
| `internship:publish`       | ✅    | ✅\*       | ❌      |
| **Content Access**         |
| `content:lectures:view`    | ✅    | ✅         | ✅\*\*  |
| `content:assignments:view` | ✅    | ✅         | ✅\*\*  |
| `content:resources:view`   | ✅    | ✅         | ✅\*\*  |
| `content:download`         | ✅    | ✅         | ✅\*\*  |
| **User Management**        |
| `users:manage`             | ✅    | ❌         | ❌      |
| `users:view:all`           | ✅    | ❌         | ❌      |
| **System Operations**      |
| `system:config`            | ✅    | ❌         | ❌      |
| `analytics:view`           | ✅    | ✅\*\*\*   | ❌      |

**Notes:**

- `*` Instructors: Limited by ABAC policies (own content only)
- `**` Students: Limited by ABAC policies (enrolled content only)
- `***` Instructors: Own analytics only

## Best Practices

### 1. **Role Design**

- **Keep roles simple** - avoid complex hierarchies
- **Use meaningful names** that reflect business roles
- **Limit role count** - typically 3-7 roles maximum

### 2. **Permission Granularity**

- **Resource-action format**: `resource:action` (e.g., `internship:create`)
- **Specific enough** for security, **general enough** for maintainability
- **Avoid over-segmentation** - use ABAC for fine-grained control

### 3. **Performance**

- **Static mapping** enables fast in-memory lookups
- **Cache-friendly** - permissions rarely change
- **O(1) complexity** for permission checks

## Integration with ABAC

RBAC works **in conjunction** with ABAC:

1. **RBAC Check First**: Fast role-based permission validation
2. **ABAC Refinement**: Context-aware policy evaluation
3. **Combined Result**: Both must allow access

```go
// Authorization flow
func (s *authService) IsAuthorized(request *AccessRequest) (bool, error) {
    // 1. Check RBAC (fast)
    if !s.rbac.HasPermission(request.Subject.Role, request.Action) {
        return false, nil // RBAC denies
    }

    // 2. Check ABAC (contextual)
    return s.abac.Evaluate(ctx, request)
}
```

## Testing

### Unit Tests

```go
func TestRBACPermissions(t *testing.T) {
    service := rbac.NewService()

    // Test admin has all permissions
    assert.True(t, service.HasPermission(
        types.UserRoleAdmin,
        auth.PermissionSystemConfig,
    ))

    // Test student has limited permissions
    assert.False(t, service.HasPermission(
        types.UserRoleStudent,
        auth.PermissionCreateInternship,
    ))
}
```

## Future Extensions

### 1. **Dynamic Role Loading**

- Load roles from database/config
- Hot-reload permission changes
- Role versioning support

### 2. **Role Inheritance**

- Hierarchical role structure
- Permission inheritance chains
- Complex organizational structures

### 3. **Permission Groups**

- Logical permission groupings
- Bulk permission assignment
- Simplified role management

---

**Next: See [ABAC Documentation](../abac/README.md) for fine-grained access control.**
