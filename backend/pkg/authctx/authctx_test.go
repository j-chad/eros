package authctx

import (
	"backend/internal/testutil"
	"context"
	"testing"
)

func TestIsAdmin_Default(t *testing.T) {
	testutil.False(t, IsAdmin(context.Background()), "fresh context should not be admin")
}

func TestIsAdmin_WithAdmin(t *testing.T) {
	ctx := WithAdmin(context.Background())
	testutil.True(t, IsAdmin(ctx), "context with WithAdmin should be admin")
}

func TestIsAdmin_NestedContext(t *testing.T) {
	ctx := WithAdmin(context.Background())
	child := context.WithValue(ctx, "other", "value")
	testutil.True(t, IsAdmin(child), "child context should inherit admin status")
}
