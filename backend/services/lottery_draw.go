package services

import (
	"backend/blockchain"
	"backend/db"
	"backend/models"
	"backend/utils"
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

// DrawLottery 开奖
func DrawLottery(issueID string) error {
	utils.Logger.Info("Drawing lottery", "issue_id", issueID)

	// Gas 估算函数
	estimateGas := func(opts *bind.TransactOpts) error {
		var issue models.LotteryIssue
		if err := db.DB.Where("issue_id = ?", issueID).First(&issue).Error; err != nil {
			return err // 如果查询失败，返回错误
		}

		var lottery models.Lottery
		if err := db.DB.Where("lottery_id = ?", issue.LotteryID).First(&lottery).Error; err != nil {
			return err
		}

		contract, err := connectLotteryContract(lottery.ContractAddress)
		if err != nil {
			return err
		}

		// 估算 TransState 的 Gas
		_, err = contract.TransState(opts, uint8(2))
		if err != nil {
			return err
		}
		return nil
	}

	// 执行交易的函数
	executeTx := func() (common.Hash, error) {
		var issue models.LotteryIssue
		if err := db.DB.Where("issue_id = ?", issueID).First(&issue).Error; err != nil {
			utils.Logger.Warn("Issue not found", "issue_id", issueID)
			return common.Hash{}, utils.NewServiceError("issue not found", err)
		}

		var lottery models.Lottery
		if err := db.DB.Where("lottery_id = ?", issue.LotteryID).First(&lottery).Error; err != nil {
			utils.Logger.Warn("Lottery not found", "lottery_id", issue.LotteryID)
			return common.Hash{}, utils.NewServiceError("lottery not found", err)
		}

		contract, err := connectLotteryContract(lottery.ContractAddress)
		if err != nil {
			return common.Hash{}, err
		}

		currentState, err := contract.GetState(nil)
		if err != nil {
			utils.Logger.Error("Failed to get contract state", "error", err)
			return common.Hash{}, utils.NewServiceError("failed to get contract state", err)
		}
		if currentState != 2 {
			tx, err := contract.TransState(blockchain.Auth, uint8(2))
			if err != nil {
				utils.Logger.Error("Failed to set state to Rollout", "error", err)
				return tx.Hash(), utils.NewServiceError("failed to set state to Rollout", err)
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
		}

		if time.Now().Before(issue.DrawTime) {
			utils.Logger.Warn("Draw time not reached", "issue_id", issueID)
			return common.Hash{}, utils.NewServiceError("draw time not reached", nil)
		}

		// 为第二次交易准备新的 Nonce 和 Gas
		if _, err := blockchain.BlockchainMgr.GetNextNonce(context.Background()); err != nil {
			utils.Logger.Error("Failed to get next nonce for rollout", "error", err)
			return common.Hash{}, utils.NewServiceError("failed to get next nonce for rollout", err)
		}
		if _, err := blockchain.BlockchainMgr.GetCurrentGasPrice(context.Background()); err != nil {
			utils.Logger.Error("Failed to get gas price for rollout", "error", err)
			return common.Hash{}, utils.NewServiceError("failed to get gas price for rollout", err)
		}

		results := []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}
		tx, err := contract.RolloutCallback(blockchain.Auth, results)
		if err != nil {
			utils.Logger.Error("Failed to call rolloutCallback", "error", err)
			return tx.Hash(), utils.NewServiceError("failed to call rolloutCallback", err)
		}

		receipt, err := bind.WaitMined(context.Background(), blockchain.Client, tx)
		if err != nil {
			utils.Logger.Error("Transaction failed", "tx_hash", tx.Hash().Hex(), "error", err)
			return tx.Hash(), utils.NewServiceError("transaction failed", err)
		}
		if receipt.Status != 1 {
			utils.Logger.Error("Rollout transaction failed", "tx_hash", tx.Hash().Hex(), "status", receipt.Status)
			return tx.Hash(), utils.NewServiceError(fmt.Sprintf("rollout transaction failed with status: %d", receipt.Status), nil)
		}

		issue.WinningNumbers = fmt.Sprintf("%d,%d,%d", results[0], results[1], results[2])
		issue.DrawTxHash = tx.Hash().Hex()
		issue.UpdatedAt = time.Now()
		if err := db.DB.Save(&issue).Error; err != nil {
			utils.Logger.Error("Failed to update issue", "error", err)
			return common.Hash{}, utils.NewServiceError("failed to update issue", err)
		}

		var tickets []models.LotteryTicket
		if err := db.DB.Where("issue_id = ?", issueID).Find(&tickets).Error; err != nil {
			utils.Logger.Error("Failed to get tickets", "error", err)
			return common.Hash{}, utils.NewServiceError("failed to get tickets", err)
		}
		for _, ticket := range tickets {
			ticketNumbers := parseBetContent(ticket.BetContent)
			if len(ticketNumbers) == 3 && ticketNumbers[0].Cmp(results[0]) == 0 &&
				ticketNumbers[1].Cmp(results[1]) == 0 && ticketNumbers[2].Cmp(results[2]) == 0 {
				winner := models.Winner{
					WinnerID:    uuid.NewString(),
					IssueID:     issueID,
					TicketID:    ticket.TicketID,
					Address:     ticket.BuyerAddress,
					PrizeLevel:  "First Prize",
					PrizeAmount: ticket.PurchaseAmount,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				if err := db.DB.Create(&winner).Error; err != nil {
					utils.Logger.Error("Failed to save winner", "ticket_id", ticket.TicketID, "error", err)
					return common.Hash{}, utils.NewServiceError("failed to save winner", err)
				}
			}
		}

		utils.Logger.Info("Lottery drawn successfully", "issue_id", issueID)
		return tx.Hash(), nil
	}

	return blockchain.WithBlockchain(context.Background(), estimateGas, executeTx)
}

// GetDrawnLotteryByIssueID 根据期号ID获取已开奖的彩票信息
func GetDrawnLotteryByIssueID(issueID string) (*models.LotteryIssue, error) {
	utils.Logger.Info("Fetching drawn lottery", "issue_id", issueID)
	var issue models.LotteryIssue
	if err := db.DB.Where("issue_id = ?", issueID).First(&issue).Error; err != nil {
		utils.Logger.Warn("Drawn lottery not found", "issue_id", issueID)
		return nil, utils.NewServiceError("drawn lottery not found", err)
	}
	utils.Logger.Info("Fetched drawn lottery", "issue_id", issueID)
	return &issue, nil
}
