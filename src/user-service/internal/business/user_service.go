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
	Register(user *models.User) error
	HasPermission(userID uuid.UUID, permissionName string) bool
	GetUser(userID uuid.UUID) (*models.User, error)
	GetListUser(paging *common.Paging) ([]models.User, error)
	ChangePassword(username, password string) error

	CreateProfile(profile *models.Profile) error
	UpdateProfile(profile *models.Profile) error
	GetProfile(idUser string) (*models.Profile, error)
	GetListProfile(paging *common.Paging) ([]models.Profile, error)
}

type UserServiceImpl struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{Repo: repo}
}

// Login handles user authentication
func (s *UserServiceImpl) Login(username, password string) (*models.User, error) {
	if username == "" || password == "" {
		return nil, errors.New("username and password are required")
	}

	user, err := s.Repo.Login(username, password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Register handles new user registration
func (s *UserServiceImpl) Register(user *models.User) error {
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
	return s.Repo.Register(user)
}

// HasPermission checks if user has specific permission
func (s *UserServiceImpl) HasPermission(userID uuid.UUID, permissionName string) bool {
	return s.Repo.HasPermission(userID, permissionName)
}

// GetUserProfile retrieves user profile (without sensitive data)
func (s *UserServiceImpl) GetUser(userID uuid.UUID) (*models.User, error) {
	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserServiceImpl) GetListUser(paging *common.Paging) ([]models.User, error){
	paging.Process()
	return s.Repo.GetListUser(paging)
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

	return s.Repo.ChangePassword(username, password)
}

