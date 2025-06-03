# ABAC (Attribute-Based Access Control) System

## Overview

The ABAC system provides **fine-grained, context-aware authorization** through dynamic policy evaluation. It works **in conjunction with RBAC** to enable complex access control scenarios based on user attributes, resource properties, environmental conditions, and business rules.

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Attributes    │    │     Policies     │    │    Decision     │
│                 │    │                  │    │                 │
│ • User attrs    │────│ • Enrollment     │────│   Allow/Deny    │
│ • Resource      │    │ • Ownership      │    │                 │
│ • Environment   │    │ • Time-based     │    │  + Audit Log    │
│ • Context       │    │ • Progress       │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

## Key Features

### 1. **Dynamic Policy Evaluation**

- **Runtime policy execution** based on current context
- **Multiple policy support** with configurable combination logic
- **Rich attribute model** for complex decision making

### 2. **Policy Types**

- **Enrollment Policy**: Students access only enrolled content
- **Ownership Policy**: Users modify only their own resources
- **Time-Based Policy**: Content available within time windows
- **Progress Policy**: Access gated by completion requirements

### 3. **Flexible Architecture**

- **Plugin-based policies** - easy to add new policy types
- **Configurable policy chains** - customize evaluation order
- **Extensible attribute providers** - load attributes from any source

## Core Concepts

### 1. **Access Request**

Every authorization request contains:

```go
type AccessRequest struct {
    Subject  *AuthContext    // Who is requesting access
    Resource *Resource       // What resource is being accessed
    Action   Permission      // What action is being performed
    Context  map[string]interface{} // Additional context
}
```

### 2. **Policies**

Policies implement business rules:

```go
type Policy interface {
    Evaluate(ctx context.Context, request *AccessRequest) (Decision, error)
    GetName() string
    GetPriority() int
}

type Decision struct {
    Allow  bool
    Reason string
    Metadata map[string]interface{}
}
```

### 3. **Attribute Providers**

Load dynamic attributes:

```go
type AttributeProvider interface {
    LoadUserAttributes(ctx context.Context, userID string) (map[string]interface{}, error)
    LoadResourceAttributes(ctx context.Context, resourceType, resourceID string) (map[string]interface{}, error)
}
```

## Built-in Policies

### 1. **Enrollment Based Access Policy**

**Purpose**: Students can only access content of internships they're enrolled in.

**Attributes Required**:

- `Subject.Attributes["enrolled_internships"]` - List of internship IDs
- `Resource.Attributes["internship_id"]` - Target internship ID

**Logic**:

```go
// Applies only to students accessing content
if user.Role != Student { return Allow }

// Check if accessing content
if !isContentPermission(action) { return Allow }

// Verify enrollment
enrolledInternships := user.Attributes["enrolled_internships"]
targetInternship := resource.Attributes["internship_id"]

if targetInternship in enrolledInternships {
    return Allow
}
return Deny
```

**Example Usage**:

```go
// Student enrolled in internships ["internship_123", "internship_456"]
// Accessing lecture in internship_123 -> ALLOW
// Accessing lecture in internship_789 -> DENY
```

### 2. **Ownership Policy**

**Purpose**: Users can only modify resources they created.

**Attributes Required**:

- `Resource.Attributes["created_by"]` - Resource owner ID
- `Subject.UserID` - Current user ID

**Logic**:

```go
// Applies only to modify actions (update, delete)
if !isModifyAction(action) { return Allow }

// Admins can modify anything
if user.Role == Admin { return Allow }

// Check ownership
resourceOwner := resource.Attributes["created_by"]
return resourceOwner == user.UserID
```

### 3. **Time-Based Access Policy**

**Purpose**: Control access based on time windows (e.g., exam periods, course schedules).

**Attributes Required**:

- `Resource.Attributes["access_start_time"]` - When access begins
- `Resource.Attributes["access_end_time"]` - When access ends

**Logic**:

```go
now := time.Now()

if hasStartTime && now.Before(startTime) {
    return Deny("Access not yet available")
}

if hasEndTime && now.After(endTime) {
    return Deny("Access period has expired")
}

return Allow
```

### 4. **Progress-Based Policy**

**Purpose**: Gate content access based on completion of prerequisites.

**Attributes Required**:

- `Subject.Attributes["progress"]` - User's progress score (0-100)
- `Resource.Attributes["required_progress"]` - Minimum required progress

**Logic**:

```go
// Applies only to students
if user.Role != Student { return Allow }

userProgress := user.Attributes["progress"]
requiredProgress := resource.Attributes["required_progress"]

if userProgress >= requiredProgress {
    return Allow
}

return Deny("Insufficient progress")
```

## Usage Examples

### 1. **Basic Policy Evaluation**

```go
abacService := abac.NewService(logger)

request := &auth.AccessRequest{
    Subject: &auth.AuthContext{
        UserID: "user_123",
        Role:   types.UserRoleStudent,
        Attributes: map[string]interface{}{
            "enrolled_internships": []string{"internship_456"},
            "progress": 75.0,
        },
    },
    Resource: &auth.Resource{
        Type: "content",
        ID:   "lecture_789",
        Attributes: map[string]interface{}{
            "internship_id": "internship_456",
            "required_progress": 50.0,
        },
    },
    Action: auth.PermissionViewLectures,
}

allowed, err := abacService.Evaluate(ctx, request)
// Returns: true (student enrolled + sufficient progress)
```

### 2. **Custom Policy Registration**

```go
// Create custom policy
type CompanyAccessPolicy struct{}

func (p *CompanyAccessPolicy) Evaluate(ctx context.Context, request *auth.AccessRequest) (abac.Decision, error) {
    userCompany := request.Subject.Attributes["company_id"]
    resourceCompany := request.Resource.Attributes["created_by_company"]

    if userCompany == resourceCompany {
        return abac.Decision{Allow: true, Reason: "Same company access"}, nil
    }

    return abac.Decision{Allow: false, Reason: "Cross-company access denied"}, nil
}

// Register policy
abacService.RegisterPolicy(&CompanyAccessPolicy{})
```

### 3. **Attribute Provider Integration**

```go
type EnrollmentAttributeProvider struct {
    enrollmentRepo EnrollmentRepository
}

func (p *EnrollmentAttributeProvider) LoadUserAttributes(ctx context.Context, userID string) (map[string]interface{}, error) {
    enrollments, err := p.enrollmentRepo.GetActiveByUserID(ctx, userID)
    if err != nil {
        return nil, err
    }

    internshipIDs := make([]string, len(enrollments))
    for i, enrollment := range enrollments {
        internshipIDs[i] = enrollment.InternshipID
    }

    return map[string]interface{}{
        "enrolled_internships": internshipIDs,
        "enrollment_dates": extractDates(enrollments),
        "progress": calculateProgress(enrollments),
    }, nil
}

// Register provider
abacService.RegisterAttributeProvider(provider)
```

## Policy Combination Strategies

### 1. **All Must Allow (Default)**

All policies must return `Allow` for access to be granted.

```go
type AllMustAllowCombiner struct{}

func (c *AllMustAllowCombiner) Combine(decisions []Decision) Decision {
    for _, decision := range decisions {
        if !decision.Allow {
            return decision // Return first denial
        }
    }
    return Decision{Allow: true, Reason: "All policies allowed"}
}
```

### 2. **Any Can Allow**

If any policy allows access, grant access.

```go
type AnyCanAllowCombiner struct{}

func (c *AnyCanAllowCombiner) Combine(decisions []Decision) Decision {
    for _, decision := range decisions {
        if decision.Allow {
            return decision // Return first allow
        }
    }
    return Decision{Allow: false, Reason: "All policies denied"}
}
```

### 3. **Priority-Based**

Policies with higher priority override lower priority decisions.

```go
type PriorityBasedCombiner struct{}

func (c *PriorityBasedCombiner) Combine(decisions []Decision) Decision {
    // Sort by priority and return highest priority decision
    sort.Slice(decisions, func(i, j int) bool {
        return decisions[i].Priority > decisions[j].Priority
    })
    return decisions[0]
}
```

## Performance Considerations

### 1. **Caching Strategies**

```go
// Cache attribute lookups
type CachedAttributeProvider struct {
    provider AttributeProvider
    cache    map[string]CacheEntry
    ttl      time.Duration
}

// Cache policy decisions for identical requests
type PolicyDecisionCache struct {
    cache map[string]CachedDecision
    ttl   time.Duration
}
```

### 2. **Lazy Evaluation**

```go
// Only load attributes when needed
func (s *service) Evaluate(ctx context.Context, request *AccessRequest) (bool, error) {
    for _, policy := range s.policies {
        // Check if policy applies before loading attributes
        if !policy.Applies(request) {
            continue
        }

        // Load required attributes only for applicable policies
        if err := s.loadRequiredAttributes(ctx, request, policy); err != nil {
            return false, err
        }

        decision, err := policy.Evaluate(ctx, request)
        // ... handle decision
    }
}
```

### 3. **Batch Operations**

```go
// Batch attribute loading for multiple requests
func (s *service) EvaluateBatch(ctx context.Context, requests []*AccessRequest) ([]Decision, error) {
    // Group requests by required attributes
    attributeRequests := s.groupByAttributes(requests)

    // Batch load attributes
    attributes, err := s.attributeProvider.LoadBatch(ctx, attributeRequests)
    if err != nil {
        return nil, err
    }

    // Evaluate all requests with loaded attributes
    decisions := make([]Decision, len(requests))
    for i, request := range requests {
        s.enrichWithAttributes(request, attributes)
        decisions[i], _ = s.evaluateSingle(ctx, request)
    }

    return decisions, nil
}
```

## Testing

### 1. **Policy Unit Tests**

```go
func TestEnrollmentPolicy(t *testing.T) {
    policy := &EnrollmentBasedAccessPolicy{}

    // Test enrolled student access
    request := createTestRequest(
        userRole: types.UserRoleStudent,
        enrolledInternships: []string{"internship_123"},
        targetInternship: "internship_123",
        action: auth.PermissionViewLectures,
    )

    decision, err := policy.Evaluate(ctx, request)
    assert.NoError(t, err)
    assert.True(t, decision.Allow)

    // Test non-enrolled student access
    request.Resource.Attributes["internship_id"] = "internship_456"
    decision, err = policy.Evaluate(ctx, request)
    assert.NoError(t, err)
    assert.False(t, decision.Allow)
}
```

### 2. **Integration Tests**

```go
func TestABACIntegration(t *testing.T) {
    // Setup service with all policies
    service := abac.NewService(logger)
    service.RegisterPolicy(&EnrollmentBasedAccessPolicy{})
    service.RegisterPolicy(&OwnershipPolicy{})

    // Test complex scenario: student accessing own progress in enrolled course
    request := &auth.AccessRequest{
        Subject: &auth.AuthContext{
            UserID: "student_123",
            Role:   types.UserRoleStudent,
            Attributes: map[string]interface{}{
                "enrolled_internships": []string{"internship_456"},
            },
        },
        Resource: &auth.Resource{
            Type: "progress",
            Attributes: map[string]interface{}{
                "internship_id": "internship_456",
                "created_by": "student_123",
            },
        },
        Action: auth.PermissionViewProgress,
    }

    allowed, err := service.Evaluate(ctx, request)
    assert.NoError(t, err)
    assert.True(t, allowed)
}
```

## Future Extensions

### 1. **Machine Learning Integration**

- **Behavioral analysis** for anomaly detection
- **Risk scoring** based on access patterns
- **Adaptive policies** that learn from user behavior

### 2. **External Policy Engines**

- **OPA (Open Policy Agent)** integration
- **XACML** policy support
- **Cloud provider** policy services

### 3. **Advanced Policy Types**

- **Geolocation-based** access control
- **Device trust** policies
- **Network-based** restrictions
- **Compliance** policy frameworks

---

**Next: See [Integration Guide](../integration/README.md) for combining RBAC and ABAC.**
