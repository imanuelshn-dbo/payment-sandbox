package payment

import (
	"gorm.io/gorm"
	"payment-sandbox/internal/models"
)

func createLedger(tx *gorm.DB, userID int64, ref string, amount int64, typ string, balance int64) error {
	ledger := models.Ledger{
		UserID:    userID,
		Reference: ref,
		Type:      typ,
		Amount:    amount,
		Balance:   balance,
	}

	return tx.Create(&ledger).Error
}
