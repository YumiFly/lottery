// models/user.go
package models

import (
	"time"
)

// Customer 表示用户表结构
type Customer struct {
	CustomerAddress  string                   `gorm:"column:customer_address" json:"customer_address"`   // 用户地址
	IsVerified       bool                     `gorm:"column:is_verified" json:"is_verified"`             // 是否通过 KYC 验证
	VerifierAddress  string                   `gorm:"column:verifier_address" json:"verifier_address"`   // 验证者地址
	VerificationTime time.Time                `gorm:"column:verification_time" json:"verification_time"` // 验证时间
	RegistrationTime time.Time                `gorm:"column:registration_time" json:"registration_time"` // 注册时间
	RoleID           uint                     `gorm:"column:role_id" json:"role_id"`                     // 角色 ID
	AssignedDate     time.Time                `gorm:"column:assigned_date" json:"assigned_date"`         // 分配日期
	KYCData          KYCData                  `gorm:"-" json:"kyc_data"`                                 // KYC 数据，忽略 GORM 映射
	KYCVerifications []KYCVerificationHistory `gorm:"-" json:"kyc_verifications"`                        // KYC 验证历史，忽略 GORM 映射
	Role             Role                     `gorm:"-" json:"role"`                                     // 角色信息，忽略 GORM 映射
}

// KYCData 表示 KYC 数据表结构
type KYCData struct {
	CustomerAddress    string    `gorm:"column:customer_address" json:"customer_address"`       // 用户地址
	Name               string    `gorm:"column:name" json:"name"`                               // 姓名
	BirthDate          time.Time `gorm:"column:birth_date" json:"birth_date"`                   // 出生日期
	Nationality        string    `gorm:"column:nationality" json:"nationality"`                 // 国籍
	ResidentialAddress string    `gorm:"column:residential_address" json:"residential_address"` // 居住地址
	PhoneNumber        string    `gorm:"column:phone_number" json:"phone_number"`               // 电话号码
	Email              string    `gorm:"column:email" json:"email"`                             // 电子邮件
	DocumentType       string    `gorm:"column:document_type" json:"document_type"`             // 证件类型
	DocumentNumber     string    `gorm:"column:document_number" json:"document_number"`         // 证件号码
	FilePath           string    `gorm:"column:file_path" json:"file_path"`                     // 证件文件路径
	SubmissionDate     time.Time `gorm:"column:submission_date" json:"submission_date"`         // 提交日期
	RiskLevel          string    `gorm:"column:risk_level" json:"risk_level"`                   // 风险等级
	SourceOfFunds      string    `gorm:"column:source_of_funds" json:"source_of_funds"`         // 资金来源
	Occupation         string    `gorm:"column:occupation" json:"occupation"`                   // 职业
}

// KYCVerificationHistory 表示 KYC 验证历史表结构
type KYCVerificationHistory struct {
	HistoryID        uint      `gorm:"column:history_id" json:"history_id"`               // 历史记录 ID
	CustomerAddress  string    `gorm:"column:customer_address" json:"customer_address"`   // 用户地址
	VerifyStatus     string    `gorm:"column:verify_status" json:"verify_status"`         // 验证状态：Pending, Approved, Rejected
	VerifierAddress  string    `gorm:"column:verifier_address" json:"verifier_address"`   // 验证者地址
	VerificationDate time.Time `gorm:"column:verification_date" json:"verification_date"` // 验证日期
	Comments         string    `gorm:"column:comments" json:"comments"`                   // 备注
}
