package models

import (
	"time"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Media struct {
	URL  string `json:"url" bson:"url"`
	Type string `json:"type" bson:"type"` // "image" or "video"
}

type Post struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title        *string            `json:"title" bson:"title"`         // Nullable
	Content      *string            `json:"content" bson:"content"`     // Nullable
	AuthorID     uuid.UUID          `json:"author_id" bson:"author_id"` // Non-nullable
	Categories   []string           `json:"categories" bson:"categories"`
	Tags         []string           `json:"tags" bson:"tags"`
	Media        []Media            `json:"media" bson:"media"`
	Status       string             `json:"status" bson:"status"` // Non-nullable: "draft", "pending", "published"
	LikeCount    int                `json:"like_count" bson:"like_count"`
	CommentCount int                `json:"comment_count" bson:"comment_count"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}
