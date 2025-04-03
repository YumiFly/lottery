// models/role.go
package models

import "time"

// Role 表示角色表结构
type Role struct {
	RoleID      int        `gorm:"primaryKey;autoIncrement:false" json:"role_id"`
	RoleName    string     `gorm:"size:50;not null" json:"role_name"`
	RoleType    string     `gorm:"size:50" json:"role_type"`
	Description string     `gorm:"type:text" json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Menus       []RoleMenu `gorm:"-" json:"menus"` // 角色菜单，忽略 GORM 映射
}
