# E-commerce Backend Service

A Go-based e-commerce backend service built using clean architecture principles.

## Table of Contents

- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Development Guide](#development-guide)
- [Error Handling](#error-handling)
- [Adding New Features](#adding-new-features)
- [Best Practices](#best-practices)

## Project Structure

```
├── internal/
│   ├── api/            # HTTP handlers and middleware
│   │   ├── middleware/ # HTTP middleware components
│   │   ├── router/    # Router setup and configuration
│   │   └── v1/        # API version 1 handlers
│   ├── bootstrap/     # Application bootstrapping
│   ├── config/        # Configuration management
│   ├── domain/        # Business domain models and interfaces
│   ├── pkg/           # Shared utilities
│   ├── repository/    # Data access implementations
│   ├── service/       # Business logic implementation
│   └── validator/     # Request validation
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- MongoDB
- Git

### Setup Steps

1. **Clone the Repository**

```bash
git clone https://github.com/your-username/ecommerce.git
cd ecommerce
```

2. **Environment Setup**
   Create a `.env` file in the root directory:

```env
PORT=8080
MODE=development
DB_USER=your_db_user
DB_USER_PWD=your_db_password
DB_HOST=your_db_host
DB_NAME=your_db_name
JWT_SECRET=your_jwt_secret
```

3. **Install Dependencies**

```bash
go mod download
```

4. **Run the Application**

```bash
go run cmd/main.go
```

## Development Guide

### Architecture Overview

This project follows Clean Architecture principles with three main layers:

1. **Domain Layer** (`internal/domain/`)

   - Contains business models and interfaces
   - Defines repository interfaces
   - Contains domain-specific errors

2. **Service Layer** (`internal/service/`)

   - Implements business logic
   - Orchestrates data flow between repositories
   - Handles error wrapping

3. **API Layer** (`internal/api/`)
   - Handles HTTP requests/responses
   - Manages request validation
   - Uses middleware for cross-cutting concerns

### Error Handling

The project implements a hierarchical error handling system:

1. **Domain Errors** (`internal/domain/errors.go`)

   - Base error types for business rules
   - Used across all layers

2. **Service Errors** (`internal/service/errors.go`)

   - Wraps domain errors with context
   - Adds operation and entity information

3. **API Error Handling** (`internal/api/middleware/error_handler.go`)
   - Converts errors to HTTP responses
   - Maintains consistent error format

## Adding New Features

### Step 1: Domain Layer

1. Create a new directory in `internal/domain/<feature>/`
2. Create `model.go`:

```go
package feature

type YourModel struct {
    ID          string    `json:"id" bson:"_id,omitempty"`
    Name        string    `json:"name" bson:"name"`
    CreatedAt   time.Time `json:"created_at" bson:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}
```

3. Create `repository.go`:

```go
package feature

type Repository interface {
    Create(ctx context.Context, model *YourModel) error
    FindByID(ctx context.Context, id string) (*YourModel, error)
    // Add other methods as needed
}
```

### Step 2: Repository Implementation

1. Create `internal/repository/<feature>/repository.go`:

```go
package feature

import (
    "context"
    "github.com/your-username/ecommerce/internal/domain/feature"
    "go.mongodb.org/mongo-driver/mongo"
)

type repository struct {
    db *mongo.Database
}

func NewRepository(db *mongo.Database) feature.Repository {
    return &repository{db: db}
}

// Implement the interface methods
```

### Step 3: Service Layer

1. Create `internal/service/<feature>.go`:

```go
package service

import (
    "context"
    "github.com/your-username/ecommerce/internal/domain/feature"
)

type FeatureService struct {
    repo feature.Repository
}

func NewFeatureService(repo feature.Repository) *FeatureService {
    return &FeatureService{repo: repo}
}

// Add service methods with business logic
```

### Step 4: API Handler

1. Create `internal/api/v1/<feature>.go`:

```go
package v1

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

type FeatureHandler struct {
    service *service.FeatureService
}

func NewFeatureHandler(service *service.FeatureService) *FeatureHandler {
    return &FeatureHandler{service: service}
}

// Add handler methods
```

### Step 5: Update Router

Add your new routes in `internal/api/router/router.go`:

```go
feature := v1.Group("/feature")
{
    feature.POST("/", handlers.CreateFeature)
    feature.GET("/:id", handlers.GetFeature)
    // Add other routes
}
```

## Best Practices

### Code Organization

- Keep packages small and focused
- Use meaningful package names
- Follow Go naming conventions
- Use interfaces for abstraction

### Error Handling

- Use domain errors for business rule violations
- Wrap errors with context at service layer
- Return appropriate HTTP status codes

### Testing

1. **Unit Tests**

   - Test domain logic
   - Test service layer business rules
   - Use mocks for dependencies

2. **Integration Tests**
   - Test repository implementations
   - Test API endpoints

### Validation

- Implement request validation using validator package
- Add domain-level validation in models
- Add service-level business rule validation

### Documentation

- Add godoc comments to exported types and functions
- Keep README updated with new features
- Document API endpoints

## Contributing

1. Create a new branch for your feature

```bash
git checkout -b feature/your-feature-name
```

2. Follow the feature implementation steps above

3. Write tests for your changes

4. Submit a pull request

## License

[Your License Here]

# MongoDB Database Utilities

This package provides a set of utilities and extensions for working with MongoDB in a more streamlined way, while still leveraging the official MongoDB Go driver.

## Features

- **FilterBuilder**: Create MongoDB filters with chainable methods
- **PipelineBuilder**: Build MongoDB aggregation pipelines with a fluent API
- **QueryOptions**: Simplify pagination, sorting, and other query options
- **ExtendedRepository**: Enhanced repository with additional helper methods

## Usage Examples

### Filter Builder

```go
// Create a filter with multiple conditions
filter := mongoUtils.NewFilter().
    Eq("status", "active").
    Gt("age", 18).
    In("roles", []string{"user", "admin"}).
    Build()

// Use the filter with your repository
users, err := userRepo.Find(ctx, filter)
```

### Pipeline Builder

```go
// Create an aggregation pipeline
pipeline := mongoUtils.NewPipeline().
    Match(bson.M{"status": "active"}).
    Lookup("departments", "department_id", "_id", "department").
    Unwind("$department", true).
    Sort(bson.D{{"created_at", -1}}).
    Skip(10).
    Limit(20).
    Build()

// Convert to mongo.Pipeline and use with repository
mongoPipeline := mongoUtils.ToPipeline(pipeline)
var results []UserWithDepartment
err := repo.AggregateWithPipeline(ctx, mongoPipeline, &results)
```

### Query Options

```go
// Create query options for pagination and sorting
opts := mongoUtils.QueryOptions{
    Page:  2,
    Limit: 25,
    Sort:  bson.D{{"created_at", -1}},
}

// Use with repository
var users []User
err := repo.FindWithQueryOptions(ctx, filter, opts, &users)
```

### Extended Repository

```go
// Create an extended repository
repo := mongoRepository.NewExtendedRepository(db.Collection("users"))

// Use extended methods
totalCount, err := repo.PaginatedFind(ctx, filter, 1, 20, &users)

// Create indexes easily
_, err := repo.CreateIndex(ctx,
    bson.D{{"email", 1}},
    options.Index().SetUnique(true))
```

## Hybrid Approach

This package supports a hybrid approach to database operations:

1. Use the utility helpers for common patterns and simplified queries
2. Fall back to raw MongoDB queries for complex operations
3. Extend the repositories as needed for domain-specific operations

This gives you the best of both worlds - convenience for common tasks and full flexibility when needed.
# ksp-go-backend
# codegeeky-lms-backend
