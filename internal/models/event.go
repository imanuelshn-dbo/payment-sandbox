package models

import "time"

type EventLog struct {
	Event     string                 `bson:"event"`
	UserID    int64                  `bson:"user_id"`
	RefID     int64                  `bson:"ref_id"`
	Data      map[string]interface{} `bson:"data"`
	CreatedAt time.Time              `bson:"created_at"`
}
