package models

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Like struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	UserID    uuid.UUID  		 `json:"user_id" bson:"user_id"`
	TargetID  primitive.ObjectID `json:"target_id" bson:"target_id"`
	TargetType string            `json:"target_type" bson:"target_type"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}