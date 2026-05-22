package payment

import (
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

// @Summary Get Public Payment Page
// @Tags Payment
// @Produce json
// @Param token path string true "Payment Token"
// @Router /pay/{token} [get]
func (h *Handler) GetInvoice(c *gin.Context) {
	token := c.Param("token")

	inv, err := h.service.GetInvoiceByToken(token)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, inv)
}

// @Summary Create Payment
// @Tags Payment
// @Accept json
// @Param token path string true "Token"
// @Param body body payment.CreatePaymentRequest true "Payment"
// @Router /pay/{token} [post]
func (h *Handler) Pay(c *gin.Context) {
	token := c.Param("token")

	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	inv, err := h.service.GetInvoiceByToken(token)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.service.CreatePayment(inv.ID, req.PaymentMethod)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "payment created"})
}

// @Summary Update Payment (Admin)
// @Tags Payment
// @Security BearerAuth
// @Param id path int true "Payment ID"
// @Param body body payment.UpdatePaymentStatusRequest true "Status"
// @Router /admin/payment/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req UpdatePaymentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	err := h.service.UpdatePayment(int64(id), req.Status)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, gin.H{"message": "updated"})
}
