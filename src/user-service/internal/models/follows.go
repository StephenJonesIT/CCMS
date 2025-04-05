package models

import (
	"time"

	"github.com/google/uuid"
)

type Follows struct {
	FollowerID uuid.UUID `gorm:"column:follower_id" json:"follower_id,omitempty"`
	FolloweeID uuid.UUID `gorm:"column:followee_id" json:"followee_id,omitempty"`
	CreatedAt  *time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
}

func (Follows) TableName() string {
	return "follows"
}