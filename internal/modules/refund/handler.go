package refund

import (
	apperror "payment-sandbox/pkg/app-error"
	"payment-sandbox/pkg/utils"
	"strconv"

	"payment-sandbox/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{s}
}

// @Summary Request Refund
// @Tags Refund
// @Security BearerAuth
// @Accept json
// @Param body body refund.CreateRefundRequest true "Refund"
// @Router /merchant/refund [post]
func (h *Handler) Request(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req CreateRefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.Validation(utils.FormatValidationError(err)))
		return
	}

	err := h.service.RequestRefund(userID, req.InvoiceID)
	if err != nil {
		c.Error(apperror.Validation(utils.FormatValidationError(err)))
		return
	}

	response.Success(c, gin.H{"message": "refund requested"})
}

// @Summary Approve/Reject Refund
// @Tags Refund
// @Security BearerAuth
// @Param id path int true "Refund ID"
// @Param body body refund.UpdateRefundStatusRequest true "Status"
// @Router /admin/refund/{id}/status [put]
func (h *Handler) UpdateStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req UpdateRefundStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.Validation(utils.FormatValidationError(err)))
		return
	}

	err := h.service.UpdateStatus(int64(id), req.Status)
	if err != nil {
		c.Error(apperror.Validation(utils.FormatValidationError(err)))
		return
	}

	response.Success(c, gin.H{"message": "updated"})
}

// @Summary Process Refund
// @Tags Refund
// @Security BearerAuth
// @Param id path int true "Refund ID"
// @Param body body refund.UpdateRefundStatusRequest true "Process"
// @Router /admin/refund/{id}/process [put]
func (h *Handler) Process(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req UpdateRefundStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.Validation(utils.FormatValidationError(err)))
		return
	}

	err := h.service.ProcessRefund(int64(id), req.Status)
	if err != nil {
		c.Error(apperror.Validation(utils.FormatValidationError(err)))
		return
	}

	response.Success(c, gin.H{"message": "processed"})
}
