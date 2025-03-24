// models/menu.go
package models

// RoleMenu 表示角色菜单表结构
type RoleMenu struct {
	RoleMenuID uint   `gorm:"column:role_menu_id" json:"role_menu_id"` // 角色菜单 ID
	RoleID     uint   `gorm:"column:role_id" json:"role_id"`           // 角色 ID
	MenuName   string `gorm:"column:menu_name" json:"menu_name"`       // 菜单名称
	MenuPath   string `gorm:"column:menu_path" json:"menu_path"`       // 菜单路径
}
