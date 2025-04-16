package services

import (
	"backend/blockchain"
	"backend/db"
	"backend/models"
	"backend/utils"
	"context"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// CreateIssue 创建彩票期号
func CreateIssue(issue *models.LotteryIssue) error {
	utils.Logger.Info("Creating issue", "issue_id", issue.IssueID, "lottery_id", issue.LotteryID)

	executeTx := func() (common.Hash, error) {
		var lottery models.Lottery
		if err := db.DB.Where("lottery_id = ?", issue.LotteryID).First(&lottery).Error; err != nil {
			utils.Logger.Warn("Lottery not found", "lottery_id", issue.LotteryID)
			return common.Hash{}, utils.NewServiceError("lottery not found", err)
		}

		var existingIssue models.LotteryIssue
		if err := db.DB.Where("issue_number = ?", issue.IssueNumber).Where("lottery_id = ?", issue.LotteryID).First(&existingIssue).Error; err == nil {
			utils.Logger.Warn("Issue number already exists", "issue_number", issue.IssueNumber)
			return common.Hash{}, utils.NewServiceError("issue number already exists", nil)
		}

		issue.CreatedAt = time.Now()
		issue.UpdatedAt = time.Now()
		issue.PrizePool = 0

		contract, err := connectLotteryContract(lottery.ContractAddress)
		if err != nil {
			return common.Hash{}, utils.NewServiceError("failed to connect to lottery contract", err)
		}

		// 获取当前合约状态
		currentState, err := contract.GetState(nil)
		if err != nil {
			utils.Logger.Error("Failed to get contract state", "error", err)
			return common.Hash{}, utils.NewServiceError("failed to get contract state", err)
		}
		utils.Logger.Info("Current contract state", "state", currentState) // 添加日志记录

		tx, err := contract.TransState(blockchain.Auth, uint8(1)) // 设置为 Distribute 状态
		if err != nil {
			utils.Logger.Error("Failed to set state to Distribute", "error", err)
			if tx != nil {
				return tx.Hash(), utils.NewServiceError("failed to set state to Distribute", err)
			}
			return common.Hash{}, utils.NewServiceError("failed to set state to Distribute", err)
		}

		receipt, err := bind.WaitMined(context.Background(), blockchain.Client, tx)
		if err != nil {
			utils.Logger.Error("Transaction failed", "tx_hash", tx.Hash().Hex(), "error", err)
			return tx.Hash(), utils.NewServiceError("transaction failed", err)
		}
		if receipt.Status != 1 {
			utils.Logger.Error("Transaction failed", "tx_hash", tx.Hash().Hex(), "status", receipt.Status)
			return tx.Hash(), utils.NewServiceError("transaction failed", nil)
		}

		if err := db.DB.Create(issue).Error; err != nil {
			utils.Logger.Error("Failed to save issue to database", "error", err)
			return common.Hash{}, utils.NewServiceError("failed to save issue to database", err)
		}

		utils.Logger.Info("Issue created successfully", "issue_id", issue.IssueID)
		return tx.Hash(), nil
	}

	// 使用空的 data 调用 WithBlockchain，Gas 估算依赖内部逻辑
	data := []byte{}
	return blockchain.WithBlockchain(context.Background(), data, executeTx)
}

// GetUpcomingIssues 获取即将销售的期号
func GetUpcomingIssues() ([]models.LotteryIssue, error) {
	utils.Logger.Info("Fetching upcoming issues")
	var issues []models.LotteryIssue
	if err := db.DB.Where("sale_end_time > ?", time.Now()).Find(&issues).Error; err != nil {
		utils.Logger.Error("Failed to fetch upcoming issues", "error", err)
		return nil, utils.NewServiceError("failed to fetch upcoming issues", err)
	}
	utils.Logger.Info("Fetched upcoming issues", "count", len(issues))
	return issues, nil
}

// GetLatestIssueByLotteryID 根据彩票ID获取最新期号
func GetLatestIssueByLotteryID(lotteryID string) (*models.LotteryIssue, error) {
	utils.Logger.Info("Fetching latest issue", "lottery_id", lotteryID)
	var issue models.LotteryIssue
	if err := db.DB.Where("lottery_id = ?", lotteryID).Order("issue_id desc").First(&issue).Error; err != nil {
		utils.Logger.Warn("No issue found for lottery", "lottery_id", lotteryID)
		return nil, utils.NewServiceError("issue not found", err)
	}
	utils.Logger.Info("Fetched latest issue", "issue_id", issue.IssueID)
	return &issue, nil
}

// GetExpiringIssues 获取即将到期的期号
func GetExpiringIssues() ([]models.LotteryIssue, error) {
	utils.Logger.Info("Fetching expiring issues")
	var issues []models.LotteryIssue
	now := time.Now()
	if err := db.DB.Where("draw_time > ? and draw_time <= ?", now, now.Add(24*time.Hour)).Find(&issues).Error; err != nil {
		utils.Logger.Error("Failed to fetch expiring issues", "error", err)
		return nil, utils.NewServiceError("failed to fetch expiring issues", err)
	}
	utils.Logger.Info("Fetched expiring issues", "count", len(issues))
	return issues, nil
}
