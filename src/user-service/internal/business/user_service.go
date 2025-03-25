package business

import (
	"errors"
	"user-service/internal/models"
	"user-service/internal/repository"

	"github.com/google/uuid"
)

type UserService interface {
	Login(username, password string) (*models.User, error)
	Register(user *models.User) error
	HasPermission(userID uuid.UUID, permissionName string) bool
	GetUserProfile(userID uuid.UUID) (*models.User, error)
	UpdateUserProfile(user *models.User) error
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

	// Set default role if not provided (e.g., regular user role)
	if user.RoleID == 0 {
		user.RoleID = 1 // Assuming 1 is the default user role
	}

	return s.Repo.Register(user)
}

// HasPermission checks if user has specific permission
func (s *UserServiceImpl) HasPermission(userID uuid.UUID, permissionName string) bool {
	return s.Repo.HasPermission(userID, permissionName)
}

// GetUserProfile retrieves user profile (without sensitive data)
func (s *UserServiceImpl) GetUserProfile(userID uuid.UUID) (*models.Profile, error) {
	profile, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// UpdateUserProfile updates user information
func (s *UserServiceImpl) UpdateUserProfile(profile *models.Profile) error {
	// Add any business logic/validation here before updating
	return s.Repo.UpdateProfile(profile)
}