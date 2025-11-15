package main

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	db, err := gorm.Open(sqlite.Open("payments.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect SQLite:", err)
	}

	db.AutoMigrate(&Order{}, &Payment{}, &AuditLog{})

	DB = db
}
