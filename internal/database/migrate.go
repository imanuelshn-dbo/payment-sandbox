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
