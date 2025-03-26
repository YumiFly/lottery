// controllers/auth.go
package controllers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginRequest 用户登录请求结构体
type LoginRequest struct {
	WalletAddress string `json:"wallet_address" validate:"required"`
}

// LoginResponse 用户登录响应结构体
type LoginResponse struct {
	CustomerAddress string            `json:"customer_address"`
	Role            string            `json:"role"`
	Menus           []models.RoleMenu `json:"menus"`
	Token           string            `json:"token"`
}

// RefreshResponse 刷新 Token 响应结构体
type RefreshResponse struct {
	Token string `json:"token"`
}

//TODO支持使用邮箱登录
// Register godoc

// Login godoc
// @Summary 用户登录
// @Description 认证用户并返回 JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param login body LoginRequest true "登录凭据"
// @Success 200 {object} utils.Response{data=LoginResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.WithField("error", err.Error()).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrCodeBadRequest, "Invalid request body", err.Error()))
		return
	}
	utils.Logger.WithField("wallet_address", req.WalletAddress).Info("Login request received")

	//TODO 验证请求参数是钱包地址
	// if err := utils.ValidateWalletAddress(req.WalletAddress); err != nil {
	// 	utils.Logger.WithField("error", err.Error()).Error("Invalid wallet address")
	// 	c.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrCodeBadRequest, "Invalid wallet address", err.Error()))
	// 	return
	// }

	result, err := services.Login(req.WalletAddress, c.ClientIP())
	if err != nil {
		utils.Logger.WithField("error", err.Error()).Error("Login failed")
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Login failed", err.Error()))
		return
	}
	resp := LoginResponse{
		CustomerAddress: result.Customer.CustomerAddress,
		Role:            result.Role.RoleName,
		Menus:           result.Role.Menus,
		Token:           result.Token,
	}
	utils.Logger.WithField("customer_address", result.Customer.CustomerAddress).Info("Login successful")
	c.JSON(http.StatusOK, utils.SuccessResponse("Login successful", resp))
}

// RefreshToken godoc
// @Summary 刷新 JWT
// @Description 使用现有 JWT 刷新新 Token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=RefreshResponse}
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /auth/refresh [post]
func RefreshToken(c *gin.Context) {
	customerAddress, _ := c.Get("customer_address")
	role, _ := c.Get("role")

	newToken, err := services.RefreshToken(customerAddress.(string), role.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to refresh token", err.Error()))
		return
	}
	resp := RefreshResponse{
		Token: newToken,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Token refreshed successfully", resp))
}
