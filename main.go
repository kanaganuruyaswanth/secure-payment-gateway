package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var cfg *Config

func main() {
	// Load env file
	err := godotenv.Load("env")
	if err != nil {
		log.Println("WARNING: Could not load env file:", err)
	}

	cfg = LoadConfig()
	InitDB()

	r := gin.Default()

	r.Use(RateLimit(cfg.RateLimitPerMinute))

	post := r.Group("/")
	post.Use(IdempotencyMiddleware())
	{
		post.POST("/api/v1/orders", CreateOrder)
		post.POST("/api/v1/payments/initiate", InitiatePayment)
		post.POST("/api/v1/payments/webhook", PaymentWebhook)
	}

	// GET route — NO idempotency required
	r.GET("/api/v1/orders/:order_id", GetOrderStatus)

	// Start server
	r.Run(":8090")
}
