// models/user.go
package models

import (
	"time"
)

// Customer 表示用户表结构
type Customer struct {
	CustomerAddress  string                   `gorm:"primaryKey;size:255" json:"customer_address"`
	IsVerified       bool                     `gorm:"default:false" json:"is_verified"`
	VerifierAddress  string                   `gorm:"size:255" json:"verifier_address"`
	VerificationTime time.Time                `gorm:"type:timestamptz" json:"verification_time"`
	RegistrationTime time.Time                `gorm:"type:timestamptz;default:now()" json:"registration_time"`
	RoleID           int                      `gorm:"not null" json:"role_id"`
	AssignedDate     time.Time                `gorm:"type:timestamptz;default:now()" json:"assigned_date"`
	CreatedAt        time.Time                `gorm:"type:timestamptz;default:now()" json:"created_at"`
	UpdatedAt        time.Time                `gorm:"type:timestamptz;default:now()" json:"updated_at"`
	KYCData          KYCData                  `gorm:"-" json:"kyc_data"`          // KYC 数据，忽略 GORM 映射
	KYCVerifications []KYCVerificationHistory `gorm:"-" json:"kyc_verifications"` // KYC 验证历史，忽略 GORM 映射
	Role             Role                     `gorm:"-" json:"role"`              // 角色信息，忽略 GORM 映射
}

// KYCData KYC 数据表模型
type KYCData struct {
	CustomerAddress    string    `gorm:"primaryKey;size:255" json:"customer_address"`
	Name               string    `gorm:"size:100" json:"name"`
	BirthDate          time.Time `gorm:"type:date" json:"birth_date"`
	Nationality        string    `gorm:"size:50" json:"nationality"`
	ResidentialAddress string    `gorm:"type:text" json:"residential_address"`
	PhoneNumber        string    `gorm:"size:20" json:"phone_number"`
	Email              string    `gorm:"size:255" json:"email"`
	DocumentType       string    `gorm:"size:50" json:"document_type"`
	DocumentNumber     string    `gorm:"size:50" json:"document_number"`
	FilePath           string    `gorm:"type:text" json:"file_path"`
	SubmissionDate     time.Time `gorm:"type:timestamptz;default:now()" json:"submission_date"`
	RiskLevel          string    `gorm:"size:20" json:"risk_level"`
	SourceOfFunds      string    `gorm:"type:text" json:"source_of_funds"`
	Occupation         string    `gorm:"size:100" json:"occupation"`
	CreatedAt          time.Time `gorm:"type:timestamptz;default:now()" json:"created_at"`
	UpdatedAt          time.Time `gorm:"type:timestamptz;default:now()" json:"updated_at"`
}

// KYCVerificationHistory 表示 KYC 验证历史表结构
type KYCVerificationHistory struct {
	HistoryID        int       `gorm:"primaryKey;autoIncrement:false" json:"history_id"`
	CustomerAddress  string    `gorm:"size:255;not null" json:"customer_address"`
	VerifyStatus     string    `gorm:"size:50" json:"verify_status"`
	VerifierAddress  string    `gorm:"size:255" json:"verifier_address"`
	VerificationDate time.Time `gorm:"type:timestamptz;default:now()" json:"verification_date"`
	Comments         string    `gorm:"type:text" json:"comments"`
	CreatedAt        time.Time `gorm:"type:timestamptz;default:now()" json:"created_at"`
	UpdatedAt        time.Time `gorm:"type:timestamptz;default:now()" json:"updated_at"`
}
