/*
 * @File: bussiness.profile_service.go
 * @Description: Implements Profile business functions
 * @Author: Tran Thanh Sang (tranthanhsang.it.la@gmail.com)
 */
package business

import (
	"errors"
	"fmt"
	"strings"

	"github.com/StephenJonesIT/CCMS/src/user-service/common"
	"github.com/StephenJonesIT/CCMS/src/user-service/internal/models"
	"github.com/StephenJonesIT/CCMS/src/user-service/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProfileService interface {
	CreateProfile(profile *models.Profile) error
	UpdateProfile(profile *models.Profile) error
	GetProfile(idUser string) (*models.Profile, error)
	GetListProfile(paging *common.Paging) ([]models.Profile, error)
}

type ProfileServiceImpl struct {
	Repo repository.ProfileRepository
}

func NewProfileService(repo repository.ProfileRepository) *ProfileServiceImpl {
	return &ProfileServiceImpl{Repo: repo}
}

func (s *ProfileServiceImpl) CreateProfile(profile *models.Profile) error {
	if profile.FullName == "" {
		return errors.New("fullname is required")
	}

	if profile.UserID == uuid.Nil {
        return errors.New("user ID is required")
    }

	existing, err := s.Repo.GetProfile(profile.UserID.String())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("error checking existing profile: %w", err)
	}

	if existing != nil {
		return errors.New("profile already exists for this user")
	}

	if err := s.Repo.CreateProfile(profile); err != nil{
		return fmt.Errorf("failed to create profile: %w", err)
	}

	return nil
}

func (s *ProfileServiceImpl) UpdateProfile(profile *models.Profile) error {
	if profile == nil {
        return errors.New("profile cannot be nil")
    }
    
    if profile.UserID == uuid.Nil {
        return errors.New("user ID is required")
    }

	if profile.ProfileID == 0 {
		return errors.New("profile id is required")
	}

  	// Check if profile exists
    if _, err := s.Repo.GetProfileUpdate(profile.UserID.String(), profile.ProfileID); err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return fmt.Errorf("profile not found for user %s: %w", profile.UserID, err)
        }
        return fmt.Errorf("error checking existing profile: %w", err)
    }

	// Update profile in repository
	if  err := s.Repo.UpdateProfile(profile); err != nil{
		return fmt.Errorf("failed to update profile: %w",err)
	}

	return nil
}

func (s *ProfileServiceImpl) GetProfile(idUSer string) (*models.Profile, error) {
	if strings.TrimSpace(idUSer) == "" {
		return nil, fmt.Errorf("invalid input for %s: %s", "userID", "cannot be empty")
	}

	userUUID, err := uuid.Parse(idUSer);
	if  err != nil  {
		return nil, fmt.Errorf("invalid input for %s: %s", "userID", "invalid UUID format")
	}

	// Get profile
    profile, err := s.Repo.GetProfile(userUUID.String())
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("profile not found for user %v", userUUID)
        }
        return nil, fmt.Errorf("repository error: %w", err)
    }

    if profile == nil {
        return nil, errors.New("data integrity error: nil profile returned")
    }

    return profile, nil
}

func(s *ProfileServiceImpl) GetListProfile(paging *common.Paging) ([]models.Profile, error) {
	paging.Process()
	return s.Repo.GetListProfile(paging)
}
