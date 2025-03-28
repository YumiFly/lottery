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

// CreateLottery 创建彩票
// 该方法处理创建彩票的 HTTP 请求，验证输入并调用服务层方法。
func CreateLottery(c *gin.Context) {
	var lottery models.Lottery
	if err := c.ShouldBindJSON(&lottery); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrCodeInvalidInput, "Invalid input", err.Error()))
		return
	}

	// 清洗字符串字段
	lottery.LotteryID = cleanString(lottery.LotteryID)
	lottery.TypeID = cleanString(lottery.TypeID)
	lottery.TicketName = cleanString(lottery.TicketName)
	lottery.TicketPrice = cleanString(lottery.TicketPrice)
	lottery.BettingRules = cleanString(lottery.BettingRules)
	lottery.PrizeStructure = cleanString(lottery.PrizeStructure)
	lottery.ContractAddress = cleanString(lottery.ContractAddress)

	if err := services.CreateLottery(&lottery); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to create lottery", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Lottery created successfully", lottery))
}

// CreateIssue 创建彩票期号
// 该方法处理创建彩票期号的 HTTP 请求，验证输入并调用服务层方法。
func CreateIssue(c *gin.Context) {
	var issue models.LotteryIssue
	if err := c.ShouldBindJSON(&issue); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrCodeInvalidInput, "Invalid input", err.Error()))
		return
	}

	// 清洗字符串字段
	issue.IssueID = cleanString(issue.IssueID)
	issue.LotteryID = cleanString(issue.LotteryID)
	issue.IssueNumber = cleanString(issue.IssueNumber)
	issue.PrizePool = cleanString(issue.PrizePool)
	issue.DrawStatus = cleanString(issue.DrawStatus)
	issue.WinningNumbers = cleanString(issue.WinningNumbers)
	issue.RandomSeed = cleanString(issue.RandomSeed)
	issue.DrawTxHash = cleanString(issue.DrawTxHash)

	if err := services.CreateIssue(&issue); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to create issue", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Issue created successfully", issue))
}

// PurchaseTicket 购买彩票
// 该方法处理购买彩票的 HTTP 请求，验证输入并调用服务层方法。
func PurchaseTicket(c *gin.Context) {
	var ticket models.LotteryTicket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrCodeInvalidInput, "Invalid input", err.Error()))
		return
	}

	// 清洗字符串字段
	ticket.TicketID = cleanString(ticket.TicketID)
	ticket.IssueID = cleanString(ticket.IssueID)
	ticket.BuyerAddress = cleanString(ticket.BuyerAddress)
	ticket.BetContent = cleanString(ticket.BetContent)
	ticket.PurchaseAmount = cleanString(ticket.PurchaseAmount)
	ticket.TransactionHash = cleanString(ticket.TransactionHash)
	ticket.ClaimStatus = cleanString(ticket.ClaimStatus)
	ticket.ClaimTxHash = cleanString(ticket.ClaimTxHash)

	if err := services.PurchaseTicket(&ticket); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to purchase ticket", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Ticket purchased successfully", ticket))
}

// DrawLottery 开奖
// 该方法处理彩票开奖的 HTTP 请求，验证权限并调用服务层方法。
func DrawLottery(c *gin.Context) {
	// 检查权限（确保调用者是 lottery_admin）
	role, exists := c.Get("role")
	if !exists || role != "lottery_admin" {
		c.JSON(http.StatusForbidden, utils.ErrorResponse(utils.ErrCodeForbidden, "Insufficient permissions", nil))
		return
	}

	issueID := c.Query("issue_id")
	if issueID == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrCodeInvalidInput, "Issue ID is required", nil))
		return
	}

	// 清洗 issueID
	issueID = cleanString(issueID)

	if err := services.DrawLottery(issueID); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to draw lottery", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Lottery draw initiated", nil))
}

func GetAllLotteryTypes(c *gin.Context) {
	lotteryTypes, err := services.GetAllLotteryTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to get lottery types", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Lottery types retrieved successfully", lotteryTypes))
}

func GetAllLottery(c *gin.Context) {
	lotteries, err := services.GetAllLotteries()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to get lotteries", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Lotteries retrieved successfully", lotteries))
}

func GetLotteryByType(c *gin.Context) {
	lotteryType := c.Query("typeID")
	if lotteryType == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrCodeInvalidInput, "Lottery type is required", nil))
		return
	}

	lottery, err := services.GetLotteryByTypeID(lotteryType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to get lottery", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Lottery retrieved successfully", lottery))
}

func GetLatestIssueByLotteryID(c *gin.Context) {
	lotteryID := c.Param("lottery_id")
	if lotteryID == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrCodeInvalidInput, "Lottery ID is required", nil))
		return
	}

	issue, err := services.GetLatestIssueByLotteryID(lotteryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to get latest issue", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Latest issue retrieved successfully", issue))
}

// 查询快到期的彩票，并提醒管理员开奖
func GetExpiringIssues(c *gin.Context) {
	issues, err := services.GetExpiringIssues()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to get expiring issues", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Expiring issues retrieved successfully", issues))
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

func GetDrawnLotteryByIssueID(c *gin.Context) {
	issueID := c.Param("issue_id")
	if issueID == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrCodeInvalidInput, "Issue ID is required", nil))
		return
	}

	drawnLottery, err := services.GetDrawnLotteryByIssueID(issueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to get drawn lottery", err.Error()))
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Drawn lottery retrieved successfully", drawnLottery))
}
