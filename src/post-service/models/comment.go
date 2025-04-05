package models

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reply struct {
	UserID    string    `json:"user_id" bson:"user_id"`
	Content   string    `json:"content" bson:"content"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type Comment struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	PostID    primitive.ObjectID `json:"post_id" bson:"post_id"`
	UserID    uuid.UUID          `json:"user_id" bson:"user_id"`
	Content   string             `json:"content" bson:"content"`
	Replies   []Reply            `json:"replies" bson:"replies"`
	Status    string             `json:"status" bson:"status"` // "pending", "approved", "rejected"
	LikeCount int                `json:"like_count" bson:"like_count"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}