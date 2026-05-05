package models

import "time"

type Admin struct {
	ID        int64 `gorm:"primaryKey"`
	Name      string
	Email     string `gorm:"unique"`
	MobileNo  string
	NIK       string
	CreatedAt time.Time
	UpdatedAt time.Time
}
