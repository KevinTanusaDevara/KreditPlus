package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

const (
	requestsPerMinute = 10
	burstLimit        = 5
	expirationTime    = 5 * time.Minute
)

var limiters = make(map[string]*rate.Limiter)
var lastAccess = make(map[string]time.Time)
var mutex = sync.Mutex{}

func getLimiter(ip string) *rate.Limiter {
	mutex.Lock()
	defer mutex.Unlock()

	now := time.Now()
	for k, v := range lastAccess {
		if now.Sub(v) > expirationTime {
			delete(limiters, k)
			delete(lastAccess, k)
		}
	}

	if _, exists := limiters[ip]; !exists {
		limiters[ip] = rate.NewLimiter(rate.Every(time.Minute/requestsPerMinute), burstLimit)
	}

	lastAccess[ip] = now

	return limiters[ip]
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getLimiter(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}
