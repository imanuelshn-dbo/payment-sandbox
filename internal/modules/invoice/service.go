package invoice

import (
	"fmt"
	"payment-sandbox/internal/logger"
	apperror "payment-sandbox/pkg/app-error"
	"time"

	"payment-sandbox/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db}
}

func (s *Service) CreateInvoice(userID int64, amount int64) (*models.Invoice, error) {
	invoice := models.Invoice{
		UserID:        userID,
		InvoiceNumber: generateInvoiceNumber(),
		Amount:        amount,
		Status:        "PENDING",
		PaymentToken:  uuid.New().String(),
		ExpiredAt:     time.Now().Add(24 * time.Hour),
	}

	err := s.db.Create(&invoice).Error
	if err != nil {
		err = apperror.BadRequest(err.Error())
		return nil, err
	}

	// log event create invoice
	logger.LogEvent(
		"CREATE_INVOICE",
		userID,
		invoice.ID,
		map[string]interface{}{
			"amount": amount,
		},
	)

	return &invoice, nil
}

func (s *Service) ListInvoice(userID int64, status string, page int, limit int) ([]models.Invoice, int64, error) {
	var invoices []models.Invoice
	var total int64

	// get data from DB
	query := s.db.Model(&models.Invoice{}).Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// count total
	if err := query.Count(&total).Error; err != nil {
		err = apperror.BadRequest(err.Error())
		return nil, 0, err
	}

	// pagination
	err := query.
		Limit(limit).
		Offset((page - 1) * limit).
		Order("created_at desc").
		Find(&invoices).Error

	return invoices, total, err
}

func (s *Service) UpdateInvoice(userID int64, invoiceID int64, amount int64) (*models.Invoice, error) {
	var invoice models.Invoice

	err := s.db.Where("id = ? AND user_id = ?", invoiceID, userID).First(&invoice).Error
	if err != nil {
		err = apperror.BadRequest(err.Error())
		return nil, err
	}

	if invoice.Status != "PENDING" {
		return nil, fmt.Errorf("cannot update non-pending invoice")
	}

	invoice.Amount = amount

	err = s.db.Save(&invoice).Error
	if err != nil {
		err = apperror.BadRequest(err.Error())
		return nil, err
	}

	return &invoice, nil
}

func (s *Service) DeleteInvoice(userID int64, invoiceID int64) error {
	var invoice models.Invoice

	err := s.db.Where("id = ? AND user_id = ?", invoiceID, userID).First(&invoice).Error
	if err != nil {
		err = apperror.BadRequest(err.Error())
		return err
	}

	if invoice.Status != "PENDING" {
		return fmt.Errorf("cannot delete non-pending invoice")
	}

	return s.db.Delete(&invoice).Error
}

// Auto generate invoice number
func generateInvoiceNumber() string {
	return fmt.Sprintf(
		"INV-%d-%s",
		time.Now().Unix(),
		uuid.New().String()[:8],
	)
}
