package main

import (
	"net/http"
	"sync"

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
