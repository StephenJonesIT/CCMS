/*
 * @File: models.user.go
 * @Description: Defines User information will be returned to the clients
 * @Author: Tran Thanh Sang (tranthanhsang.it.la@gmail.com)
 */
package models
import "github.com/google/uuid"
type User struct {
	UserID  	uuid.UUID 	`json:"user_id"  gorm:"column:user_id;primaryKey;type:uuid;default:gen_random_uuid()"`
	UserName 	string 		`json:"username" gorm:"column:username"`
	Password 	string		`json:"password" gorm:"column:password"`
	RoleID 		uint		`json:"role_id" gorm:"column:role_id"` 
	Role 		Role		`gorm:"foreignKey:RoleID;references:RoleID" json:"role"`
}

func(User) TableName() string {
	return "users"
}

type UserRegister struct {
	UserID  	uuid.UUID 	`json:"user_id,omitempty"  gorm:"column:user_id;primaryKey;type:uuid;default:gen_random_uuid()"`
	UserName 	string 		`json:"username" gorm:"column:username"`
	Password 	string		`json:"password" gorm:"column:password"`
	RoleID 		uint		`json:"role_id" gorm:"column:role_id"`
}

func(UserRegister) TableName() string {
	return User{}.TableName()
}