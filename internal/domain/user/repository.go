package user

import (
	"context"

	"github.com/omkar273/codegeeky/internal/types"
)

// Repository defines the interface for user persistence operations
type Repository interface {
	// Core operations
	Create(ctx context.Context, user *User) error
	Get(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)

	// List operations
	List(ctx context.Context, filter *types.UserFilter) ([]*User, error)
	ListAll(ctx context.Context, filter *types.UserFilter) ([]*User, error)
	Count(ctx context.Context, filter *types.UserFilter) (int, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, user *User) error
}
