package rate

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter controls how frequently requests are allowed for a given key.
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
	r        rate.Limit
	b        int
}

// NewRateLimiter creates a new RateLimiter.
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        r,
		b:        b,
	}
}

// GetLimiter returns a rate limiter for the provided key, creating a new one if necessary.
func (l *RateLimiter) GetLimiter(key string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	limiter, exists := l.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(l.r, l.b)
		l.limiters[key] = limiter
	}
	return limiter
}

// Middleware enforces rate limiting using an IP address as the key.
func (l *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := l.GetLimiter(c.ClientIP())
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			return
		}
		c.Next()
	}
}
