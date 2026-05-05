package admin

import (
	"time"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db}
}

type DashboardResponse struct {
	TotalInvoice     int64
	TotalPaid        int64
	TotalFailed      int64
	TotalExpired     int64
	TotalTransaction int64
	TotalRefund      int64
}

func (s *Service) GetDashboard(filter DashboardFilter) (*DashboardResponse, error) {

	var res DashboardResponse

	query := s.db.Table("invoices")

	// FILTER USER
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}

	// FILTER DATE
	if filter.StartDate != "" && filter.EndDate != "" {
		start, _ := time.Parse("2006-01-02", filter.StartDate)
		end, _ := time.Parse("2006-01-02", filter.EndDate)

		query = query.Where("created_at BETWEEN ? AND ?", start, end)
	}

	// TOTAL INVOICE
	query.Count(&res.TotalInvoice)

	// STATUS COUNT
	s.db.Table("invoices").Where("status = 'PAID'").Count(&res.TotalPaid)
	s.db.Table("invoices").Where("status = 'FAILED'").Count(&res.TotalFailed)
	s.db.Table("invoices").Where("status = 'EXPIRED'").Count(&res.TotalExpired)

	// TOTAL TRANSACTION (PAID)
	s.db.Table("invoices").
		Select("COALESCE(SUM(amount),0)").
		Where("status = 'PAID'").
		Scan(&res.TotalTransaction)

	// TOTAL REFUND (SUCCESS)
	s.db.Table("refunds").
		Select("COALESCE(SUM(amount),0)").
		Where("status = 'SUCCESS'").
		Scan(&res.TotalRefund)

	return &res, nil
}
