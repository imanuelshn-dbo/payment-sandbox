package models

import "time"

type Wallet struct {
	ID      int64 `gorm:"primaryKey"`
	UserID  int64 `gorm:"uniqueIndex"`
	Balance int64
}

type TopUp struct {
	ID           int64  `gorm:"primaryKey"`
	MerchantID   int64  `gorm:"index"`
	AdminID      *int64 `gorm:"index"`
	Amount       int64
	Status       string `gorm:"size:20"`
	ApprovalDate *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `gorm:"index"`
}
