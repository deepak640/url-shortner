package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type URL struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	ShortCode string        `bson:"short_code"`
	UserID		string        `bson:"user_id"`
	LongURL   string        `bson:"long_url"`
	CreatedAt time.Time     `bson:"created_at"`
	ExpiresAt    *time.Time         `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
  MaxClicks    int                `bson:"max_clicks" json:"max_clicks"`
  CurrentClicks int               `bson:"current_clicks" json:"current_clicks"`
  IsActive     bool               `bson:"is_active" json:"is_active"`
}
