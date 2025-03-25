package models

type Permission struct {
	PermissionID uint `gorm:"column:permission_id"`
	PermissionName string `gorm:"column:permission_name"`
}

func(Permission) TableName() string {
	return "permissions"
}