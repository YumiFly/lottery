// services/user.go
package services

import (
	"backend/db"
	"backend/models"
	"backend/utils"
	"errors"
	"time"
)

// GetCustomers 获取所有用户及其关联数据
func GetCustomers() ([]models.Customer, error) {
	var customers []models.Customer
	// 查询所有用户记录
	if err := db.DB.Find(&customers).Error; err != nil {
		return nil, err
	}
	utils.Logger.Info("customers from database ", customers)
	// 遍历用户，手动查询关联数据
	for i := range customers {
		// 查询 KYCData
		var kycData models.KYCData
		if err := db.DB.Where("customer_address = ?", customers[i].CustomerAddress).First(&kycData).Error; err == nil {
			customers[i].KYCData = kycData
		}

		// 查询 KYCVerifications
		var kycVerifications []models.KYCVerificationHistory
		if err := db.DB.Where("customer_address = ?", customers[i].CustomerAddress).Find(&kycVerifications).Error; err == nil {
			customers[i].KYCVerifications = kycVerifications
		}

		// 查询 Role（如果 role_id 不为 0）
		if customers[i].RoleID != 0 {
			var role models.Role
			if err := db.DB.Where("role_id = ?", customers[i].RoleID).First(&role).Error; err == nil {
				customers[i].Role = role
			}

			// 查询 Role.Menus
			var menus []models.RoleMenu
			if err := db.DB.Where("role_id = ?", customers[i].RoleID).Find(&menus).Error; err == nil {
				customers[i].Role.Menus = menus
			}
		}
	}

	return customers, nil
}

// GetCustomerByAddress 根据 CustomerAddress 获取用户及其关联数据
func GetCustomerByAddress(customerAddress string) (*models.Customer, error) {
	var customer models.Customer
	// 查询指定用户
	if err := db.DB.Where("customer_address = ?", customerAddress).First(&customer).Error; err != nil {
		return nil, err
	}
	utils.Logger.Info("GetCustomerByAddress's customer from database ,", customer)

	// 手动查询关联数据
	// 查询 KYCData
	var kycData models.KYCData
	if err := db.DB.Where("customer_address = ?", customer.CustomerAddress).First(&kycData).Error; err == nil {
		customer.KYCData = kycData
	}
	utils.Logger.Info("GetCustomerByAddress's kycData from database ,", kycData)

	// 查询 KYCVerifications
	var kycVerifications []models.KYCVerificationHistory
	if err := db.DB.Where("customer_address = ?", customer.CustomerAddress).Find(&kycVerifications).Error; err == nil {
		customer.KYCVerifications = kycVerifications
	}
	utils.Logger.Info("GetCustomerByAddress's kycVerifications from database ,", kycVerifications)

	// 查询 Role（如果 role_id 不为 0）
	if customer.RoleID != 0 {
		var role models.Role
		if err := db.DB.Where("role_id = ?", customer.RoleID).First(&role).Error; err == nil {
			customer.Role = role
		}

		// 查询 Role.Menus
		var menus []models.RoleMenu
		if err := db.DB.Where("role_id = ?", customer.RoleID).Find(&menus).Error; err == nil {
			customer.Role.Menus = menus
		}
	}
	return &customer, nil
}

// CreateCustomer 创建用户，仅插入 customers 和 kyc_data 表
func CreateCustomer(customer *models.Customer) error {
	// 使用事务确保数据一致性
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 插入用户记录
	if err := tx.Create(customer).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 插入 KYCData
	if customer.KYCData.CustomerAddress != "" {
		if err := tx.Create(&customer.KYCData).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 不插入 KYCVerifications，留给验证流程处理
	return tx.Commit().Error
}

// VerifyCustomer 验证用户 KYC 信息
func VerifyCustomer(verification *models.KYCVerificationHistory) error {
	// 插入验证记录（无论 Approved 还是 Rejected 都会插入）
	if err := db.DB.Create(verification).Error; err != nil {
		return err
	}

	// 查询用户
	var customer models.Customer
	if err := db.DB.Where("customer_address = ?", verification.CustomerAddress).First(&customer).Error; err != nil {
		return err
	}
	utils.Logger.Info("VerifyCustomer'customer from database , ", customer)

	// 根据验证状态更新用户记录
	if verification.VerifyStatus == "Approved" {
		// 验证通过，更新 Customer 记录
		updates := map[string]interface{}{
			"is_verified":       true,
			"verifier_address":  verification.VerifierAddress,
			"verification_time": verification.VerificationDate,
			"role_id":           2, // 分配 normal_user 角色
			"assigned_date":     time.Now(),
		}
		// 明确指定更新条件
		if err := db.DB.Model(&customer).Where("customer_address = ?", customer.CustomerAddress).Updates(updates).Error; err != nil {
			return err
		}

		//TODO 调用链上KYC的相关函数，更新KYC状态
		updateKycStatusOnChain()
	} else if verification.VerifyStatus != "Rejected" {
		// 如果 verify_status 既不是 Approved 也不是 Rejected，返回错误
		return errors.New("invalid verify_status, must be 'Approved' or 'Rejected'")
	}

	// Rejected 状态下无需更新 Customer 表，仅插入验证记录
	return nil
}
func updateKycStatusOnChain() {

}
