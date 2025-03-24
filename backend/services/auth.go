// services/auth.go
package services

import (
	"backend/config"
	"backend/db"
	"backend/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// LoginResult 定义登录返回结果
type LoginResult struct {
	Customer *models.Customer // 用户信息
	Role     *models.Role     // 角色信息
	Token    string           // JWT 令牌
}

// Login 登录逻辑，返回用户、角色和 JWT 令牌
func Login(email, walletAddress, ip string) (*LoginResult, error) {
	var kycData models.KYCData
	// 查询 KYCData 以验证 Email 和 WalletAddress
	if err := db.DB.Where("email = ? AND customer_address = ?", email, walletAddress).First(&kycData).Error; err != nil {
		if err.Error() == "record not found" {
			// 如果用户不存在，创建新用户
			newCustomer := models.Customer{
				CustomerAddress:  walletAddress,
				RoleID:           0, // 不分配角色
				RegistrationTime: time.Now(),
				AssignedDate:     time.Now(),
				KYCData: models.KYCData{
					CustomerAddress: walletAddress,
					Email:           email,
				},
			}
			if err := CreateCustomer(&newCustomer); err != nil {
				return nil, err
			}
			// 重新查询新创建的用户
			var customer models.Customer
			if err := db.DB.Where("customer_address = ?", walletAddress).First(&customer).Error; err != nil {
				return nil, err
			}
			// 创建默认角色对象（未分配角色）
			role := &models.Role{
				RoleID:   0,
				RoleName: "unassigned",
			}
			// 生成 JWT 令牌
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"customer_address": customer.CustomerAddress,
				"email":            newCustomer.KYCData.Email,
				"role":             role.RoleName,
				"exp":              time.Now().Add(time.Hour * 24).Unix(),
			})
			tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
			if err != nil {
				return nil, err
			}
			return &LoginResult{
				Customer: &customer,
				Role:     role,
				Token:    tokenString,
			}, nil
		}
		return nil, err
	}

	// 查询现有用户
	var customer models.Customer
	if err := db.DB.Where("customer_address = ?", walletAddress).First(&customer).Error; err != nil {
		return nil, err
	}

	// 查询角色（如果 role_id 不为 0）
	role := &models.Role{
		RoleID:   0,
		RoleName: "unassigned",
	}
	if customer.RoleID != 0 {
		if err := db.DB.Where("role_id = ?", customer.RoleID).First(role).Error; err != nil {
			return nil, err
		}
		// 查询角色菜单
		var menus []models.RoleMenu
		if err := db.DB.Where("role_id = ?", role.RoleID).Find(&menus).Error; err == nil {
			role.Menus = menus
		}
	}

	// 生成 JWT 令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"customer_address": customer.CustomerAddress,
		"email":            kycData.Email,
		"role":             role.RoleName,
		"exp":              time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		Customer: &customer,
		Role:     role,
		Token:    tokenString,
	}, nil
}

// RefreshToken 刷新 JWT 令牌
func RefreshToken(customerAddress, email, role string) (string, error) {
	// 生成新的 JWT 令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"customer_address": customerAddress,
		"email":            email,
		"role":             role,
		"exp":              time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
