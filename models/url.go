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
}
