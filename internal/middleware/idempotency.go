package middleware

import (
	"net/http"
	"payment-sandbox/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Idempotency(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		key := c.GetHeader("Idempotency-Key")
		if key == "" {
			c.Next()
			return
		}

		var record models.IdempotencyKey
		err := db.Where("key = ?", key).First(&record).Error

		if err == nil {
			// sudah pernah → return cached response
			c.Data(http.StatusOK, "application/json", []byte(record.Response))
			c.Abort()
			return
		}

		// capture response
		writer := &responseWriter{body: []byte{}, ResponseWriter: c.Writer}
		c.Writer = writer

		c.Next()

		// save response
		db.Create(&models.IdempotencyKey{
			Key:      key,
			Endpoint: c.FullPath(),
			Response: string(writer.body),
		})
	}
}

type responseWriter struct {
	gin.ResponseWriter
	body []byte
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return w.ResponseWriter.Write(b)
}
