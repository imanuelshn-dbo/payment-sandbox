package config

import (
	"fmt"
	"log"
	"os"
	"payment-sandbox/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	log.Println("HOST:", os.Getenv("DB_HOST"))
	log.Println("PORT:", os.Getenv("DB_PORT"))
	log.Println("USER:", os.Getenv("DB_USER"))
	log.Println("PASS:", os.Getenv("DB_PASS"))
	log.Println("DB:", os.Getenv("DB_NAME"))

	if port == "" {
		port = "5432"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("failed to migrate:", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Wallet{},
		&models.TopUp{},
		&models.IdempotencyKey{},
		&models.Ledger{},
		&models.PaymentIntent{},
		&models.Invoice{},
		&models.Refund{},
		&models.Merchant{},
		&models.Admin{},
	)
	if err != nil {
		log.Fatal("failed to migrate:", err)
	}

	log.Println("database connected")

	return db
}
