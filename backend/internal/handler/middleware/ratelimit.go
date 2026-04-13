package middleware

import (
	"backend/pkg/apierror"
	"backend/pkg/response"
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

type PerNodeRateLimit struct {
	buckets  map[string]*tokenBucket
	capacity int
	mu       sync.RWMutex
}

func NewPerNodeRateLimit(capacity int) *PerNodeRateLimit {
	return &PerNodeRateLimit{
		buckets:  make(map[string]*tokenBucket),
		capacity: capacity,
	}
}

func (r *PerNodeRateLimit) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		nodeID := req.PathValue("id")
		if nodeID == "" {
			// No node ID in path, skip rate limiting
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
			response.Error(req.Context(), w, apierror.TooManyRequests("Too many unlock attempts. Try again shortly."))
			return
		}

		next.ServeHTTP(w, req)
	})
}
