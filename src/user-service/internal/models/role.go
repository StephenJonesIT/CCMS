/*
 * @File: models.role.go
 * @Description: Defines Role information will be returned to the clients
 * @Author: Tran Thanh Sang (tranthanhsang.it.la@gmail.com)
 */
package models

type Role struct {
	RoleID 		uint 	`json:"role_id" gorm:"column:role_id"`
	RoleName 	string	`json:"role_name" gorm:"column:role_name"`
}

func(Role) TableName() string {
	return "roles"
}

type Role_Permissions struct {
	RoleID  		uint 	`gorm:"column:role_id;primaryKey"`
	PermissionID	uint	`gorm:"column:permission_id;primaryKey"`
}

func(Role_Permissions) TableName() string{
	return "role_permissions"
}