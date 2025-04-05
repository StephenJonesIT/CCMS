/*
 * @File: models.profile.go
 * @Description: Defines Profile information will be returned to the clients
 * @Author: Tran Thanh Sang (tranthanhsang.it.la@gmail.com)
 */
package models

import "github.com/google/uuid"

type Profile struct {
	ProfileID   int64  		`gorm:"column:profile_id;primaryKey" json:"profile_id,omitempty"`
	FullName    string 		`gorm:"column:fullname" json:"fullname,omitempty"`
	Email       string 		`gorm:"column:email" json:"email,omitempty"`
	Bio         string 		`gorm:"column:bio" json:"bio,omitempty"`
	Picture_URL string 		`gorm:"column:picture_url" json:"image_url,omitempty"`
	UserID      uuid.UUID	`gorm:"foreignKey:user_id" json:"user_id"`
	Interests   []string 	`gorm:"column:interests" json:"interests"`
	Follower_Count int64	`gorm:"column:follower_count" json:"follower_count"`
	Followee_Count int64	`gorm:"column:followee_count" json:"followee_count"`
}

func(Profile) TableName() string{
	return "user_profiles"
}
