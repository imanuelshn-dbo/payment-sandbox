package admin

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

// @Summary Admin Dashboard
// @Tags Admin
// @Security BearerAuth
// @Produce json
// @Param user_id query int false "User ID"
// @Param start_date query string false "Start Date"
// @Param end_date query string false "End Date"
// @Router /admin/dashboard [get]
func (h *Handler) Dashboard(c *gin.Context) {

	var filter DashboardFilter

	if userID := c.Query("user_id"); userID != "" {
		id, _ := strconv.Atoi(userID)
		uid := int64(id)
		filter.UserID = &uid
	}

	filter.StartDate = c.Query("start_date")
	filter.EndDate = c.Query("end_date")

	data, err := h.service.GetDashboard(filter)
	if err != nil {
		c.Error(apperror.Validation(utils.FormatValidationError(err)))
		return
	}

	response.Success(c, data)
}
