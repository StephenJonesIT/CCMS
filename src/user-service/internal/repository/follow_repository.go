package repository

import (
	"context"

	"github.com/StephenJonesIT/CCMS/src/user-service/internal/models"
	"gorm.io/gorm"
)

type FollowRepository interface {
	Follow(context context.Context, follow *models.Follows) error
	Unfollow(context context.Context, follow *models.Follows) error
}

type FollowRepositoryImpl struct {
	DB *gorm.DB
}

func NewFollowRepository(db *gorm.DB) *FollowRepositoryImpl {
	return &FollowRepositoryImpl{DB: db}
}

func (repo *FollowRepositoryImpl) Follow(context context.Context, follow *models.Follows) error {
	return repo.DB.WithContext(context).Table(models.Follows{}.TableName()).Create(&follow).Error
}

func (repo *FollowRepositoryImpl) Unfollow(context context.Context, follow *models.Follows) error {
	return repo.DB.WithContext(context).Table(models.Follows{}.TableName()).
		Where("follower_id = ? AND followee_id = ?", follow.FollowerID, follow.FolloweeID).
		Delete(models.Follows{}).Error
}