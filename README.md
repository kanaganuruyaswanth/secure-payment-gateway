
This project implements a **secure backend payment gateway simulator"" in ""Golang"", supporting Indonesian Virtual Account (VA) and QRIS payments.  
It follows best practices inspired by "Bank Indonesia (BI)"" and ""OJK"" standards.

---

Features

- Create orders and initiate payments (VA / QRIS)
- Mock bank webhook (with HMAC-SHA256 verification)
- Merchant callback delivery with "retry + exponential backoff"
- "Idempotency-Key" protection for duplicate payments
- **Replay attack protection (timestamp ≤ 5 mins)
- "Rate limiting" using sliding window algorithm
- ""Audit logging"" for all critical transitions
- Secrets stored only in environment variables

---

 API Endpoints

 1. Create Order
`POST /api/v1/orders`

2. Initiate Payment
`POST /api/v1/payments/initiate`

3. Payment Webhook (Bank → Gateway)
`POST /api/v1/payments/webhook`

4. Get Order Status
`GET /api/v1/orders/{order_id}`

---

Security Implemented

| Security Requirement | Status |
|----------------------|--------|
| HMAC verification for webhook | ✔ |
| Replay attack protection | ✔ |
| Idempotency keys | ✔ |
| Input validation (IDR only) | ✔ |
| No sensitive logs | ✔ |
| Rate limiting | ✔ |
| Retry logic for merchant callback | ✔ |
| Env-based secrets | ✔ |

---

 Database Tables

orders
- id  
- customer_id  
- amount  
- currency  
- description  
- status  
- created_at  
- updated_at  

payments
- id  
- order_id  
- method  
- bank_code  
- va_number  
- qris_payload  
- status  
- created_at  
- updated_at  

audit_logs
- id  
- entity_type  
- entity_id  
- action  
- metadata_json  
- created_at  

---

 Data Flow Diagram (DFD)

```
 Customer
    │ POST /orders
    ▼
 Payment Gateway API (Go)
    │ Creates Order + Payment
    │
    ├──► Mock Bank Server
    │       Sends webhook (SUCCESS/FAILED)
    │
    ▼
 Webhook Handler
    │ Validate HMAC + Timestamp
    │ Update Order & Payment
    │ Log Audit Entry
    │ Retry callback to merchant
    ▼
 Merchant Callback URL
```

---

Threat Modeling (Summary)

| Threat | Mitigation |
|--------|------------|
| Replay attacks | Timestamp ≤5 min, HMAC |
| Forged webhooks | HMAC-SHA256 |
| Duplicate payments | Idempotency-Key |
| API brute force | Sliding-window rate limiting |
| Invalid amounts | Input validation |
| Sensitive data exposure | Masked logging |
| Failed merchant updates | Retry logic |

---

 Environment Variables

Create a `.env` file:

```
WEBHOOK_HMAC_SECRET=your-secret
RATE_LIMIT_PER_MINUTE=5
SIMULATE_BANK_CALLBACK=false
MERCHANT_CALLBACK_URL=http://localhost:9999/dummy
```

---

 Run the Server

```bash
go mod tidy
go run main.go
```

Server runs on:

```
http://localhost:8080
```

---

 Full Flow  
1. Create Order  
2. Initiate Payment  
3. Trigger Webhook  
4. Get Order Status  

---

 Project Structure

```
/config
/middleware
/controllers
/models
retry.go
utils.go
main.go
```

---

Summary

This project demonstrates:

- Secure payment processing  
- Strong defensive design  
- Webhook reliability with retry logic  
- Compliance-style implementation  
- Clean and modular Golang codebase  


