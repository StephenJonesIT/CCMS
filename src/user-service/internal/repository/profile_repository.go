/*
 * @File: repository.profile_repository.go
 * @Description: Implements Profile CRUD functions for Posgres
 * @Author: Tran Thanh Sang (tranthanhsang.it.la@gmail.com)
 */
package repository

import (
	"context"
	"fmt"

	"github.com/StephenJonesIT/CCMS/src/user-service/common"
	"github.com/StephenJonesIT/CCMS/src/user-service/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProfileRepository interface {
	GetProfile(idUser string) (*models.Profile, error)
	GetListProfile(paging *common.Paging) ([]models.Profile, error)
	CreateProfile(profile *models.Profile) error
	UpdateProfile(profile *models.Profile) error
	GetProfileUpdate(idUser string, idProfile int64) (*models.Profile, error)

	IncrementFollowersCount(ctx context.Context, userID uuid.UUID) error
	DecrementFollowersCount(ctx context.Context, userID uuid.UUID) error
	IncrementFolloweeCount(ctx context.Context, userID uuid.UUID) error
	DecrementFolloweeCount(ctx context.Context, userID uuid.UUID) error
	GetFollowers(ctx context.Context, userID uuid.UUID) ([]models.Profile, error)
	GetFollowees(ctx context.Context, userID uuid.UUID) ([]models.Profile, error)
}

type ProfileRepositoryImpl struct {
	DB *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *ProfileRepositoryImpl {
	return &ProfileRepositoryImpl{
		DB: db,
	}
}

func (repo *ProfileRepositoryImpl) UpdateProfile(profile *models.Profile) error {
	return repo.DB.Table(models.Profile{}.TableName()).Where("profile_id = ?",profile.ProfileID).Updates(profile).Error
}

func (repo *ProfileRepositoryImpl) CreateProfile(profile *models.Profile) error {
	return repo.DB.Create(profile).Error
}

func (repo *ProfileRepositoryImpl) GetProfile(idUser string) (*models.Profile, error) {
	var profile models.Profile
	err := repo.DB.Table(models.Profile{}.TableName()).
					Where("user_id = ?", idUser).
					First(&profile).
					Error
	
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (repo *ProfileRepositoryImpl)GetListProfile(paging *common.Paging) ([]models.Profile, error) {
	var profiles []models.Profile

	if err := repo.DB.Table(models.Profile{}.TableName()).Count(&paging.Total).Error; err != nil {
		return nil, fmt.Errorf("failed to count profileprofiles: %w", err)
	}

	query := repo.DB.Table(models.Profile{}.TableName()).
		Order("profile_id DESC").
		Offset((paging.Page - 1) * paging.Limit).
		Limit(paging.Limit)

    if err := query.Find(&profiles).Error; err != nil {
        return nil, fmt.Errorf("failed to get user list: %w", err)
    }
	return profiles, nil
}

func (repo *ProfileRepositoryImpl) GetProfileUpdate(idUser string, idProfile int64) (*models.Profile, error){
	var profile models.Profile
	err := repo.DB.Table(models.Profile{}.TableName()).
					Where("user_id = ? AND profile_id = ?", idUser, idProfile).
					First(&profile).
					Error
	
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (repo *ProfileRepositoryImpl) IncrementFollowersCount(ctx context.Context, userID uuid.UUID) error {
	return repo.DB.WithContext(ctx).Table(models.Profile{}.TableName()).
		Where("user_id = ?", userID).
		UpdateColumn("follower_count", gorm.Expr("followers_count + 1")).Error
}

func (repo *ProfileRepositoryImpl) DecrementFollowersCount(ctx context.Context, userID uuid.UUID) error {
	return repo.DB.WithContext(ctx).Table(models.Profile{}.TableName()).
		Where("user_id = ?", userID).
		UpdateColumn("followers_count", gorm.Expr("followers_count - 1")).Error
}

func (repo *ProfileRepositoryImpl) IncrementFolloweeCount(ctx context.Context, userID uuid.UUID) error {
	return repo.DB.WithContext(ctx).Table(models.Profile{}.TableName()).
		Where("user_id = ?", userID).
		UpdateColumn("followee_count", gorm.Expr("followee_count + 1")).Error
}

func (repo *ProfileRepositoryImpl) DecrementFolloweeCount(ctx context.Context, userID uuid.UUID) error {
	return repo.DB.WithContext(ctx).Table(models.Profile{}.TableName()).
		Where("user_id = ?", userID).
		UpdateColumn("followee_count", gorm.Expr("followee_count - 1")).Error
}

func (repo *ProfileRepositoryImpl) GetFollowers(ctx context.Context, userID uuid.UUID) ([]models.Profile, error) {
	var followers []models.Profile
	err := repo.DB.WithContext(ctx).Table(models.Profile{}.TableName()).
		Select("user_id, fullname, email, bio, picture_url, interests, follower_count, followee_count").
		Joins("JOIN follows ON user_profiles.user_id = follows.follower_id").
		Where("follows.followee_id = ?", userID).
		Find(&followers).Error
	if err != nil {
		return nil, err
	}
	return followers, nil
}

func (repo *ProfileRepositoryImpl) GetFollowees(ctx context.Context, userID uuid.UUID) ([]models.Profile, error) {
	var followees []models.Profile
	err := repo.DB.WithContext(ctx).Table(models.Profile{}.TableName()).
		Select("user_id, fullname, email, bio, picture_url, interests, follower_count, followee_count").
		Joins("JOIN follows ON user_profiles.user_id = follows.followee_id").
		Where("follows.follower_id = ?", userID).
		Find(&followees).Error
	if err != nil {
		return nil, err
	}
	return followees, nil
}