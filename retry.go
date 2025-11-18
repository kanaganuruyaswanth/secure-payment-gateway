package main

import (
	"bytes"
	"log"
	"net/http"
	"time"
)

func SendWebhookWithRetry(url string, body []byte) {
	maxRetries := 3

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err == nil && res.StatusCode == 200 {
			log.Println("Webhook delivered successfully")
			return
		}

		log.Printf("Retry %d failed, backing off...", attempt)
		time.Sleep(time.Duration(attempt) * time.Second)
	}

	log.Println("Webhook delivery failed after all retries")
}
