package services

import (
	"backend/blockchain"
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
func PurchaseTicket(ticket *models.LotteryTicket) error {
	utils.Logger.Info("Purchasing ticket", "issue_id", ticket.IssueID, "buyer", ticket.BuyerAddress)

	// Gas 估算函数
	estimateGas := func(opts *bind.TransactOpts) error {
		var issue models.LotteryIssue
		if err := db.DB.Where("issue_id = ?", ticket.IssueID).First(&issue).Error; err != nil {
			return err
		}

		var lottery models.Lottery
		if err := db.DB.Where("lottery_id = ?", issue.LotteryID).First(&lottery).Error; err != nil {
			return err
		}

		contract, err := connectLotteryContract(lottery.ContractAddress)
		if err != nil {
			return err
		}

		// 估算 RecordPlaceBet 的 Gas
		targets := parseBetContentV2bigIntSlice(ticket.BetContent)
		amount := big.NewInt(1)
		_, err = contract.RecordPlaceBet(opts, common.HexToAddress(ticket.BuyerAddress), amount, targets)
		return err
	}

	// 执行交易的函数
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
		amount := big.NewInt(1)
		ticket.PurchaseAmount = float64(new(big.Int).Mul(price, amount).Int64())
		ticket.PurchaseTime = time.Now()
		ticket.CreatedAt = time.Now()
		ticket.UpdatedAt = time.Now()

		contract, err := connectLotteryContract(lottery.ContractAddress)
		if err != nil {
			return common.Hash{}, err
		}

		targets := parseBetContentV2bigIntSlice(ticket.BetContent)
		if len(targets) != 3 {
			utils.Logger.Error("Invalid bet content", "content", ticket.BetContent)
			return common.Hash{}, utils.NewServiceError("invalid bet content: must contain exactly 3 numbers", nil)
		}

		tx, err := contract.RecordPlaceBet(blockchain.Auth, common.HexToAddress(ticket.BuyerAddress), amount, targets)
		if err != nil {
			utils.Logger.Error("Failed to record bet", "error", err)
			return tx.Hash(), utils.NewServiceError("failed to record bet", err)
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

	return blockchain.WithBlockchain(context.Background(), estimateGas, executeTx)
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
