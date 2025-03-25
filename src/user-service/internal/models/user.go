package models
import "github.com/google/uuid"
type User struct {
	UserID  	uuid.UUID 	`json:"user_id"  gorm:"column:user_id;primaryKey;type:uuid;default:gen_random_uuid()"`
	UserName 	string 		`json:"username" gorm:"column:username"`
	Password 	string		`json:"password" gorm:"column:password"`
	RoleID 		uint		`json:"role_id,omitempty" gorm:"column:role_id"` 
	Role 		Role		
}

func(User) TableName() string {
	return "users"
}