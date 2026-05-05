package models

import "time"

type Invoice struct {
	ID            int64 `gorm:"primaryKey"`
	UserID        int64
	InvoiceNumber string `gorm:"uniqueIndex"`
	Amount        int64
	Status        string // PENDING, PAID, FAILED
	PaymentToken  string `gorm:"uniqueIndex"`
	ExpiredAt     time.Time
	CreatedAt     time.Time
}

type InvoiceFilter struct {
	Status string `form:"status"`
}
