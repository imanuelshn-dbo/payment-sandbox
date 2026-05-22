package payment

import (
	"errors"
	"gorm.io/gorm/clause"
	"payment-sandbox/internal/logger"
	apperror "payment-sandbox/pkg/app-error"
	"time"

	"payment-sandbox/internal/models"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db}
}

// GET INVOICE BY TOKEN (PUBLIC)
func (s *Service) GetInvoiceByToken(token string) (*models.Invoice, error) {
	var inv models.Invoice
	err := s.db.Where("payment_token = ?", token).First(&inv).Error
	return &inv, err
}

// CREATE PAYMENT INTENT
func (s *Service) CreatePayment(invoiceID int64, method string) error {
	var invoice models.Invoice
	if err := s.db.First(&invoice, invoiceID).Error; err != nil {
		return apperror.BadRequest("invoice not found")
	}

	// CHECK INVOICE STATUS
	if invoice.Status != "PENDING" {
		return apperror.BadRequest("invoice already processed")
	}

	// RETRY PAYMENT
	if invoice.Status == "PAID" {
		return apperror.BadRequest("invoice already paid")
	}

	// CHECK EXPIRED TIME
	if time.Now().After(invoice.ExpiredAt) {
		return apperror.BadRequest("invoice expired")
	}

	intent := models.PaymentIntent{
		InvoiceID:     invoiceID,
		PaymentMethod: method,
		Status:        "PENDING",
	}

	return s.db.Create(&intent).Error
}

// ADMIN UPDATE PAYMENT
func (s *Service) UpdatePayment(id int64, status string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {

		// LOCK PAYMENT INTENT
		var intent models.PaymentIntent
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&intent, id).Error; err != nil {
			err = apperror.BadRequest(err.Error())
			return err
		}

		// STATUS CHECK
		if intent.Status != "PENDING" {
			return errors.New("payment already processed")
		}

		// LOCK INVOICE
		var invoice models.Invoice
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&invoice, intent.InvoiceID).Error; err != nil {
			return err
		}

		// EXPIRED CHECK
		if time.Now().After(invoice.ExpiredAt) {
			intent.Status = "EXPIRED"

			// log payment expired
			logger.LogEvent(
				"PAYMENT_EXPIRED",
				invoice.UserID,
				invoice.ID,
				map[string]interface{}{
					"amount": invoice.Amount,
				},
			)

			return tx.Save(&intent).Error
		}

		// UPDATE STATUS
		intent.Status = status
		now := time.Now()

		if status == "SUCCESS" {
			intent.PaidAt = &now
		}

		if err := tx.Save(&intent).Error; err != nil {
			err = apperror.BadRequest(err.Error())
			return err
		}

		// log payment failed
		if status == "FAILED" {
			logger.LogEvent(
				"PAYMENT_FAILED",
				invoice.UserID,
				invoice.ID,
				map[string]interface{}{
					"amount": invoice.Amount,
					"method": intent.PaymentMethod,
				},
			)
		}

		// SUCCESS FLOW
		if status == "SUCCESS" {

			// UPDATE INVOICE
			if invoice.Status != "PENDING" {
				return errors.New("invoice already processed")
			}

			invoice.Status = "PAID"
			if err := tx.Save(&invoice).Error; err != nil {
				err = apperror.BadRequest(err.Error())
				return err
			}

			// WALLET DEDUCTION
			if intent.PaymentMethod == "WALLET" {

				var wallet models.Wallet
				if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
					Where("user_id = ?", invoice.UserID).
					First(&wallet).Error; err != nil {
					err = apperror.BadRequest(err.Error())
					return err
				}

				if wallet.Balance < invoice.Amount {
					// masuk error ini
					if wallet.Balance < invoice.Amount {
						return apperror.BadRequestWithError(
							"insufficient balance",
							map[string]interface{}{
								"wallet_balance": wallet.Balance,
								"invoice_amount": invoice.Amount,
							},
						)
					}
				}

				wallet.Balance -= invoice.Amount
				if err := tx.Save(&wallet).Error; err != nil {
					err = apperror.BadRequest(err.Error())
					return err
				}
			}

			// logging event payment success
			logger.LogEvent(
				"PAYMENT_SUCCESS",
				invoice.UserID,
				invoice.ID,
				map[string]interface{}{
					"amount": invoice.Amount,
					"method": intent.PaymentMethod,
				},
			)
		}

		return nil
	})
}
