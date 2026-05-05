package refund

import (
	"errors"
	"gorm.io/gorm/clause"
	"payment-sandbox/internal/logger"
	apperror "payment-sandbox/pkg/app-error"

	"payment-sandbox/internal/models"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db}
}

// MERCHANT REQUEST REFUND
func (s *Service) RequestRefund(userID int64, invoiceID int64) error {

	var inv models.Invoice
	if err := s.db.First(&inv, invoiceID).Error; err != nil {
		err = apperror.BadRequest(err.Error())
		return err
	}

	if inv.Status != "PAID" {
		return errors.New("invoice not paid")
	}

	refund := models.Refund{
		InvoiceID: invoiceID,
		UserID:    userID,
		Amount:    inv.Amount,
		Status:    "REQUESTED",
	}

	// log refund request
	logger.LogEvent(
		"REFUND_REQUEST",
		refund.UserID,
		refund.ID,
		map[string]interface{}{
			"amount": refund.Amount,
		},
	)

	return s.db.Create(&refund).Error
}

// ADMIN UPDATE (APPROVE / REJECT)
func (s *Service) UpdateStatus(id int64, status string) error {

	var refund models.Refund
	if err := s.db.First(&refund, id).Error; err != nil {
		err = apperror.BadRequest(err.Error())
		return err
	}

	if refund.Status != "REQUESTED" {
		return errors.New("invalid state")
	}

	if status != "APPROVED" && status != "REJECTED" {
		return errors.New("invalid status")
	}

	// log refund success
	logger.LogEvent(
		"REFUND_SUCCESS",
		refund.UserID,
		refund.ID,
		map[string]interface{}{
			"amount": refund.Amount,
		},
	)

	refund.Status = status
	return s.db.Save(&refund).Error
}

// ADMIN PROCESS (SUCCESS / FAILED)
func (s *Service) ProcessRefund(id int64, status string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {

		// 🔒 LOCK REFUND
		var refund models.Refund
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&refund, id).Error; err != nil {
			err = apperror.BadRequest(err.Error())
			return err
		}

		// ✅ STATE MACHINE
		switch refund.Status {
		case "APPROVED":
			if status != "SUCCESS" && status != "FAILED" {
				return errors.New("invalid transition")
			}
		default:
			return errors.New("invalid state")
		}

		refund.Status = status
		if err := tx.Save(&refund).Error; err != nil {
			err = apperror.BadRequest(err.Error())
			return err
		}

		// SUCCESS FLOW
		if status == "SUCCESS" {

			// 🔒 LOCK WALLET
			var wallet models.Wallet
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("user_id = ?", refund.UserID).
				First(&wallet).Error; err != nil {
				err = apperror.BadRequest(err.Error())
				return err
			}

			if wallet.Balance < refund.Amount {
				return errors.New("insufficient balance")
			}

			wallet.Balance -= refund.Amount
			if err := tx.Save(&wallet).Error; err != nil {
				err = apperror.BadRequest(err.Error())
				return err
			}
		}

		return nil
	})
}
