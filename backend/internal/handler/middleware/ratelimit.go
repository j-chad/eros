package middleware

import (
	"backend/internal/logging"
	"backend/internal/scheduler"
	"backend/pkg/apierror"
	"backend/pkg/response"
	"context"
	"net"
	"net/http"
	"sync"
	"time"
)

type tokenBucket struct {
	tokens   int
	capacity int
	refillAt time.Time
	mu       sync.Mutex
}

func (b *tokenBucket) allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	// Refill tokens every minute
	if now.After(b.refillAt) {
		b.tokens = b.capacity
		b.refillAt = now.Add(time.Minute)
	}

	if b.tokens > 0 {
		b.tokens--
		return true
	}

	return false
}

// IPRateLimit is a per-IP rate limiter using token buckets.
type IPRateLimit struct {
	buckets  map[string]*tokenBucket
	capacity int
	mu       sync.RWMutex
}

func NewIPRateLimit(sched *scheduler.Scheduler, capacity int) *IPRateLimit {
	limiter := &IPRateLimit{
		buckets:  make(map[string]*tokenBucket),
		capacity: capacity,
	}

	sched.MustAddTask(scheduler.Task{
		Name:    "rate_limit_cleanup",
		Fn:      limiter.clean,
		Timeout: 10 * time.Second,
		Cron:    scheduler.MustParseCronExpression("0 * * * *"), // every hour
	})

	return limiter
}

func (r *IPRateLimit) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		logger := logging.FromContext(ctx)

		ip, _, _ := net.SplitHostPort(req.RemoteAddr)
		if ip == "" {
			ip = req.RemoteAddr
		}

		r.mu.Lock()
		bucket, exists := r.buckets[ip]
		if !exists {
			bucket = &tokenBucket{
				tokens:   r.capacity,
				capacity: r.capacity,
				refillAt: time.Now().Add(time.Minute),
			}
			r.buckets[ip] = bucket
		}
		r.mu.Unlock()

		if !bucket.allow() {
			logger.WarnContext(ctx, "rate limit exceeded for IP", "ip", ip)
			response.Error(ctx, w, apierror.TooManyRequests("Too many requests. Try again shortly."))
			return
		}

		next.ServeHTTP(w, req)
	})
}

// clean removes expired token buckets to prevent unbounded memory growth.
// It should be run periodically by the scheduler.
func (r *IPRateLimit) clean(ctx context.Context) error {
	logger := logging.FromContext(ctx)

	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for ip, bucket := range r.buckets {
		if now.After(bucket.refillAt) {
			logger.DebugContext(ctx, "removing expired token bucket for IP", "ip", ip)
			delete(r.buckets, ip)
		}
	}

	return nil
}

type PerNodeRateLimit struct {
	buckets  map[string]*tokenBucket
	capacity int
	mu       sync.RWMutex
}

func NewPerNodeRateLimit(sched *scheduler.Scheduler, capacity int) *PerNodeRateLimit {
	limiter := &PerNodeRateLimit{
		buckets:  make(map[string]*tokenBucket),
		capacity: capacity,
	}

	sched.MustAddTask(scheduler.Task{
		Name:    "rate_limit_cleanup",
		Fn:      limiter.clean,
		Timeout: 10 * time.Second,
		Cron:    scheduler.MustParseCronExpression("0 * * * *"), // every hour
	})

	return limiter
}

func (r *PerNodeRateLimit) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		logger := logging.FromContext(ctx)

		nodeID := req.PathValue("id")
		if nodeID == "" {
			// No node ID in path, skip rate limiting
			logger.DebugContext(ctx, "no node ID in path, skipping rate limit")
			next.ServeHTTP(w, req)
			return
		}

		r.mu.Lock()
		bucket, exists := r.buckets[nodeID]
		if !exists {
			bucket = &tokenBucket{
				tokens:   r.capacity,
				capacity: r.capacity,
				refillAt: time.Now().Add(time.Minute),
			}
			r.buckets[nodeID] = bucket
		}
		r.mu.Unlock()

		if !bucket.allow() {
			logger.WarnContext(ctx, "rate limit exceeded for node ID", "nodeID", nodeID)
			response.Error(ctx, w, apierror.TooManyRequests("Too many unlock attempts. Try again shortly."))
			return
		}

		next.ServeHTTP(w, req)
	})
}

func (r *PerNodeRateLimit) clean(ctx context.Context) error {
	logger := logging.FromContext(ctx)

	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for nodeID, bucket := range r.buckets {
		if now.After(bucket.refillAt) {
			logger.DebugContext(ctx, "removing expired token bucket for node ID", "nodeID", nodeID)
			delete(r.buckets, nodeID)
		}
	}

	return nil
}
