package models

import "time"

// RoleMenu 角色菜单表模型
type RoleMenu struct {
	RoleMenuID int       `gorm:"primaryKey;autoIncrement:false" json:"role_menu_id"`
	RoleID     int       `gorm:"not null" json:"role_id"`
	MenuName   string    `gorm:"size:50;not null" json:"menu_name"`
	MenuPath   string    `gorm:"size:255;not null" json:"menu_path"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
