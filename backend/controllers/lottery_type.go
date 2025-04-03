// controllers/lottery.go
package controllers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateLotteryType 创建彩票类型
// 该方法处理创建彩票类型的 HTTP 请求，验证输入并调用服务层方法。
func CreateLotteryType(c *gin.Context) {
	// 检查权限（确保调用者是 lottery_admin）
	role, exists := c.Get("role")
	if !exists || role != "lottery_admin" {
		c.JSON(http.StatusForbidden, utils.ErrorResponse(utils.ErrCodeForbidden, "Insufficient permissions", nil))
		return
	}

	var lotteryType models.LotteryType
	if err := c.ShouldBindJSON(&lotteryType); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrCodeInvalidInput, "Invalid input", err.Error()))
		return
	}

	lotteryType.TypeID = uuid.NewString() // 生成新的 UUID 作为彩票类型的 ID

	if err := services.CreateLotteryType(&lotteryType); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to create lottery type", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Lottery type created successfully", lotteryType))
}

// GetAllLotteryTypes 获取所有彩票类型
func GetAllLotteryTypes(c *gin.Context) {
	lotteryTypes, err := services.GetAllLotteryTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to get lottery types", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Lottery types retrieved successfully", lotteryTypes))
}
