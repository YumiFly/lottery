// controllers/lottery.go
package controllers

import (
	"backend/blockchain"
	"backend/db"
	"backend/models"
	"backend/services"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateLottery 创建彩票
// 该方法处理创建彩票的 HTTP 请求，验证输入并调用服务层方法。
func CreateLottery(c *gin.Context) {
	// 检查权限（确保调用者是 lottery_admin）
	// role, exists := c.Get("role")
	// if !exists || (role != "lottery_admin" && role != "admin") {
	// 	c.JSON(http.StatusForbidden, utils.ErrorResponse(utils.ErrCodeForbidden, "Insufficient permissions", nil))
	// 	return
	// }

	var lottery models.Lottery
	if err := c.ShouldBindJSON(&lottery); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrCodeInvalidInput, "Invalid input", err.Error()))
		return
	}

	lottery.LotteryID = uuid.NewString() // 生成新的 UUID 作为彩票的 ID

	if err := services.CreateLottery(&lottery); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to create lottery", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Lottery created successfully", lottery))
}

// CreateIssue 创建彩票期号
// 该方法处理创建彩票期号的 HTTP 请求，验证输入并调用服务层方法。
func CreateIssue(c *gin.Context) {
	// role, exists := c.Get("role")
	// if !exists || (role != "lottery_admin" && role != "admin") {
	// 	c.JSON(http.StatusForbidden, utils.ErrorResponse(utils.ErrCodeForbidden, "Insufficient permissions", nil))
	// 	return
	// }
	var issue models.LotteryIssue
	if err := c.ShouldBindJSON(&issue); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrCodeInvalidInput, "Invalid input", err.Error()))
		return
	}

	issue.IssueID = uuid.NewString() // 生成新的 UUID 作为彩票期号的 ID
	issue.Status = models.IssueStatusPending

	if err := services.CreateIssue(&issue); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to create issue", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Issue created successfully", issue))
}

// DrawLottery 开奖
// 该方法处理彩票开奖的 HTTP 请求，验证权限并调用服务层方法。
func DrawLottery(c *gin.Context) {
	// 检查权限（确保调用者是 lottery_admin）
	// role, exists := c.Get("role")
	// if !exists || (role != "lottery_admin" && role != "admin") {
	// 	c.JSON(http.StatusForbidden, utils.ErrorResponse(utils.ErrCodeForbidden, "Insufficient permissions", nil))
	// 	return
	// }

	issueID := c.Query("issue_id")
	if issueID == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrCodeInvalidInput, "Issue ID is required", nil))
		return
	}

	service := services.NewLotteryService(blockchain.Client, blockchain.Auth, db.DB)

	if err := service.DrawLotteryAsync(issueID); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to draw lottery", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Lottery draw initiated", nil))
}

func GetAllPools(c *gin.Context) {
	pools, err := services.GetAllPools()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to get pools", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Pools retrieved successfully", pools))
}

// GetLotteries 获取所有彩票
func GetAllLottery(c *gin.Context) {
	lotteries, err := services.GetAllLotteries()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to get lotteries", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Lotteries retrieved successfully", lotteries))
}

// GetLotteryByID 根据ID获取彩票
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

// GetUpcomingIssues 获取即将开奖的彩票
func GetUpcomingIssues(c *gin.Context) {
	issues, err := services.GetUpcomingIssues()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to get upcoming issues", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Upcoming issues retrieved successfully", issues))
}

// GetLatestIssueByLotteryID 获取最新期号
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

// GetDrawnLotteryByIssueID 根据期号ID获取开奖结果
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

// GetLatestDrawnLottery 获取最新开奖结果
func GetLatestDrawnLottery(c *gin.Context) {
	drawnLottery, err := services.GetLatestDrawnLottery()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to get latest drawn lottery", err.Error()))
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Latest drawn lottery retrieved successfully", drawnLottery))
}
