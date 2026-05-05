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
