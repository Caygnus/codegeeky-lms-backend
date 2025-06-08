package webhookdto

import "github.com/omkar273/codegeeky/internal/domain/user"

type InternalUserEvent struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type UserWebhookPayload struct {
	user.User
}

func NewUserWebhookPayload(u *user.User) *UserWebhookPayload {
	return &UserWebhookPayload{
		User: *u,
	}
}
