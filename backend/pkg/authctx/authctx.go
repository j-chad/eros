package authctx

import "context"

type ctxKey int

const (
	adminKey ctxKey = iota
	deviceIDKey
)

func WithAdmin(ctx context.Context) context.Context {
	return context.WithValue(ctx, adminKey, true)
}

func IsAdmin(ctx context.Context) bool {
	v, ok := ctx.Value(adminKey).(bool)
	return ok && v
}

func WithDeviceID(ctx context.Context, deviceID string) context.Context {
	return context.WithValue(ctx, deviceIDKey, deviceID)
}

func DeviceID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(deviceIDKey).(string)
	return v, ok && v != ""
}
