// controllers/user.go
package controllers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetCustomers godoc
// @Summary 获取所有用户
// @Description 获取用户列表
// @Tags customers
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.Customer}
// @Failure 500 {object} utils.Response
// @Router /customers [get]
func GetCustomers(c *gin.Context) {
	customers, err := services.GetCustomers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to retrieve customers", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Customers retrieved successfully", customers))
}

// GetCustomerByAddress godoc
// @Summary 根据 CustomerAddress 获取用户
// @Description 根据用户地址获取详细信息
// @Tags customers
// @Accept json
// @Produce json
// @Param customer_address path string true "用户地址"
// @Success 200 {object} utils.Response{data=models.Customer}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /customers/{customer_address} [get]
func GetCustomerByAddress(c *gin.Context) {
	customerAddress := c.Param("customer_address")
	if customerAddress == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrCodeInvalidInput, "Customer address is required", nil))
		return
	}
	customer, err := services.GetCustomerByAddress(customerAddress)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, utils.ErrorResponse(utils.ErrCodeInvalidInput, "Customer not found", nil))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to retrieve customer", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Customer retrieved successfully", customer))
}

// CreateCustomer godoc
// @Summary 注册用户
// @Description 创建新用户，不需要权限控制，仅插入 customers 和 kyc_data 表
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body models.Customer true "用户信息"
// @Success 201 {object} utils.Response{data=models.Customer}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /customers [post]
func CreateCustomer(c *gin.Context) {
	// 从上下文获取 ValidationMiddleware 绑定的 customer 对象
	customer, exists := c.Get("validated_obj")
	if !exists {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to retrieve validated customer", nil))
		return
	}

	// 转换为 models.Customer 类型
	cust, ok := customer.(*models.Customer)
	if !ok {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Invalid customer type", nil))
		return
	}

	// 注册时不分配角色，role_id 设为 0（未分配）
	cust.RoleID = 0

	// 调用服务层创建用户
	if err := services.CreateCustomer(cust); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to create customer", err.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.SuccessResponse("Customer registered successfully", cust))
}

// VerifyCustomer godoc
// @Summary 验证用户 KYC 信息
// @Description 管理员验证用户 KYC 信息，更新验证状态
// @Tags customers
// @Accept json
// @Produce json
// @Param verification body models.KYCVerificationHistory true "验证信息"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /auth/verify [post]
func VerifyCustomer(c *gin.Context) {
	// 检查权限（确保调用者是 lottery_admin）
	role, exists := c.Get("role")
	if !exists || role != "lottery_admin" {
		c.JSON(http.StatusForbidden, utils.ErrorResponse(utils.ErrCodeForbidden, "Insufficient permissions", nil))
		return
	}

	// 从上下文获取 ValidationMiddleware 绑定的 verification 对象
	var verification models.KYCVerificationHistory
	//verification.HistoryID = uint(uuid.New()[8:]).
	// 绑定 JSON 请求体到结构体
	if err := c.ShouldBindJSON(&verification); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrCodeInvalidInput, "Invalid input", err.Error()))
		return
	}
	logrus.Info("verification: ", verification)
	// 调用服务层进行验证
	if err := services.VerifyCustomer(&verification); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to verify customer", err.Error()))
		return
	}

	// 根据验证状态返回响应
	if verification.VerifyStatus == "Approved" {
		c.JSON(http.StatusOK, utils.SuccessResponse("Verification successful", nil))
	} else if verification.VerifyStatus == "Rejected" {
		c.JSON(http.StatusOK, utils.SuccessResponse("Verification failed", nil))
	}
}
