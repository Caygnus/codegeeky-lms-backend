package types

import "context"

// ContextKey is a type for the keys of values stored in the context
type ContextKey string

const (
	CtxRequestID     ContextKey = "ctx_request_id"
	CtxUserID        ContextKey = "ctx_user_id"
	CtxJWT           ContextKey = "ctx_jwt"
	CtxDBTransaction ContextKey = "ctx_db_transaction"

	// Default values
	DefaultUserID = "00000000-0000-0000-0000-000000000000"
)

func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(CtxUserID).(string); ok {
		return userID
	}
	return ""
}

func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(CtxRequestID).(string); ok {
		return requestID
	}
	return ""
}

func GetJWT(ctx context.Context) string {
	if jwt, ok := ctx.Value(CtxJWT).(string); ok {
		return jwt
	}
	return ""
}
