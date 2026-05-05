package invoice

import (
	"payment-sandbox/internal/models/request"
	apperror "payment-sandbox/pkg/app-error"
	"payment-sandbox/pkg/pagination"
	"payment-sandbox/pkg/response"
	"payment-sandbox/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{s}
}

// @Summary Create Invoice
// @Tags Invoice
// @Security BearerAuth
// @Accept json
// @Param body body request.CreateInvoiceRequest true "Create Invoice"
// @Success 200 {object} models.Invoice
// @Router /merchant/invoice [post]
func (h *Handler) Create(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req request.CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.Validation(utils.FormatValidationError(err)))
		return
	}

	invoice, err := h.service.CreateInvoice(userID, req.Amount)
	if err != nil {
		c.Error(apperror.Validation(utils.FormatValidationError(err)))
		return
	}

	response.Success(c, invoice)
}

// @Summary List Invoice
// @Tags Invoice
// @Security BearerAuth
// @Produce json
// @Param status query string false "Filter status"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Router /merchant/invoice [get]
func (h *Handler) List(c *gin.Context) {
	userID := c.GetInt64("user_id")
	status := c.Query("status")

	p := pagination.GetPagination(c)

	data, total, err := h.service.ListInvoice(userID, status, p.Page, p.Limit)
	if err != nil {
		c.Error(apperror.Validation(utils.FormatValidationError(err)))
		return
	}

	response.Success(c, gin.H{
		"data":  data,
		"page":  p.Page,
		"limit": p.Limit,
		"total": total,
	})
}

// @Summary Update Invoice
// @Tags Invoice
// @Security BearerAuth
// @Accept json
// @Param id path int true "Invoice ID"
// @Param body body request.UpdateInvoiceRequest true "Update"
// @Router /merchant/invoice/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	userID := c.GetInt64("user_id")

	id, _ := strconv.Atoi(c.Param("id"))

	var req request.UpdateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.Validation(utils.FormatValidationError(err)))
		return
	}

	invoice, err := h.service.UpdateInvoice(userID, int64(id), req.Amount)
	if err != nil {
		c.Error(apperror.Validation(utils.FormatValidationError(err)))
		return
	}

	response.Success(c, invoice)
}

// @Summary Delete Invoice
// @Tags Invoice
// @Security BearerAuth
// @Param id path int true "Invoice ID"
// @Router /merchant/invoice/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	userID := c.GetInt64("user_id")

	id, _ := strconv.Atoi(c.Param("id"))

	err := h.service.DeleteInvoice(userID, int64(id))
	if err != nil {
		c.Error(apperror.Validation(utils.FormatValidationError(err)))
		return
	}

	response.Success(c, gin.H{"message": "invoice deleted"})
}
