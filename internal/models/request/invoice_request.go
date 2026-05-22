package request

import "time"

type CreateInvoiceRequest struct {
	Amount    int64     `json:"amount" binding:"required,gt=0"`
	ExpiredAt time.Time `json:"expired_at"`
}

type UpdateInvoiceRequest struct {
	Amount int64 `json:"amount" binding:"required,gt=0"`
}
