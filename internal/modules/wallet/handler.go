package wallet

import (
	"payment-sandbox/internal/models/request"
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

// @Summary Get Wallet Balance
// @Tags Wallet
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.Wallet
// @Router /merchant/wallet [get]
func (h *Handler) GetBalance(c *gin.Context) {
	userID := c.GetInt64("user_id")

	balance, _ := h.service.GetBalance(userID)

	response.Success(c, gin.H{
		"balance": balance,
	})
}

// @Summary Top Up Wallet
// @Tags Wallet
// @Security BearerAuth
// @Accept json
// @Param body body request.TopUpRequest true "TopUp"
// @Router /merchant/topup [post]
func (h *Handler) TopUp(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req request.TopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.Validation(utils.FormatValidationError(err)))
		return
	}

	err := h.service.TopUp(userID, req.Amount)
	if err != nil {
		c.Error(apperror.Validation(utils.FormatValidationError(err)))
		return
	}

	response.Success(c, gin.H{
		"message": "topup success submitted, wait for admin aprrove",
		"status":  "PENDING",
	})
}

// @Summary Approval Topup Request
// @Tags Wallet
// @Security BearerAuth
// @Accept json
// @Param body body request.UpdateTopUpStatusRequest true "TopUpApproval"
// @Router /admin/topup/{id} [post]
func (h *Handler) UpdateTopUpStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	// Get id admin from token
	adminID := int64(c.GetFloat64("user_id"))

	err := h.service.UpdateTopUpStatus(int64(id), adminID, req.Status)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "updated"})
}
