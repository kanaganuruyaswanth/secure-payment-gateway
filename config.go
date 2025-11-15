package main

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	WebhookSecret      string
	RateLimitPerMinute int
	SimulateCallback   bool
}

func LoadConfig() *Config {
	fmt.Println("DEBUG SECRET =", os.Getenv("WEBHOOK_HMAC_SECRET"))

	secret := os.Getenv("WEBHOOK_HMAC_SECRET")
	if secret == "" {
		log.Fatal("WEBHOOK_HMAC_SECRET is required")
	}

	return &Config{
		WebhookSecret: secret,
	}

}
