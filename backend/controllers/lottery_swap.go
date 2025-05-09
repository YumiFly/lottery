// controllers/lottery_swap.go
package controllers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 设置稳定币
// 该方法处理设置稳定币的 HTTP 请求，验证输入并调用服务层方法。
func SetStableCoin(c *gin.Context) {
	var stbCoin models.LotterySTBCoin
	if err := c.ShouldBindJSON(&stbCoin); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrCodeInvalidInput, "Invalid input", err.Error()))
		return
	}

	err := services.SetStableCoin(&stbCoin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to set stable coin", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("SetStableCoin successfully", stbCoin))
}

// 删除稳定币
// 该方法处理设置稳定币的 HTTP 请求，验证输入并调用服务层方法。
func RemoveStableCoin(c *gin.Context) {
	var stbCoin models.LotterySTBCoin
	if err := c.ShouldBindJSON(&stbCoin); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrCodeInvalidInput, "Invalid input", err.Error()))
		return
	}

	err := services.RemoveStableCoin(&stbCoin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to set stable coin", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("RemoveStableCoin successfully", stbCoin))
}
