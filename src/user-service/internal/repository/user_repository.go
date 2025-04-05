/*
 * @File: repository.user_repository.go
 * @Description: Implements User CRUD functions for Posgres
 * @Author: Tran Thanh Sang (tranthanhsang.it.la@gmail.com)
 */
package repository

import (
	"errors"
	"fmt"

	"github.com/StephenJonesIT/CCMS/src/user-service/common"
	"github.com/StephenJonesIT/CCMS/src/user-service/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)


type UserRepository interface {
	Login(username, password string) (*models.User, error)
	Register(item *models.UserRegister) (err error)
	HasPermission(userID uuid.UUID, permissionName string) bool
	GetUserByID(userID uuid.UUID) (*models.User, error)
	GetListUser(paging *common.Paging) ([]models.User, error)
	ChangePassword(username, password string) error
}

type UserRepoImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepoImpl {
	return &UserRepoImpl{DB: db}
}

func (repo *UserRepoImpl) Login(username, password string) (*models.User, error) {
	var user models.User
	err := repo.DB.Preload("Role").
		Table(models.User{}.TableName()).
		Where("username = ?", username).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database error")
	}

	if err := common.CheckPassword(user.Password, password); err != nil {
		return nil, fmt.Errorf("invalid password") // More specific error
	}

	return &user, nil
}

func (repo *UserRepoImpl) Register(user *models.UserRegister) error {
	if exists, _ := repo.usernameExists(user.UserName); exists {
		return fmt.Errorf("username already exists")
	}

	user.UserID = uuid.New()

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
		Where("users.user_id = ? AND permissions.permission_name = ?", userID, permissionName).
		Count(&count).Error
	return err == nil && count > 0
}

func (repo *UserRepoImpl) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User

	if err := repo.DB.Table(models.User{}.TableName()).Preload("Role").Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepoImpl) GetListUser(paging *common.Paging) ([]models.User, error) {
	var users []models.User
	if err := repo.DB.Table(models.User{}.TableName()).Count(&paging.Total).Error; err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	query := repo.DB.Table(models.User{}.TableName()).
		Preload("Role").
		Order("created_at DESC").
		Offset((paging.Page - 1) * paging.Limit).
		Limit(paging.Limit)

    if err := query.Find(&users).Error; err != nil {
        return nil, fmt.Errorf("failed to get user list: %w", err)
    }
	return users, nil
}

func (repo *UserRepoImpl) ChangePassword(username, newPassword string) error {
    // 1. Start a transaction for atomic operation
    tx := repo.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // 2. Find the user by username
    var user models.User
    if err := tx.Table(models.User{}.TableName()).
        Where("username = ?", username).
        First(&user).Error; err != nil {
        tx.Rollback()
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return fmt.Errorf("user not found")
        }
        return fmt.Errorf("database error: %w", err)
    }

    // 3. Hash the new password
    hashedPassword, err := common.HassPassword(newPassword)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to hash password: %w", err)
    }

    // 4. Update the password
    if err := tx.Table(models.User{}.TableName()).
        Where("user_id = ?", user.UserID).
        Update("password", hashedPassword).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to update password: %w", err)
    }

    // 5. Commit the transaction
    if err := tx.Commit().Error; err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    return nil
}
