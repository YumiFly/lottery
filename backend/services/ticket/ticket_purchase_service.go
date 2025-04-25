package ticket

import (
	"context"
	"math/big"
	"strings"
	"time"

	"backend/blockchain"
	"backend/config"
	"backend/models"
	"backend/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// PurchaseTicketParams defines the parameters for purchasing a ticket
type PurchaseTicketParams struct {
	TicketID       string
	IssueID        string
	BuyerAddress   string
	PurchaseAmount uint64
	BetContent     string
}

// TicketService encapsulates ticket purchasing business logic
type TicketPurchaseService struct {
	db *gorm.DB
}

// NewTicketService creates a new TicketService instance
func NewTicketPurchaseService(db *gorm.DB) *TicketPurchaseService {
	return &TicketPurchaseService{db: db}
}

// validatePurchaseTicketParams validates the parameters for purchasing a ticket
func (s *TicketPurchaseService) validatePurchaseTicketParams(ctx context.Context, params PurchaseTicketParams) error {
	// Validate issue_id exists and sale is active
	var issue models.LotteryIssue
	if err := s.db.WithContext(ctx).
		Where("issue_id = ? AND status = ?", params.IssueID, models.IssueStatusPending).
		First(&issue).Error; err != nil {
		utils.Logger.Error("Failed to find active issue", "issue_id", params.IssueID, "error", err)
		return utils.NewBadRequestError("Invalid or inactive issue_id", nil)
	}

	// Validate lottery exists
	var lottery models.Lottery
	if err := s.db.WithContext(ctx).
		Where("lottery_id = ?", issue.LotteryID).
		First(&lottery).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Logger.Warn("Lottery not found", "lottery_id", issue.LotteryID)
			return utils.NewBadRequestError("Lottery not found", nil)
		}
		return utils.NewInternalError("Failed to check lottery ID", errors.Wrap(err, "database error"))
	}

	// Validate ticket supply
	var totalTickets uint64
	if err := s.db.WithContext(ctx).
		Model(&models.LotteryTicket{}).
		Where("issue_id = ?", params.IssueID).
		Select("COALESCE(SUM(purchase_amount), 0)").
		Scan(&totalTickets).Error; err != nil {
		utils.Logger.Error("Failed to calculate total tickets", "issue_id", params.IssueID, "error", err)
		return utils.NewInternalError("Failed to validate ticket supply", err)
	}
	if totalTickets+params.PurchaseAmount > uint64(lottery.TicketSupply) {
		return utils.NewBadRequestError("Purchase amount exceeds available ticket supply", nil)
	}

	if time.Now().After(issue.SaleEndTime) {
		utils.Logger.Warn("Sale has ended", "issue_id", params.IssueID)
		return utils.NewBadRequestError("Sale has ended", nil)
	}

	// Validate buyer_address
	if !common.IsHexAddress(params.BuyerAddress) {
		return utils.NewBadRequestError("Invalid buyer address", nil)
	}

	// Validate purchase_amount
	if params.PurchaseAmount <= 0 {
		return utils.NewBadRequestError("Purchase amount must be positive", nil)
	}

	// Validate bet_content
	if len(params.BetContent) == 0 || len(params.BetContent) > 100 {
		return utils.NewBadRequestError("Bet content must be between 1 and 100 characters", nil)
	}
	targets := strings.Split(params.BetContent, ",")
	if len(targets) != 3 {
		utils.Logger.Warn("Invalid bet content", "content", params.BetContent)
		return utils.NewBadRequestError("Bet content must contain exactly 3 numbers", nil)
	}
	for _, target := range targets {
		if _, ok := new(big.Int).SetString(strings.TrimSpace(target), 10); !ok {
			return utils.NewBadRequestError("Bet content contains invalid numbers", nil)
		}
	}
	// Validate ticket_id
	if len(params.TicketID) == 0 || len(params.TicketID) > 36 {
		return utils.NewBadRequestError("Invalid ticket ID", nil)
	}

	return nil
}

// PurchaseTicket purchases a lottery ticket
//
// Parameters:
//   - ctx: Request context
//   - params: Purchase parameters, including issue_id, buyer_address, etc.
//
// Returns:
//   - *LotteryTicket: The purchased ticket record
//   - common.Hash: Blockchain transaction hash
//   - error: Purchase error or invalid parameters
func (s *TicketPurchaseService) PurchaseTicket(ctx context.Context, params PurchaseTicketParams) (*models.LotteryTicket, common.Hash, error) {
	// Validate parameters
	if err := s.validatePurchaseTicketParams(ctx, params); err != nil {
		return nil, common.Hash{}, err
	}

	// Execute blockchain transaction
	data := []byte{}
	ticket := models.LotteryTicket{}

	executeTx := func() (common.Hash, error) {
		// Fetch issue and lottery
		var issue models.LotteryIssue
		if err := s.db.WithContext(ctx).
			Where("issue_id = ?", params.IssueID).
			First(&issue).Error; err != nil {
			return common.Hash{}, utils.NewBadRequestError("Issue not found", err)
		}

		var lottery models.Lottery
		if err := s.db.WithContext(ctx).
			Where("lottery_id = ?", issue.LotteryID).
			First(&lottery).Error; err != nil {
			return common.Hash{}, utils.NewBadRequestError("Lottery not found", err)
		}

		// Convert ticket price to big.Int (assume TicketPrice is in ETH, convert to wei)
		price, ok := new(big.Int).SetString(big.NewFloat(lottery.TicketPrice).Text('f', 0), 10)
		if !ok {
			utils.Logger.Error("Invalid ticket price", "price", lottery.TicketPrice)
			return common.Hash{}, utils.NewServiceError("invalid ticket price in lottery", nil)
		}
		totalPrice := new(big.Int).SetInt64(int64(ticket.PurchaseAmount)) // 获取总价
		// 计算数量
		amount := new(big.Int).Div(totalPrice, price)

		// Construct ticket record
		ticket = models.LotteryTicket{
			TicketID:       params.TicketID,
			IssueID:        params.IssueID,
			BuyerAddress:   params.BuyerAddress,
			PurchaseAmount: float64(params.PurchaseAmount), // Store as number of tickets
			BetContent:     params.BetContent,
			PurchaseTime:   time.Now(),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Log purchase attempt
		utils.Logger.Info("Purchasing ticket",
			"ticket_id", ticket.TicketID,
			"issue_id", ticket.IssueID,
			"buyer", ticket.BuyerAddress,
			"amount", amount.String(),
			"total_price", totalPrice.String())

		// Parse bet content
		targets := utils.ParseBetContentV2bigIntSlice(ticket.BetContent)
		if len(targets) != 3 {
			utils.Logger.Error("Invalid bet content", "content", ticket.BetContent)
			return common.Hash{}, utils.NewBadRequestError("Invalid bet content: must contain exactly 3 numbers", nil)
		}

		// Connect to token contract
		tokenContract, err := blockchain.ConnectTokenContract(config.AppConfig.TokenContractAddress)
		if err != nil {
			utils.Logger.Error("Failed to connect to LOTToken contract", "error", err)
			return common.Hash{}, utils.NewInternalError("Failed to connect to LOTToken contract", err)
		}

		// Call Buy function on token contract
		tx, err := tokenContract.Buy(
			blockchain.Auth,
			common.HexToAddress(lottery.ContractAddress),
			amount,
			targets,
		)
		if err != nil {
			utils.Logger.Error("Failed to buy ticket",
				"error", err,
				"amount", amount.String(),
				"total_price", totalPrice.String())
			if tx != nil {
				return tx.Hash(), utils.NewInternalError("Failed to buy ticket", err)
			}
			return common.Hash{}, utils.NewInternalError("Failed to buy ticket", err)
		}

		// Wait for transaction confirmation
		receipt, err := bind.WaitMined(ctx, blockchain.Client, tx)
		if err != nil {
			utils.Logger.Error("Transaction failed", "tx_hash", tx.Hash().Hex(), "error", err)
			return tx.Hash(), utils.NewInternalError("Transaction failed", err)
		}
		if receipt.Status != 1 {
			utils.Logger.Error("Transaction failed", "tx_hash", tx.Hash().Hex(), "status", receipt.Status)
			return tx.Hash(), utils.NewInternalError("Transaction failed", nil)
		}

		// Update ticket and issue
		ticket.TransactionHash = tx.Hash().Hex()
		issue.PrizePool += float64(totalPrice.Int64()) / 1e18 // Convert wei to ETH for prize pool

		// Save to database within a transaction
		err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Save(&issue).Error; err != nil {
				utils.Logger.Error("Failed to update issue prize pool", "error", err)
				return utils.NewInternalError("Failed to update issue prize pool", err)
			}
			if err := tx.Create(&ticket).Error; err != nil {
				utils.Logger.Error("Failed to save ticket to database", "error", err)
				return utils.NewInternalError("Failed to save ticket to database", err)
			}
			return nil
		})
		if err != nil {
			return common.Hash{}, err
		}

		utils.Logger.Info("Ticket purchased successfully",
			"ticket_id", ticket.TicketID,
			"tx_hash", tx.Hash().Hex())
		return tx.Hash(), nil
	}

	txHash, err := blockchain.WithBlockchain(ctx, data, executeTx)
	if err != nil {
		return nil, common.Hash{}, err
	}

	return &ticket, txHash, nil
}
