package logger

import (
	"backend/internal/config"
	"context"
	"log/slog"
	"os"
)

// Init sets up the default slog logger for the application.
// Call this once at startup (e.g. in main()).
func Init(config config.LoggingConfig) {
	opts := &slog.HandlerOptions{
		Level:     config.Level,
		AddSource: config.AddSource,
	}

	var handler slog.Handler
	if config.JSON {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	if config.Collector.Enabled {
		DefaultCollector = NewCollector(config.Collector)
		handler = newCollectingHandler(handler, DefaultCollector)
	}

	slog.SetDefault(slog.New(handler))
}

// FromContext extracts a logger from context, falling back to the default.
// Pair with NewContext to propagate request-scoped fields.
type ctxKey struct{}

func NewContext(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}
