package authctx

import "context"

type ctxKey int

const adminKey ctxKey = iota

func WithAdmin(ctx context.Context) context.Context {
	return context.WithValue(ctx, adminKey, true)
}

func IsAdmin(ctx context.Context) bool {
	v, ok := ctx.Value(adminKey).(bool)
	return ok && v
}
