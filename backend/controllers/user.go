// controllers/user.go
package controllers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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

func UploadPhoto(c *gin.Context) {
	// 获取上传的文件（表单字段名为 "photo"）
	file, err := c.FormFile("idPhoto")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrCodeInvalidInput, "Failed to get photo from form", err.Error()))
		return
	}

	// 验证文件类型（只允许图片）
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrCodeInvalidInput, "Only JPG, JPEG, and PNG files are allowed", nil))
		return
	}

	// 验证文件大小（例如限制为 5MB）
	const maxSize = 5 << 20 // 5MB
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrCodeInvalidInput, "File size exceeds 5MB limit", nil))
		return
	}

	// 检查是否配置了S3
	var fileURL string
	if utils.S3Client != nil {
		// 使用S3存储
		url, err := utils.UploadFileToS3(file, "photos")
		if err != nil {
			utils.Logger.WithField("error", err.Error()).Error("Failed to upload file to S3")
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to upload file to S3", err.Error()))
			return
		}
		fileURL = url
	} else {
		// 使用本地存储作为备选
		// 创建存储路径（uploads/目录）
		uploadDir := "uploads"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to create upload directory", err.Error()))
			return
		}

		// 生成唯一的文件名（使用时间戳和原始文件名）
		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		filePath := filepath.Join(uploadDir, filename)

		// 保存文件到指定路径
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to save photo", err.Error()))
			return
		}

		// 构造文件访问地址（假设服务端运行在 localhost:8080）
		// 实际部署时，应使用真实的域名或 CDN 地址
		fileURL = fmt.Sprintf("/%s", filePath)
	}

	// 返回成功响应
	c.JSON(http.StatusOK, utils.SuccessResponse("Photo uploaded successfully", map[string]string{
		"file_url": fileURL,
	}))
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

	// 注册时间
	cust.RegistrationTime = time.Now()

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
	if !exists || role != "admin" {
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

	verification.VerificationDate = time.Now()
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
