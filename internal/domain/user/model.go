package user

import (
	"github.com/omkar273/codegeeky/ent"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
)

type User struct {
	ID       string         `json:"id,omitempty" db:"id"`
	Email    string         `json:"email,omitempty" db:"email"`
	Phone    string         `json:"phone,omitempty" db:"phone"`
	Role     types.UserRole `json:"role,omitempty" db:"role"`
	FullName string         `json:"full_name,omitempty" db:"full_name"`
	types.BaseModel
}

func FromEnt(user *ent.User) *User {
	return &User{
		ID:       user.ID,
		Email:    user.Email,
		Phone:    user.PhoneNumber,
		FullName: user.FullName,
		Role:     types.UserRole(user.Role),
		BaseModel: types.BaseModel{
			Status:    types.Status(user.Status),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			CreatedBy: user.CreatedBy,
			UpdatedBy: user.UpdatedBy,
		},
	}
}

func FromEntList(users []*ent.User) []*User {
	return lo.Map(users, func(user *ent.User, _ int) *User {
		return FromEnt(user)
	})
}
