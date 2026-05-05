package wallet

import (
	"errors"
	"payment-sandbox/internal/models"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db}
}

func (s *Service) GetBalance(userID int64) (int64, error) {
	var wallet models.Wallet

	err := s.db.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		wallet = models.Wallet{
			UserID:  userID,
			Balance: 0,
		}
		s.db.Create(&wallet)
	}

	return wallet.Balance, nil
}

func (s *Service) TopUp(userID int64, amount int64) error {

	topup := models.TopUp{
		MerchantID: userID,
		Amount:     amount,
		Status:     "PENDING",
	}

	return s.db.Create(&topup).Error
}

func (s *Service) UpdateTopUpStatus(id int64, adminID int64, status string) error {

	return s.db.Transaction(func(tx *gorm.DB) error {

		var topup models.TopUp
		if err := tx.First(&topup, id).Error; err != nil {
			return err
		}

		// validate status, only process when status pending
		if topup.Status != "PENDING" {
			return errors.New("already processed")
		}

		now := time.Now()

		topup.Status = status
		topup.AdminID = &adminID
		topup.ApprovalDate = &now

		if err := tx.Save(&topup).Error; err != nil {
			return err
		}

		// topup balance
		if status == "SUCCESS" {

			var wallet models.Wallet
			if err := tx.Where("user_id = ?", topup.MerchantID).
				First(&wallet).Error; err != nil {
				return err
			}

			wallet.Balance += topup.Amount

			if err := tx.Save(&wallet).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
