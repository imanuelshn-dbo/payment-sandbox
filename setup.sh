#!/bin/bash

echo "🚀 Generating FULL Payment Sandbox Project..."

mkdir -p cmd/app
mkdir -p internal/{config,database,middleware,models}
mkdir -p internal/modules/{auth,wallet,invoice,payment,refund,admin}
mkdir -p pkg/{response,utils}

# go.mod
cat << 'EOM' > go.mod
module payment-sandbox

go 1.22

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/golang-jwt/jwt/v5 v5.2.0
	gorm.io/driver/postgres v1.5.0
	gorm.io/gorm v1.25.5
)
EOM

# main.go
cat << 'EOM' > cmd/app/main.go
package main

import (
	"payment-sandbox/internal/config"
	"payment-sandbox/internal/database"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.InitDB()
	database.Migrate(db)
	database.Seed(db)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.Run(":8080")
}
EOM

# db config
cat << 'EOM' > internal/config/db.go
package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}
EOM

# migrate
cat << 'EOM' > internal/database/migrate.go
package database

import (
	"payment-sandbox/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&models.User{},
		&models.Wallet{},
		&models.Invoice{},
		&models.PaymentIntent{},
		&models.Refund{},
		&models.BalanceRequest{},
	)
}
EOM

# seed
cat << 'EOM' > internal/database/seed.go
package database

import (
	"payment-sandbox/internal/models"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	db.Create(&models.User{
		Name: "Admin",
		Email: "admin@mail.com",
		Password: "hashed",
		Role: "ADMIN",
	})
}
EOM

# models
cat << 'EOM' > internal/models/user.go
package models

type User struct {
	ID int64 `gorm:"primaryKey"`
	Name string
	Email string
	Password string
	Role string
}
EOM

cat << 'EOM' > internal/models/wallet.go
package models

type Wallet struct {
	ID int64 `gorm:"primaryKey"`
	UserID int64
	Balance float64
}
EOM

cat << 'EOM' > internal/models/invoice.go
package models

type Invoice struct {
	ID int64 `gorm:"primaryKey"`
	MerchantID int64
	InvoiceNumber string
	Amount float64
	Status string
	PaymentLinkToken string
}
EOM

cat << 'EOM' > internal/models/payment.go
package models

type PaymentIntent struct {
	ID int64 `gorm:"primaryKey"`
	InvoiceID int64
	Method string
	Status string
}
EOM

cat << 'EOM' > internal/models/refund.go
package models

type Refund struct {
	ID int64 `gorm:"primaryKey"`
	InvoiceID int64
	Reason string
	Status string
}
EOM

cat << 'EOM' > internal/models/topup.go
package models

type BalanceRequest struct {
	ID int64 `gorm:"primaryKey"`
	MerchantID int64
	Amount float64
	Status string
}
EOM

echo "✅ DONE! Project generated."
