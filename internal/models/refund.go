package models

import "time"

type Refund struct {
	ID        int64 `gorm:"primaryKey"`
	InvoiceID int64
	UserID    int64
	Amount    int64
	Status    string // REQUESTED, APPROVED, REJECTED, SUCCESS, FAILED
	CreatedAt time.Time
}
