package models

import "time"

type PaymentIntent struct {
	ID            int64 `gorm:"primaryKey"`
	InvoiceID     int64
	PaymentMethod string // WALLET, VA_DUMMY, EWALLET_DUMMY
	Status        string // PENDING, SUCCESS, FAILED, EXPIRED
	PaidAt        *time.Time
	CreatedAt     time.Time
}
