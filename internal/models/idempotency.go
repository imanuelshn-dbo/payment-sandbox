package models

import "time"

type IdempotencyKey struct {
	ID        int64     `gorm:"primaryKey"`
	Key       string    `gorm:"uniqueIndex"`
	Endpoint  string
	Response  string    `gorm:"type:text"`
	CreatedAt time.Time
}