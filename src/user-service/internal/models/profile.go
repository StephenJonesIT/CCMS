/*
 * @File: models.profile.go
 * @Description: Defines Profile information will be returned to the clients
 * @Author: Tran Thanh Sang (tranthanhsang.it.la@gmail.com)
 */
package models

import "github.com/google/uuid"

type Profile struct {
	ProfileID   int64  		`gorm:"column:profile_id;primaryKey" json:"profile_id,omitempty"`
	FullName    string 		`gorm:"column:fullname" json:"fullname"`
	Email       string 		`gorm:"column:email" json:"email,omitempty"`
	Bio         string 		`gorm:"column:bio" json:"bio,omitempty"`
	Picture_URL string 		`gorm:"column:picture_url" json:"image_url,omitempty"`
	UserID      uuid.UUID	`gorm:"foreignKey:user_id" json:"user_id"`
}

func(Profile) TableName() string{
	return "user_profiles"
}
