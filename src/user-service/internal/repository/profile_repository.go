/*
 * @File: repository.profile_repository.go
 * @Description: Implements Profile CRUD functions for Posgres
 * @Author: Tran Thanh Sang (tranthanhsang.it.la@gmail.com)
 */
package repository

import (
	"fmt"
	"user-service/common"
	"user-service/internal/models"
)

func (repo *UserRepoImpl) UpdateProfile(profile *models.Profile) error {
	return repo.DB.Table(models.Profile{}.TableName()).Where("profile_id = ?",profile.ProfileID).Updates(profile).Error
}

func (repo *UserRepoImpl) CreateProfile(profile *models.Profile) error {
	return repo.DB.Create(profile).Error
}

func (repo *UserRepoImpl) GetProfile(idUser string) (*models.Profile, error) {
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

func (repo *UserRepoImpl)GetListProfile(paging *common.Paging) ([]models.Profile, error) {
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