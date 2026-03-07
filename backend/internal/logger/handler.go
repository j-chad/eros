package logger

import (
	"context"
	"log/slog"
)

// collectingHandler wraps an inner slog.Handler and also forwards log
// records to the Collector when the log comes from within an active span.
type collectingHandler struct {
	inner     slog.Handler
	collector *Collector
	// precomputed attrs from WithAttrs / WithGroup
	attrs []slog.Attr
}

func newCollectingHandler(inner slog.Handler, c *Collector) *collectingHandler {
	return &collectingHandler{inner: inner, collector: c}
}

func (h *collectingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *collectingHandler) Handle(ctx context.Context, rec slog.Record) error {
	// Always delegate to the inner handler (stdout logging)
	if err := h.inner.Handle(ctx, rec); err != nil {
		return err
	}

	// If we're inside a span, also capture the log into the collector
	if s, ok := ctx.Value(spanKey{}).(*Span); ok && h.collector != nil {
		attrs := make(map[string]any)

		// Collect pre-set attrs from WithAttrs
		for _, a := range h.attrs {
			attrs[a.Key] = a.Value.Any()
		}

		// Collect attrs from this specific record
		rec.Attrs(func(a slog.Attr) bool {
			// Skip span-related attrs that are already in the span record
			switch a.Key {
			case "trace_id", "span", "span_id", "parent_span_id":
				return true
			}
			attrs[a.Key] = a.Value.Any()
			return true
		})

		// Remove empty attrs map
		if len(attrs) == 0 {
			attrs = nil
		}

		h.collector.RecordLog(s.TraceID, s.SpanID, rec.Level.String(), rec.Message, attrs)
	}

	return nil
}

func (h *collectingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &collectingHandler{
		inner:     h.inner.WithAttrs(attrs),
		collector: h.collector,
		attrs:     append(h.attrs, attrs...),
	}
}

func (h *collectingHandler) WithGroup(name string) slog.Handler {
	return &collectingHandler{
		inner:     h.inner.WithGroup(name),
		collector: h.collector,
		attrs:     h.attrs,
	}
}
