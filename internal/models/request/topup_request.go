package request

type TopUpRequest struct {
	Amount int64 `json:"amount" binding:"required,gt=0"`
}

type UpdateTopUpStatusRequest struct {
	Status string `json:"status" binding:"required"` // SUCCESS / FAILED
}
