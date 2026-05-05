package models

type BalanceRequest struct {
	ID         int64 `gorm:"primaryKey"`
	MerchantID int64
	Amount     float64
	Status     string
}
