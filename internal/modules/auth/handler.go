package auth

import (
	apperror "payment-sandbox/pkg/app-error"
	"payment-sandbox/pkg/response"
	"payment-sandbox/pkg/utils"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{s}
}

// @Summary Register Merchant
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body auth.RegisterRequest true "Register"
// @Success 200 {object} map[string]interface{}
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.Validation(utils.FormatValidationError(err)))
		return
	}

	if err := h.service.Register(req); err != nil {
		c.Error(apperror.Validation(utils.FormatValidationError(err)))
		return
	}

	response.Success(c, gin.H{"message": "registered"})
}

// @Summary Login
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body auth.LoginRequest true "Login"
// @Success 200 {object} map[string]interface{}
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.Validation(utils.FormatValidationError(err)))
		return
	}

	token, err := h.service.Login(req)
	if err != nil {
		c.Error(apperror.Validation(utils.FormatValidationError(err)))
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    token,
	})
}

// REFRESH TOKEN
func (h *Handler) Refresh(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(apperror.Validation(utils.FormatValidationError(err)))
		return
	}

	token, err := h.service.RefreshToken(body.RefreshToken)
	if err != nil {
		c.Error(apperror.Validation(utils.FormatValidationError(err)))
		return
	}

	response.Success(c, token)
}
