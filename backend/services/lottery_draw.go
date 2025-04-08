package services

import (
	"backend/blockchain"
	blockLottery "backend/blockchain/lottery"
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
			return common.Hash{}, utils.NewServiceError("failed to connect to lottery contract", err)
		}

		// 检查当前状态并设置为 Rollout (状态 2)
		currentState, err := contract.GetState(nil)
		if err != nil {
			utils.Logger.Error("Failed to get contract state", "error", err)
			return common.Hash{}, utils.NewServiceError("failed to get contract state", err)
		}
		utils.Logger.Info("Current contract state", "state", currentState) // 添加日志记录
		var txHash common.Hash
		if currentState != 2 {
			tx, err := contract.TransState(blockchain.Auth, uint8(2))
			if err != nil {
				utils.Logger.Error("Failed to set state to Rollout", "error", err)
				if tx != nil {
					return tx.Hash(), utils.NewServiceError("failed to set state to Rollout", err)
				}
				return common.Hash{}, utils.NewServiceError("failed to set state to Rollout", err)
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
			txHash = tx.Hash()
			utils.Logger.Info("State set to Rollout", "tx_hash", txHash.Hex())

			// 获取初始 nonce
			nonce, err := blockchain.BlockchainMgr.GetNextNonce(context.Background())
			if err != nil {
				return common.Hash{}, utils.NewServiceError("failed to get next nonce", err)
			}
			blockchain.Auth.Nonce = big.NewInt(int64(nonce))
		}

		rolloutContract, err := connectRolloutContract(lottery.RolloutContractAddress)
		if err != nil {
			return common.Hash{}, utils.NewServiceError("failed to connect to rollout contract", err)
		}

		tx, err := rolloutContract.RolloutCall(blockchain.Auth, common.HexToAddress(lottery.ContractAddress))
		if err != nil {
			utils.Logger.Error("Failed to call rolloutCall", "error", err)
			if tx != nil {
				utils.Logger.Info("RolloutCall tx hash", "tx_hash", tx.Hash().Hex())
				return tx.Hash(), utils.NewServiceError("failed to call rolloutCall", err)
			}
			return common.Hash{}, utils.NewServiceError("failed to call rolloutCall", err)
		}
		utils.Logger.Info("RolloutCall transaction submitted", "tx_hash", tx.Hash().Hex(), "nonce", blockchain.Auth.Nonce)

		receipt, err := bind.WaitMined(context.Background(), blockchain.Client, tx)
		if err != nil {
			utils.Logger.Error("Transaction failed", "tx_hash", tx.Hash().Hex(), "error", err)
			return tx.Hash(), utils.NewServiceError("transaction failed", err)
		}
		if receipt.Status != 1 {
			utils.Logger.Error("Rollout transaction failed", "tx_hash", tx.Hash().Hex(), "status", receipt.Status)
			return tx.Hash(), utils.NewServiceError(fmt.Sprintf("rollout transaction failed with status: %d", receipt.Status), nil)
		}
		utils.Logger.Info("RolloutCall transaction confirmed", "tx_hash", tx.Hash().Hex(), "gas_used", receipt.GasUsed)
		startBlock := receipt.BlockNumber

		// 检查交易日志
		for _, log := range receipt.Logs {
			utils.Logger.Info("Transaction log", "topics", log.Topics, "data", log.Data)
		}

		// 检查 DiceRolled 事件
		filterer, err := blockLottery.NewSimpleRolloutFilterer(common.HexToAddress(lottery.RolloutContractAddress), blockchain.Client)
		if err != nil {
			return common.Hash{}, utils.NewServiceError("failed to create filterer", err)
		}
		diceRolledOpts := &bind.FilterOpts{
			Start:   startBlock.Uint64(),
			Context: context.Background(),
		}
		diceRolledIter, err := filterer.FilterDiceRolled(diceRolledOpts, nil, nil)
		if err != nil {
			return common.Hash{}, utils.NewServiceError("failed to filter dice rolled event", err)
		}
		for diceRolledIter.Next() {
			utils.Logger.Info("DiceRolled event captured", "request_id", diceRolledIter.Event.RequestId, "epoch", diceRolledIter.Event.Epoch)
		}
		diceRolledIter.Close()

		// 等待 DiceLanded 事件
		timeout := time.After(300 * time.Second)
		var results []*big.Int
		elapsed := 0 * time.Second
		for {
			select {
			case <-timeout:
				utils.Logger.Error("Timeout waiting for DiceLanded event", "issue_id", issueID, "elapsed", elapsed)
				return tx.Hash(), utils.NewServiceError("timeout waiting for dice landed event", nil)
			default:
				opts := &bind.FilterOpts{
					Start:   startBlock.Uint64(),
					Context: context.Background(),
				}
				iter, err := filterer.FilterDiceLanded(opts, nil, nil)
				if err != nil {
					return common.Hash{}, utils.NewServiceError("failed to filter dice landed event", err)
				}
				for iter.Next() {
					results = iter.Event.Results // 确保 iter.Event.Results 被正确解析为 []*big.Int
					utils.Logger.Info("DiceLanded event captured", "results", results)
				}
				iter.Close()

				if len(results) > 0 {
					goto processResults
				}
				time.Sleep(1 * time.Second)
				elapsed += 1 * time.Second
				utils.Logger.Debug("Waiting for DiceLanded event", "elapsed", elapsed)
			}
		}

	processResults:
		// 更新期号信息
		issue.WinningNumbers = fmt.Sprintf("%d,%d,%d", results[0], results[1], results[2])
		issue.DrawTxHash = tx.Hash().Hex()
		issue.UpdatedAt = time.Now()
		if err := db.DB.Save(&issue).Error; err != nil {
			utils.Logger.Error("Failed to update issue", "error", err)
			return common.Hash{}, utils.NewServiceError("failed to update issue", err)
		}

		// 检查中奖者
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

	data := []byte{}
	return blockchain.WithBlockchain(context.Background(), data, executeTx)
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
