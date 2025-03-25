package models

import "github.com/google/uuid"

type Profile struct {
	ProfileID   int64  		`gorm:"column:profile_id;primaryKey"`
	FullName    string 		`gorm:"column:fullname"`
	Email       string 		`gorm:"column:email"`
	Bio         string 		`gorm:"column:bio"`
	Picture_URL string 		`gorm:"column:picture_url"`
	UserID      uuid.UUID	`gorm:"foreignKey:user_id"`
}

func(Profile) TableName() string{
	return "user_profiles"
}