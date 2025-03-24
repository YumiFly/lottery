// models/role.go
package models

// Role 表示角色表结构
type Role struct {
	RoleID      uint       `gorm:"column:role_id" json:"role_id"`         // 角色 ID
	RoleName    string     `gorm:"column:role_name" json:"role_name"`     // 角色名称
	RoleType    string     `gorm:"column:role_type" json:"role_type"`     // 角色类型
	Description string     `gorm:"column:description" json:"description"` // 角色描述
	Menus       []RoleMenu `gorm:"-" json:"menus"`                        // 角色菜单，忽略 GORM 映射
}
