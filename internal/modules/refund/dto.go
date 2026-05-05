package refund

type CreateRefundRequest struct {
	InvoiceID int64 `json:"invoice_id" binding:"required"`
}

type UpdateRefundStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
