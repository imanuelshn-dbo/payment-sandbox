package models

import "time"

type Ledger struct {
	ID        int64 `gorm:"primaryKey"`
	UserID    int64
	Reference string // invoice / payment id
	Type      string // DEBIT / CREDIT
	Amount    int64
	Balance   int64 // balance after transaction
	CreatedAt time.Time
}
