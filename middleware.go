package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var idempotencyStore = sync.Map{}

func IdempotencyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("Idempotency-Key")
		if key == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Idempotency-Key header required"})
			c.Abort()
			return
		}
		if _, exists := idempotencyStore.Load(key); exists {
			c.JSON(http.StatusConflict, gin.H{"error": "Duplicate request"})
			c.Abort()
			return
		}
		idempotencyStore.Store(key, true)
		c.Next()
	}
}

var requestCounts = struct {
	sync.Mutex
	data map[string][]time.Time
}{data: make(map[string][]time.Time)}

func RateLimit(limit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if limit == 0 { // disabled
			c.Next()
			return
		}

		ip := c.ClientIP()
		now := time.Now()

		requestCounts.Lock()

		timestamps := requestCounts.data[ip]
		var filtered []time.Time
		cutoff := now.Add(-1 * time.Minute)

		for _, t := range timestamps {
			if t.After(cutoff) {
				filtered = append(filtered, t)
			}
		}

		if len(filtered) >= limit {
			requestCounts.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			c.Abort()
			return
		}

		filtered = append(filtered, now)
		requestCounts.data[ip] = filtered
		requestCounts.Unlock()

		c.Next()
	}
}
