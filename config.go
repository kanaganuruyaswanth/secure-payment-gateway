package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Config struct {
	WebhookSecret       string
	RateLimitPerMinute  int
	SimulateCallback    bool
	MerchantCallbackURL string
}

func LoadConfig() *Config {
	fmt.Println("DEBUG SECRET =", os.Getenv("WEBHOOK_HMAC_SECRET"))

	secret := os.Getenv("WEBHOOK_HMAC_SECRET")
	if secret == "" {
		log.Fatal("WEBHOOK_HMAC_SECRET is required")
	}

	rateLimitStr := os.Getenv("RATE_LIMIT_PER_MINUTE")
	rateLimit, _ := strconv.Atoi(rateLimitStr)
	if rateLimit == 0 {
		rateLimit = 20
	}

	merchantURL := os.Getenv("MERCHANT_CALLBACK_URL")
	if merchantURL == "" {
		merchantURL = "http://dummy-merchant-server.com/callback"
	}

	return &Config{
		WebhookSecret:       secret,
		RateLimitPerMinute:  rateLimit,
		MerchantCallbackURL: merchantURL,
	}
}
