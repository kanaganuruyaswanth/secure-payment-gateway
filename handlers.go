package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ------------------- CREATE ORDER -------------------
func CreateOrder(c *gin.Context) {
	var req struct {
		CustomerID  string `json:"customer_id"`
		Amount      int64  `json:"amount"`
		Currency    string `json:"currency"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Amount <= 0 || req.Currency != "IDR" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	order := Order{
		ID:          uuid.NewString(),
		CustomerID:  req.CustomerID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Description: req.Description,
		Status:      "INITIATED",
	}

	DB.Create(&order)
	c.JSON(200, order)
}

// ------------------- INITIATE PAYMENT -------------------
func InitiatePayment(c *gin.Context) {
	var req struct {
		OrderID  string `json:"order_id"`
		Method   string `json:"method"`
		BankCode string `json:"bank_code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	if req.Method != "VA" && req.Method != "QRIS" {
		c.JSON(400, gin.H{"error": "invalid payment method"})
		return
	}

	payment := Payment{
		ID:        uuid.NewString(),
		OrderID:   req.OrderID,
		Method:    req.Method,
		Status:    "PENDING",
		CreatedAt: time.Now(),
	}

	if req.Method == "VA" {
		payment.BankCode = req.BankCode
		payment.VANumber = "9001" + time.Now().Format("150405")
	}
	if req.Method == "QRIS" {
		payment.QRPayload = "000201010212" + uuid.NewString()
	}

	DB.Create(&payment)
	c.JSON(200, payment)
}

// ------------------- WEBHOOK -------------------
func PaymentWebhook(c *gin.Context) {
	log.Println("--- STARTING WEBHOOK PROCESSING ---") // New line
	bodyBytes, err := c.GetRawData()
	if err != nil {
		log.Printf("ERROR getting raw data: %v", err)
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	// This logs the exact string used to generate the signature in GenerateHMAC.
	log.Printf("WEBHOOK PAYLOAD RAW BODY: %s", string(bodyBytes))

	var data struct {
		PaymentID string `json:"payment_id"`
		Status    string `json:"status"`
		Timestamp int64  `json:"timestamp"`
	}

	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		log.Printf("ERROR unmarshaling JSON: %v", err)
		c.JSON(400, gin.H{"error": "invalid json format"})
		return
	}

	// Timestamp Check
	if !IsTimestampValid(data.Timestamp) {
		c.JSON(400, gin.H{"error": "replay attack blocked"})
		return
	}

	// HMAC Check
	expected := GenerateHMAC(cfg.WebhookSecret, string(bodyBytes))
	if c.GetHeader("X-Signature") != expected {
		c.JSON(401, gin.H{"error": "invalid signature"})
		return
	}

	var payment Payment
	if err := DB.First(&payment, "id = ?", data.PaymentID).Error; err != nil {
		c.JSON(404, gin.H{"error": "payment not found"})
		return
	}

	payment.Status = data.Status
	DB.Save(&payment)

	c.JSON(200, gin.H{"message": "webhook processed"})
}

// GetOrderStatus handles GET /api/v1/orders/{order_id}
func GetOrderStatus(c *gin.Context) {
	// 1. Get the order_id from the URL path parameter
	orderID := c.Param("order_id")

	// 2. Retrieve the order from the database
	var order Order
	// Assuming DB is initialized and points to your GORM connection
	if err := DB.First(&order, "id = ?", orderID).Error; err != nil {
		// Handle case where order is not found
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{"error": "Order not found"})
			return
		}
		// Handle other database errors
		log.Printf("DB error fetching order %s: %v", orderID, err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	// 3. Construct the response
	// We only want to return relevant status information
	response := gin.H{
		"order_id":   order.ID,
		"amount":     order.Amount,
		"currency":   order.Currency,
		"status":     "SUCCESS",
		"created_at": order.CreatedAt.Format(time.RFC3339),
	}

	// 4. Return 200 OK with the order details
	c.JSON(200, response)
}
