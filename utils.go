package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func GenerateHMAC(secret, body string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(body))
	return hex.EncodeToString(h.Sum(nil))
}

func IsTimestampValid(ts int64) bool {
	now := time.Now().Unix()
	fmt.Printf("GOT TIMESTAMP %d", now)
	diff := now - ts
	return diff >= 0 && diff <= 300
}
