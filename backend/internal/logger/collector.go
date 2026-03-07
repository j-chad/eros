package logger

import (
	"backend/internal/config"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// LogRecord is a single log line captured within a span.
type LogRecord struct {
	Level   string         `json:"level"`
	Message string         `json:"msg"`
	AtMS    float64        `json:"at_ms"` // offset from trace start
	SpanID  string         `json:"span_id"`
	Attrs   map[string]any `json:"attrs,omitempty"`
}

// SpanRecord is the finished span data we store for visualization.
type SpanRecord struct {
	TraceID  string      `json:"trace_id"`
	SpanID   string      `json:"span_id"`
	ParentID string      `json:"parent_id,omitempty"`
	Name     string      `json:"name"`
	StartMS  float64     `json:"start_ms"` // offset from trace start
	DurMS    float64     `json:"dur_ms"`
	Error    string      `json:"error,omitempty"`
	Logs     []LogRecord `json:"logs,omitempty"`
}

// Collector holds recent completed spans in a ring buffer.
type Collector struct {
	mu    sync.Mutex
	spans []SpanRecord

	// track the earliest start time per trace so we can compute offsets
	traceStarts map[string]time.Time

	// pendingLogs accumulates logs for in-flight spans
	pendingLogs map[string][]LogRecord // keyed by span_id
}

func NewCollector(config config.CollectorConfig) *Collector {
	if config.MaxSpans <= 0 {
		config.MaxSpans = 1000
	}

	return &Collector{
		spans:       make([]SpanRecord, 0, config.MaxSpans),
		traceStarts: make(map[string]time.Time),
		pendingLogs: make(map[string][]LogRecord),
	}
}

// RecordLog captures a log entry for an in-flight span.
func (c *Collector) RecordLog(traceID, spanID, level, msg string, attrs map[string]any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	traceStart, ok := c.traceStarts[traceID]
	if !ok {
		traceStart = time.Now()
	}

	rec := LogRecord{
		Level:   level,
		Message: msg,
		AtMS:    float64(time.Since(traceStart).Microseconds()) / 1000,
		SpanID:  spanID,
		Attrs:   attrs,
	}
	c.pendingLogs[spanID] = append(c.pendingLogs[spanID], rec)
}

// Record adds a finished span to the ring buffer.
func (c *Collector) Record(s *Span, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Track the earliest span start per trace
	if first, ok := c.traceStarts[s.TraceID]; !ok || s.Start.Before(first) {
		c.traceStarts[s.TraceID] = s.Start
	}
	traceStart := c.traceStarts[s.TraceID]

	rec := SpanRecord{
		TraceID:  s.TraceID,
		SpanID:   s.SpanID,
		ParentID: s.ParentID,
		Name:     s.Name,
		StartMS:  float64(s.Start.Sub(traceStart).Microseconds()) / 1000,
		DurMS:    float64(time.Since(s.Start).Microseconds()) / 1000,
		Logs:     c.pendingLogs[s.SpanID],
	}
	delete(c.pendingLogs, s.SpanID)
	if err != nil {
		rec.Error = err.Error()
	}

	if len(c.spans) >= cap(c.spans) {
		// Drop oldest
		c.spans = c.spans[1:]
	}
	c.spans = append(c.spans, rec)
}

// TraceSummary groups spans by trace for the API response.
type TraceSummary struct {
	TraceID string       `json:"trace_id"`
	Spans   []SpanRecord `json:"spans"`
}

func (c *Collector) Traces() []TraceSummary {
	c.mu.Lock()
	defer c.mu.Unlock()

	grouped := make(map[string][]SpanRecord)
	order := make([]string, 0)
	for _, s := range c.spans {
		if _, exists := grouped[s.TraceID]; !exists {
			order = append(order, s.TraceID)
		}
		grouped[s.TraceID] = append(grouped[s.TraceID], s)
	}

	out := make([]TraceSummary, 0, len(grouped))
	// Reverse order so the newest traces appear first
	for i := len(order) - 1; i >= 0; i-- {
		tid := order[i]
		out = append(out, TraceSummary{TraceID: tid, Spans: grouped[tid]})
	}
	return out
}

// ServeHTTP handles GET /debug/traces and GET /debug/traces?id=<trace_id>
func (c *Collector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if tid := r.URL.Query().Get("id"); tid != "" {
		for _, t := range c.Traces() {
			if t.TraceID == tid {
				if err := json.NewEncoder(w).Encode(t); err != nil {
					http.Error(w, "failed to encode trace", http.StatusInternalServerError)
				}
				return
			}
		}
		http.Error(w, "trace not found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(c.Traces()); err != nil {
		http.Error(w, "failed to encode traces", http.StatusInternalServerError)
	}
}
