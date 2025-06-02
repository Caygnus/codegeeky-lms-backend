package user

import (
	"github.com/omkar273/police/ent"
	"github.com/omkar273/police/internal/types"
	"github.com/samber/lo"
)

type User struct {
	ID       string           `json:"id" db:"id"`
	Email    string           `json:"email" db:"email"`
	Phone    string           `json:"phone" db:"phone"`
	Roles    []types.UserRole `json:"roles" db:"roles"`
	FullName string           `json:"full_name" db:"full_name"`
	types.BaseModel
}

func FromEnt(user *ent.User) *User {
	return &User{
		ID:       user.ID,
		Email:    user.Email,
		Phone:    user.PhoneNumber,
		FullName: user.FullName,
		Roles: lo.Map(user.Roles, func(role string, _ int) types.UserRole {
			return types.UserRole(role)
		}),
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
