package models

type User struct {
	ID         int64  `gorm:"primaryKey"`
	Email      string `gorm:"unique"`
	Password   string
	Role       string
	MerchantID int64
	AdminID    int64
}
