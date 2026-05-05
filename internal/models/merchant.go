package models

import "time"

type Merchant struct {
	ID             int64   `gorm:"primaryKey"`
	MerchantName   string  `gorm:"size:100;not null"`
	OwnerName      string  `gorm:"size:100"`
	Category       string  `gorm:"size:50"`
	Address        string  `gorm:"type:text"`
	PhoneNumber    string  `gorm:"size:20"`
	Email          string  `gorm:"size:100;unique"`
	Status         string  `gorm:"size:20"`
	CommissionRate float64 `gorm:"type:decimal(5,2)"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
