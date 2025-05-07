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
	// var ticket models.LotteryTicket
	// if err := c.ShouldBindJSON(&ticket); err != nil {
	// 	c.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrCodeInvalidInput, "Invalid input", err.Error()))
	// 	return
	// }

	var stbCoin models.LotterySTBCoin
	stbCoin.STBCoinName = "FTK"
	stbCoin.STBCoinAddr = "0x8A791620dd6260079BF849Dc5567aDC3F2FdC318"
	stbCoin.STB2LOTRate = 2
	stbCoin.STBReceiverAddr = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"

	err := services.SetStableCoin(&stbCoin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to purchase ticket", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("SetStableCoin successfully", stbCoin))
}
