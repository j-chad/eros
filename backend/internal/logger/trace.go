package logger

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"
)

type traceKey struct{}
type spanKey struct{}

func newID(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// WithTrace starts a new trace and attaches the trace ID to the context logger.
// Call this once at the start of a request (e.g. in middleware).
func WithTrace(ctx context.Context) (context.Context, string) {
	traceID := newID(16)
	return WithExistingTrace(ctx, traceID), traceID
}

func WithExistingTrace(ctx context.Context, traceID string) context.Context {
	ctx = context.WithValue(ctx, traceKey{}, traceID)

	l := FromContext(ctx).With("trace_id", traceID)
	return NewContext(ctx, l)
}

// Span represents a timed operation within a trace.
type Span struct {
	Name     string
	TraceID  string
	SpanID   string
	ParentID string
	Start    time.Time
	logger   *slog.Logger
}

// StartSpan begins a named span. If there's already an active span on the
// context, the new span becomes its child. The returned context carries the
// new span so that any deeper StartSpan calls will parent to it.
func StartSpan(ctx context.Context, name string) (context.Context, *Span) {
	traceID, _ := ctx.Value(traceKey{}).(string)
	spanID := newID(8)

	var parentID string
	if parent, ok := ctx.Value(spanKey{}).(*Span); ok {
		parentID = parent.SpanID
	}

	attrs := []any{"span", name, "span_id", spanID}
	if parentID != "" {
		attrs = append(attrs, "parent_span_id", parentID)
	}

	l := FromContext(ctx).With(attrs...)

	s := &Span{
		Name:     name,
		TraceID:  traceID,
		SpanID:   spanID,
		ParentID: parentID,
		Start:    time.Now(),
		logger:   l,
	}

	// Push this span onto context so children can find it
	ctx = context.WithValue(ctx, spanKey{}, s)
	ctx = NewContext(ctx, l)

	l.Info("span started")
	return ctx, s
}

func durationMS(start time.Time) string {
	return fmt.Sprintf("%.2f", float64(time.Since(start).Microseconds())/1000)
}

// DefaultCollector is the global span collector.
// If non-nil, spans will be recorded to it when they end. This is set up by Init() if the collector is enabled in config.
var DefaultCollector *Collector

// End completes the span and logs its duration.
func (s *Span) End(attrs ...any) {
	args := append([]any{"duration_ms", durationMS(s.Start)}, attrs...)
	s.logger.Info("span ended", args...)
	if DefaultCollector != nil {
		DefaultCollector.Record(s, nil)
	}
}

// EndWithError completes the span, logging the error if non-nil.
func (s *Span) EndWithError(err error) {
	dur := durationMS(s.Start)
	if err != nil {
		s.logger.Error("span failed", "duration_ms", dur, "err", err)
	} else {
		s.logger.Info("span ended", "duration_ms", dur)
	}
	if DefaultCollector != nil {
		DefaultCollector.Record(s, err)
	}
}
