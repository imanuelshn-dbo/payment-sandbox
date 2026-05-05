package logger

import (
	"context"
	"time"

	"payment-sandbox/internal/config"
	"payment-sandbox/internal/models"
)

func LogEvent(event string, userID int64, refID int64, data map[string]interface{}) {

	collection := config.MongoDB.Collection("events")

	logData := models.EventLog{
		Event:     event,
		UserID:    userID,
		RefID:     refID,
		Data:      data,
		CreatedAt: time.Now(),
	}

	// EL : Async flow to maximize performance
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, _ = collection.InsertOne(ctx, logData)
	}()
}
