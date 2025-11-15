package main

import "time"

type Order struct {
	ID          string `gorm:"primaryKey"`
	CustomerID  string
	Amount      int64
	Currency    string
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Payment struct {
	ID        string `gorm:"primaryKey"`
	OrderID   string
	Method    string
	BankCode  string
	VANumber  string
	QRPayload string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AuditLog struct {
	ID         uint `gorm:"primaryKey"`
	EntityType string
	EntityID   string
	Action     string
	Metadata   string
	CreatedAt  time.Time
}
