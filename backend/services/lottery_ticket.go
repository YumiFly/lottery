package services

import (
	"backend/blockchain"
	"backend/config"
	"backend/db"
	"backend/models"
	"backend/utils"
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// PurchaseTicket 购买彩票
// PurchaseTicket 购买彩票
func PurchaseTicket(ticket *models.LotteryTicket) (common.Hash, error) {
	utils.Logger.Info("Purchasing ticket", "issue_id", ticket.IssueID, "buyer", ticket.BuyerAddress)

	executeTx := func() (common.Hash, error) {
		var issue models.LotteryIssue
		if err := db.DB.Where("issue_id = ?", ticket.IssueID).First(&issue).Error; err != nil {
			utils.Logger.Warn("Issue not found", "issue_id", ticket.IssueID)
			return common.Hash{}, utils.NewServiceError("issue not found", err)
		}

		if time.Now().After(issue.SaleEndTime) {
			utils.Logger.Warn("Sale has ended", "issue_id", ticket.IssueID)
			return common.Hash{}, utils.NewServiceError("sale has ended", nil)
		}

		var lottery models.Lottery
		if err := db.DB.Where("lottery_id = ?", issue.LotteryID).First(&lottery).Error; err != nil {
			utils.Logger.Warn("Lottery not found", "lottery_id", issue.LotteryID)
			return common.Hash{}, utils.NewServiceError("lottery not found", err)
		}

		price, ok := new(big.Int).SetString(big.NewFloat(lottery.TicketPrice).Text('f', 0), 10)
		if !ok {
			utils.Logger.Error("Invalid ticket price", "price", lottery.TicketPrice)
			return common.Hash{}, utils.NewServiceError("invalid ticket price in lottery", nil)
		}

		totalPrice := new(big.Int).SetInt64(int64(ticket.PurchaseAmount)) // 获取总价

		// 计算数量
		amount := new(big.Int).Div(totalPrice, price)

		// 更新 ticket.PurchaseAmount
		ticket.PurchaseAmount = float64(new(big.Int).Mul(price, amount).Int64())
		ticket.PurchaseTime = time.Now()
		ticket.CreatedAt = time.Now()
		ticket.UpdatedAt = time.Now()

		targets := parseBetContentV2bigIntSlice(ticket.BetContent)
		if len(targets) != 3 {
			utils.Logger.Error("Invalid bet content", "content", ticket.BetContent)
			return common.Hash{}, utils.NewServiceError("invalid bet content: must contain exactly 3 numbers", nil)
		}

		// 获取 LOTToken 合约实例
		tokenContract, err := connectTokenContract(config.AppConfig.TokenContractAddress)
		if err != nil {
			utils.Logger.Error("Failed to connect to LOTToken contract", "error", err)
			return common.Hash{}, utils.NewServiceError("failed to connect to LOTToken contract", err)
		}

		// 调用 LOTToken 合约的 buy 函数
		tx, err := tokenContract.Buy(blockchain.Auth, common.HexToAddress(lottery.ContractAddress), amount, targets)
		if err != nil {
			utils.Logger.Error("Failed to buy ticket", "error", err)
			if tx != nil {
				return tx.Hash(), utils.NewServiceError("failed to buy ticket", err)
			}
			return common.Hash{}, utils.NewServiceError("failed to buy ticket", err)
		}

		receipt, err := bind.WaitMined(context.Background(), blockchain.Client, tx)
		if err != nil {
			utils.Logger.Error("Transaction failed", "tx_hash", tx.Hash().Hex(), "error", err)
			return tx.Hash(), utils.NewServiceError("transaction failed", err)
		}
		if receipt.Status != 1 {
			utils.Logger.Error("Transaction failed", "tx_hash", tx.Hash().Hex(), "status", receipt.Status)
			return tx.Hash(), utils.NewServiceError("transaction failed with status: "+string(rune(receipt.Status)), nil)
		}

		utils.Logger.Info("Ticket purchased successfully", "tx_hash", tx.Hash().Hex())

		ticket.TransactionHash = tx.Hash().Hex()
		issue.PrizePool = issue.PrizePool + ticket.PurchaseAmount
		if err := db.DB.Save(&issue).Error; err != nil {
			utils.Logger.Error("Failed to update issue prize pool", "error", err)
			return common.Hash{}, utils.NewServiceError("failed to update issue prize pool", err)
		}

		if err := db.DB.Create(ticket).Error; err != nil {
			utils.Logger.Error("Failed to save ticket to database", "error", err)
			return common.Hash{}, utils.NewServiceError("failed to save ticket to database", err)
		}

		utils.Logger.Info("Ticket purchased successfully", "ticket_id", ticket.TicketID)
		return tx.Hash(), nil
	}

	// 使用空的 data 调用 WithBlockchain，Gas 估算依赖内部逻辑
	data := []byte{}
	return blockchain.WithBlockchain(context.Background(), data, executeTx)
}

// GetPurchasedTicketsByCustomerAddress 根据客户地址获取已购彩票
func GetPurchasedTicketsByCustomerAddress(customerAddress string) ([]models.LotteryTicket, error) {
	utils.Logger.Info("Fetching tickets by customer address", "address", customerAddress)
	var tickets []models.LotteryTicket
	if err := db.DB.Where("buyer_address = ?", customerAddress).Find(&tickets).Error; err != nil {
		utils.Logger.Error("Failed to fetch tickets", "address", customerAddress, "error", err)
		return nil, utils.NewServiceError("failed to fetch tickets", err)
	}
	utils.Logger.Info("Fetched tickets", "address", customerAddress, "count", len(tickets))
	return tickets, nil
}
