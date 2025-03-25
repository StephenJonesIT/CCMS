package repository

import (
	"fmt"
	"user-service/common"
	"user-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Login(username, password string) (*models.User, error)
	Register(item *models.User) (err error)
	HasPermission(userID uuid.UUID, permissionName string) bool
	GetUserByID(userID uuid.UUID) (*models.Profile, error)
	UpdateProfile(profile *models.Profile) error
}

type UserRepoImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepoImpl {
	return &UserRepoImpl{DB: db}
}

func (repo *UserRepoImpl) Login(username, password string) (*models.User, error) {
	var user models.User
	if err := repo.DB.Preload("Role").
		Table(models.User{}.TableName()).
		Where("username = ?", username).
		First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if err := common.CheckPassword(user.Password, password); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	return &user, nil
}

func (repo *UserRepoImpl) Register(user *models.User) error {
    if exists, _ := repo.usernameExists(user.UserName); exists {
        return fmt.Errorf("username already exists")
    }

    hashPassword, err := common.HassPassword(user.Password)
    if err != nil {
        return fmt.Errorf("failed to hash password: %v", err)
    }

    user.Password = hashPassword
    return repo.DB.Create(user).Error
}

// Helper function
func (repo *UserRepoImpl) usernameExists(username string) (bool, error) {
    var count int64
    err := repo.DB.Model(&models.User{}).
        Where("username = ?", username).
        Count(&count).Error
    return count > 0, err
}

func (repo *UserRepoImpl) HasPermission(userID uuid.UUID, permissionName string) bool {
	var count int64
	err := repo.DB.Table(models.Role_Permissions{}.TableName()).
		Joins("JOIN users on role_permissions.role_id = users.role_id").
		Joins("JOIN permissions on role_permissions.permission_id = permissions.permission_id").
		Where("users.user_id = ? AND permissions.permission_name = ?",userID, permissionName).
		Count(&count).Error
	return err == nil && count > 0
}
