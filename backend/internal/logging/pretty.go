package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"
)

func newPrettyHandler(w io.Writer, level slog.Level) slog.Handler {
	return &prettyHandler{out: w, level: level, mu: &sync.Mutex{}}
}

const (
	reset     = "\033[0m"
	grey      = "\033[90m"
	lightGrey = "\033[37m"
	green     = "\033[32m"
	yellow    = "\033[33m"
	red       = "\033[31m"
)

func levelColor(l slog.Level) string {
	switch {
	case l >= slog.LevelError:
		return red
	case l >= slog.LevelWarn:
		return yellow
	case l >= slog.LevelInfo:
		return green
	default:
		return grey
	}
}

type prettyHandler struct {
	out   io.Writer
	level slog.Level
	mu    *sync.Mutex
	goas  []groupOrAttrs
}

type groupOrAttrs struct {
	group string
	attrs []slog.Attr
}

func (h *prettyHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= h.level
}

func (h *prettyHandler) Handle(_ context.Context, r slog.Record) error {
	buf := make([]byte, 0, 1024)

	color := levelColor(r.Level)
	level := r.Level.String()[:3]
	buf = fmt.Appendf(buf, "%s%s %s%s %s",
		color, r.Time.Format(time.TimeOnly), level, reset, r.Message)

	// Attrs from WithAttrs
	for _, goa := range h.goas {
		for _, a := range goa.attrs {
			buf = appendAttr(buf, a)
		}
	}

	// Attrs from this log call
	r.Attrs(func(a slog.Attr) bool {
		buf = appendAttr(buf, a)
		return true
	})

	buf = append(buf, '\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.out.Write(buf)
	return err
}

func (h *prettyHandler) withGroupOrAttrs(goa groupOrAttrs) *prettyHandler {
	h2 := *h
	h2.goas = make([]groupOrAttrs, len(h.goas)+1)
	copy(h2.goas, h.goas)
	h2.goas[len(h2.goas)-1] = goa
	return &h2
}

func (h *prettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	return h.withGroupOrAttrs(groupOrAttrs{attrs: attrs})
}

func (h *prettyHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return h.withGroupOrAttrs(groupOrAttrs{group: name})
}

func appendAttr(buf []byte, a slog.Attr) []byte {
	a.Value = a.Value.Resolve()
	if a.Equal(slog.Attr{}) {
		return buf
	}
	return fmt.Appendf(buf, " %s%s=%s%v%s", grey, a.Key, lightGrey, a.Value, reset)
}
