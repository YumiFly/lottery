// services/auth.go
package services

import (
	"backend/config"
	"backend/db"
	"backend/models"
	"errors"
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
func Login(walletAddress, ip string) (*LoginResult, error) {

	//TODO : 检查 IP 地址是否在黑名单中
	//TODO: 检查 IP 地址是否在白名单中

	//TODO: 检查 walletAddress 是否在黑名单中
	//TODO: 检查 walletAddress 是否在白名单中

	// 根据 walletAddress 从Customer表中查询用户信息,验证用户是否通过KYC
	var customer models.Customer
	if err := db.DB.Where("customer_address = ?", walletAddress).First(&customer).Error; err != nil {
		return nil, err
	}
	// 检查用户是否通过KYC,如果没有通过KYC，则返回错误
	if !customer.IsVerified {
		return nil, errors.New("用户未通过KYC")

	}
	// 根据用户角色查询角色信息
	var role models.Role
	if err := db.DB.Where("role_id = ?", customer.RoleID).First(&role).Error; err != nil {
		return nil, err
	}

	//根据角色ID查询菜单信息
	if err := db.DB.Where("role_id = ?", customer.RoleID).Find(&role.Menus).Error; err != nil {
		return nil, err
	}

	// 生成 JWT 令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"customer_address": customer.CustomerAddress,
		"role":             role.RoleName,
		"exp":              time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		Customer: &customer,
		Role:     &role,
		Token:    tokenString,
	}, nil
}

// RefreshToken 刷新 JWT 令牌
func RefreshToken(customerAddress, role string) (string, error) {
	// 生成新的 JWT 令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"customer_address": customerAddress,
		"role":             role,
		"exp":              time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
