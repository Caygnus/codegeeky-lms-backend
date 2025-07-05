package testutil

import (
	"context"

	"github.com/omkar273/codegeeky/internal/types"
)

// CtxRequestID     ContextKey = "ctx_request_id"
//
//	CtxUserID        ContextKey = "ctx_user_id"
//	CtxJWT           ContextKey = "ctx_jwt"
//	CtxUserEmail     ContextKey = "ctx_user_email"
//	CtxDBTransaction ContextKey = "ctx_db_transaction"
//	CtxUserRole      ContextKey = "ctx_user_role"
//	CtxAuthContext   ContextKey = "ctx_auth_context"
//	CtxUser          ContextKey = "ctx_user"
func SetupContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, types.CtxUserID, types.DefaultUserID)
	ctx = context.WithValue(ctx, types.CtxRequestID, types.GenerateUUID())
	ctx = context.WithValue(ctx, types.CtxUserEmail, types.DefaultUserEmail)
	ctx = context.WithValue(ctx, types.CtxUserRole, types.DefaultUserRole)
	ctx = context.WithValue(ctx, types.CtxUser, types.DefaultUserID)

	return ctx
}
