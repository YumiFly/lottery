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

// PurchaseTicket 购买彩票
// 该方法处理购买彩票的 HTTP 请求，验证输入并调用服务层方法。
func PurchaseTicket(c *gin.Context) {
	var ticket models.LotteryTicket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrCodeInvalidInput, "Invalid input", err.Error()))
		return
	}

	ticket.TicketID = uuid.NewString() // 生成新的 UUID 作为彩票的 ID

	// 清洗字符串字段
	ticket.TicketID = cleanString(ticket.TicketID)
	ticket.IssueID = cleanString(ticket.IssueID)
	ticket.BuyerAddress = cleanString(ticket.BuyerAddress)
	ticket.BetContent = cleanString(ticket.BetContent)
	ticket.PurchaseAmount = cleanString(ticket.PurchaseAmount)
	ticket.TransactionHash = cleanString(ticket.TransactionHash)
	ticket.ClaimTxHash = cleanString(ticket.ClaimTxHash)

	if err := services.PurchaseTicket(&ticket); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to purchase ticket", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Ticket purchased successfully", ticket))
}
func GetPurchasedTicketsByCustomerAddress(c *gin.Context) {
	customerAddress := c.Param("customer_address")
	if customerAddress == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrCodeInvalidInput, "Customer address is required", nil))
		return
	}

	tickets, err := services.GetPurchasedTicketsByCustomerAddress(customerAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to get purchased tickets", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Purchased tickets retrieved successfully", tickets))
}
