/*
 * @File: models.permission.go
 * @Description: Defines Permission information will be returned to the clients
 * @Author: Tran Thanh Sang (tranthanhsang.it.la@gmail.com)
 */
package models

type Permission struct {
	PermissionID uint `gorm:"column:permission_id"`
	PermissionName string `gorm:"column:permission_name"`
}

func(Permission) TableName() string {
	return "permissions"
}