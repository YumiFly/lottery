// controllers/lottery.go
package controllers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// cleanString 清洗字符串，确保只包含有效的 UTF-8 字符
func cleanString(s string) string {
	if !utf8.ValidString(s) {
		// 如果字符串包含无效的 UTF-8 字符，移除非 UTF-8 字符
		var builder strings.Builder
		for _, r := range s {
			if utf8.ValidRune(r) {
				builder.WriteRune(r)
			}
		}
		return builder.String()
	}
	return s
}

// CreateLotteryType 创建彩票类型
// 该方法处理创建彩票类型的 HTTP 请求，验证输入并调用服务层方法。
func CreateLotteryType(c *gin.Context) {
	var lotteryType models.LotteryType
	if err := c.ShouldBindJSON(&lotteryType); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrCodeInvalidInput, "Invalid input", err.Error()))
		return
	}

	lotteryType.TypeID = uuid.NewString() // 生成新的 UUID 作为彩票类型的 ID
	// 清洗字符串字段，确保只包含有效的 UTF-8 字符
	lotteryType.TypeID = cleanString(lotteryType.TypeID)
	lotteryType.TypeName = cleanString(lotteryType.TypeName)
	lotteryType.Description = cleanString(lotteryType.Description)

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
