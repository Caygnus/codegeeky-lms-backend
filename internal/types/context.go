package types

import "context"

// ContextKey is a type for the keys of values stored in the context
type ContextKey string

const (
	CtxRequestID     ContextKey = "ctx_request_id"
	CtxUserID        ContextKey = "ctx_user_id"
	CtxJWT           ContextKey = "ctx_jwt"
	CtxUserEmail     ContextKey = "ctx_user_email"
	CtxDBTransaction ContextKey = "ctx_db_transaction"
	CtxUserRole      ContextKey = "ctx_user_role"
	CtxAuthContext   ContextKey = "ctx_auth_context"
	CtxUser          ContextKey = "ctx_user"
	CtxIsGuest       ContextKey = "ctx_is_guest"

	// Default values
	DefaultUserID    = "00000000-0000-0000-0000-000000000000"
	DefaultRequestID = "00000000-0000-0000-0000-000000000000"
	DefaultUserEmail = "demo-user@test.com"
	DefaultUserRole  = UserRoleAdmin
	DefaultIsGuest   = true
)

func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(CtxUserID).(string); ok {
		return userID
	}
	return ""
}

func GetUserRole(ctx context.Context) UserRole {
	if userRole, ok := ctx.Value(CtxUserRole).(UserRole); ok {
		return userRole
	}

	// TODO: This is a temporary solution to get the user role.
	// We need to find a better way to get the user role.
	return UserRole("")
}

func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(CtxRequestID).(string); ok {
		return requestID
	}
	return ""
}

func GetUserEmail(ctx context.Context) string {
	if userEmail, ok := ctx.Value(CtxUserEmail).(string); ok {
		return userEmail
	}
	return ""
}

func GetJWT(ctx context.Context) string {
	if jwt, ok := ctx.Value(CtxJWT).(string); ok {
		return jwt
	}
	return ""
}
