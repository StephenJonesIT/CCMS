/*
 * @File: bussiness.user_service.go
 * @Description: Implements User business functions
 * @Author: Tran Thanh Sang (tranthanhsang.it.la@gmail.com)
 */
package business

import (
	"errors"

	"github.com/StephenJonesIT/CCMS/src/user-service/common"
	"github.com/StephenJonesIT/CCMS/src/user-service/internal/models"
	"github.com/StephenJonesIT/CCMS/src/user-service/internal/repository"
	"github.com/google/uuid"
)

type UserService interface {
	Login(username, password string) (*models.User, error)
	Register(user *models.UserRegister) error
	HasPermission(userID uuid.UUID, permissionName string) bool
	GetUser(userID uuid.UUID) (*models.User, error)
	GetListUser(paging *common.Paging) ([]models.User, error)
	ChangePassword(username, password string) error
}

type UserServiceImpl struct {
	userRepo repository.UserRepository
	profileRepo repository.ProfileRepository
}

func NewUserService(userRepo repository.UserRepository, profileRepo repository.ProfileRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: userRepo,
		profileRepo: profileRepo,
	}
}

// Login handles user authentication
func (s *UserServiceImpl) Login(username, password string) (*models.User, error) {
	if username == "" || password == "" {
		return nil, errors.New("username and password are required")
	}

	user, err := s.userRepo.Login(username, password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Register handles new user registration
func (s *UserServiceImpl) Register(user *models.UserRegister) error {
	// Validate input
	if user.UserName == "" {
		return errors.New("username is required")
	}
	if user.Password == "" {
		return errors.New("password is required")
	}
	if len(user.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	if user.RoleID == 0 {
		user.RoleID = 1 // Giá trị mặc định
	}

	if err := s.userRepo.Register(user); err != nil {
		return errors.New("failed to register user: " + err.Error())
	}

	if err := s.profileRepo.CreateProfile(&models.Profile{
		UserID: user.UserID,
	}); err != nil {
		return errors.New("failed to create profile: " + err.Error())
	}
	return nil
}

// HasPermission checks if user has specific permission
func (s *UserServiceImpl) HasPermission(userID uuid.UUID, permissionName string) bool {
	return s.userRepo.HasPermission(userID, permissionName)
}

// GetUserProfile retrieves user profile (without sensitive data)
func (s *UserServiceImpl) GetUser(userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserServiceImpl) GetListUser(paging *common.Paging) ([]models.User, error){
	paging.Process()
	return s.userRepo.GetListUser(paging)
}

func (s *UserServiceImpl) ChangePassword(username, password string) error {
	if username == "" {
		return errors.New("username is required")
	}

	if password == "" {
		return errors.New("password is required")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	return s.userRepo.ChangePassword(username, password)
}

