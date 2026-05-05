package request

type CreateInvoiceRequest struct {
	Amount int64 `json:"amount" binding:"required,gt=0"`
}

type UpdateInvoiceRequest struct {
	Amount int64 `json:"amount" binding:"required,gt=0"`
}
