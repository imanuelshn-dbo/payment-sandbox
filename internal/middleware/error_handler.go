package middleware

import (
	"net/http"
	"payment-sandbox/pkg/app-error"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		// custom app error
		if appErr, ok := err.(*apperror.AppError); ok {
			c.JSON(appErr.Code, gin.H{
				"success": false,
				"message": appErr.Message,
				"errors":  appErr.Errors,
			})
			return
		}

		// fallback
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "internal server error",
		})
	}
}
