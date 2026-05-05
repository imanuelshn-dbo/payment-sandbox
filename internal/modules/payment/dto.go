package payment

type CreatePaymentRequest struct {
	PaymentMethod string `json:"payment_method" binding:"required"`
}

type UpdatePaymentStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
