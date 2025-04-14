package services

import (
	"backend/blockchain"
	lotteryBlockchain "backend/blockchain/lottery"
	"backend/db"
	"backend/models"
	"backend/utils"
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Constants for configuration
const (
	LotteryResultsTimeout = 60 * time.Second // Timeout for waiting lottery results
	LotteryResultsRetries = 3                // Number of retries for event subscription
	ExpectedResultCount   = 3                // Expected number of lottery results
	TicketBatchSize       = 1000             // Batch size for ticket queries
)

// LotteryService encapsulates lottery-related operations
type LotteryService struct {
	client *ethclient.Client
	auth   *bind.TransactOpts
	db     *gorm.DB
}

// NewLotteryService creates a new LotteryService instance
func NewLotteryService(client *ethclient.Client, auth *bind.TransactOpts, db *gorm.DB) *LotteryService {
	return &LotteryService{
		client: client,
		auth:   auth,
		db:     db,
	}
}

// DrawLotteryAsync initiates an asynchronous lottery draw
// It checks the issue status to prevent duplicate draws and runs the draw in a goroutine
func (s *LotteryService) DrawLotteryAsync(issueID string) error {
	utils.Logger.Info("Starting asynchronous lottery draw", "issue_id", issueID)

	// Check if the issue has already been drawn
	var issue models.LotteryIssue
	if err := s.db.Where("issue_id = ?", issueID).Select("status").First(&issue).Error; err != nil {
		utils.Logger.Warn("Lottery issue not found", "issue_id", issueID)
		return utils.NewServiceError("lottery issue not found", err)
	}
	if issue.Status == models.IssueStatusDrawn {
		utils.Logger.Info("Lottery already drawn", "issue_id", issueID)
		return nil
	}

	// Run draw in background
	go func() {
		if err := s.executeLotteryDraw(issueID); err != nil {
			utils.Logger.Error("Failed to complete async lottery draw", "issue_id", issueID, "error", err)
		}
	}()

	return nil
}

// executeLotteryDraw performs the lottery draw for the specified issue
// It fetches data, sets contract state, executes the draw, and processes results
func (s *LotteryService) executeLotteryDraw(issueID string) error {
	// Fetch issue and lottery data
	lottery, err := s.fetchLotteryData(issueID)
	if err != nil {
		return err
	}
	// Validate contract addresses
	if lottery.ContractAddress == "" || lottery.RolloutContractAddress == "" {
		return utils.NewServiceError("invalid contract addresses", nil)
	}
	// Initialize LotteryManager contract
	contract, err := lotteryBlockchain.NewLotteryManager(common.HexToAddress(lottery.ContractAddress), s.client)
	if err != nil {
		return utils.NewServiceError("failed to connect to lottery contract", err)
	}
	// Set contract state to Rollout if needed
	if err := s.setContractState(contract, uint8(models.ContractStateRollout)); err != nil {
		return err
	}

	// 更新 issue 状态
	if err := s.db.Model(&models.LotteryIssue{}).Where("issue_id = ?", issueID).Update("status", models.IssueStatusDrawing).Error; err != nil {
		return utils.NewServiceError("failed to update lottery issue status", err)
	}

	// Subscribe to LotteryResults event before triggering rollout
	resultsChan := make(chan []*big.Int)
	errChan := make(chan error)
	go func() {
		results, err := s.subscribeToLotteryResults(contract)
		if err != nil {
			errChan <- utils.NewServiceError("failed to subscribe to LotteryResults event", err)
			return
		}
		resultsChan <- results
	}()

	// Execute the rollout call
	tx, err := s.executeRollout(lottery)
	if err != nil {
		return err
	}

	// Wait for and process results
	return s.waitAndProcessResults(issueID, contract, tx, resultsChan, errChan)
}

// fetchLotteryData retrieves lottery issue and associated lottery data
func (s *LotteryService) fetchLotteryData(issueID string) (*models.Lottery, error) {
	var lottery models.Lottery
	var issue models.LotteryIssue
	if err := s.db.Where("issue_id = ?", issueID).First(&issue).Error; err != nil {
		return nil, utils.NewServiceError("failed to fetch lottery issue data", err)
	}
	if err := s.db.Where("lottery_id = ?", issue.LotteryID).First(&lottery).Error; err != nil {
		return nil, utils.NewServiceError("failed to fetch lottery data", err)
	}
	return &lottery, nil
}

// setContractState sets the lottery contract state to the target state if needed
func (s *LotteryService) setContractState(contract *lotteryBlockchain.LotteryManager, targetState uint8) error {
	// Get current contract state
	state, err := contract.GetState(nil)
	if err != nil {
		return utils.NewServiceError("failed to get contract state", err)
	}
	utils.Logger.Info("Current contract state", "state", state)

	if state == targetState {
		return nil
	}

	// Get nonce for transaction
	nonce, err := blockchain.BlockchainMgr.GetNextNonce(context.Background())
	if err != nil {
		return utils.NewServiceError("failed to get next nonce", err)
	}
	s.auth.Nonce = big.NewInt(int64(nonce))

	// Set state to target
	utils.Logger.Info("Setting contract state", "target_state", targetState)
	tx, err := contract.TransState(s.auth, targetState)
	if err != nil {
		return utils.NewServiceError(fmt.Sprintf("failed to set state to %d", targetState), err)
	}

	// Wait for transaction confirmation
	receipt, err := bind.WaitMined(context.Background(), s.client, tx)
	if err != nil || receipt.Status != 1 {
		return utils.NewServiceError(fmt.Sprintf("failed to confirm state transition to %d", targetState), err)
	}

	// Update nonce
	newNonce := blockchain.BlockchainMgr.GetNextNonceForFunc(nonce)
	s.auth.Nonce = big.NewInt(int64(newNonce))

	return nil
}

// executeRollout calls the rollout contract to perform the lottery draw
func (s *LotteryService) executeRollout(lottery *models.Lottery) (*types.Transaction, error) {
	// Initialize Rollout contract
	rolloutContract, err := lotteryBlockchain.NewSimpleRollout(common.HexToAddress(lottery.RolloutContractAddress), s.client)
	if err != nil {
		return nil, utils.NewServiceError("failed to initialize Rollout contract", err)
	}

	// Get nonce
	nonce, err := blockchain.BlockchainMgr.GetNextNonce(context.Background())
	if err != nil {
		return nil, utils.NewServiceError("failed to get next nonce", err)
	}
	s.auth.Nonce = big.NewInt(int64(nonce))

	// Call rollout
	utils.Logger.Info("Calling rolloutCall", "rollout_contract", lottery.RolloutContractAddress, "lottery_manager", lottery.ContractAddress)
	tx, err := rolloutContract.RolloutCall(s.auth, common.HexToAddress(lottery.ContractAddress))
	if err != nil {
		return nil, utils.NewServiceError("failed to call rolloutCall", err)
	}

	// Wait for transaction confirmation
	receipt, err := bind.WaitMined(context.Background(), s.client, tx)
	if err != nil || receipt.Status != 1 {
		return nil, utils.NewServiceError("failed to confirm rolloutCall", err)
	}
	utils.Logger.Info("rolloutCall transaction confirmed", "tx_hash", tx.Hash().Hex(), "block_number", receipt.BlockNumber)

	// Update nonce
	newNonce := blockchain.BlockchainMgr.GetNextNonceForFunc(nonce)
	s.auth.Nonce = big.NewInt(int64(newNonce))

	return tx, nil
}

// waitAndProcessResults waits for lottery results and processes them
func (s *LotteryService) waitAndProcessResults(issueID string, contract *lotteryBlockchain.LotteryManager, tx *types.Transaction, resultsChan chan []*big.Int, errChan chan error) error {
	// Wait for results
	select {
	case results := <-resultsChan:
		// Verify contract state after draw
		state, err := contract.GetState(nil)
		if err != nil {
			return utils.NewServiceError("failed to get contract state after draw", err)
		}
		if state != uint8(models.ContractStateReady) {
			return utils.NewServiceError(fmt.Sprintf("contract state not Ready after draw, current state: %d", state), nil)
		}
		utils.Logger.Info("Contract state verified after draw", "state", state)

		// Process results
		if err := s.recordLotteryResults(issueID, results, tx.Hash()); err != nil {
			return err
		}
		utils.Logger.Info("Lottery draw completed successfully", "issue_id", issueID, "tx_hash", tx.Hash().Hex())
		return nil
	case err := <-errChan:
		// On error, attempt to query historical logs as a fallback
		utils.Logger.Warn("Subscription failed, attempting to query historical logs", "issue_id", issueID, "error", err)
		results, err := s.queryHistoricalResults(contract, tx.Hash())
		if err != nil {
			return utils.NewServiceError("failed to recover results from historical logs", err)
		}
		// Process historical results
		if err := s.recordLotteryResults(issueID, results, tx.Hash()); err != nil {
			return err
		}
		utils.Logger.Info("Lottery draw completed successfully using historical logs", "issue_id", issueID, "tx_hash", tx.Hash().Hex())
		return nil
	case <-time.After(LotteryResultsTimeout):
		// On timeout, attempt to query historical logs
		utils.Logger.Warn("Timeout waiting for LotteryResults event, querying historical logs", "issue_id", issueID)
		results, err := s.queryHistoricalResults(contract, tx.Hash())
		if err != nil {
			return utils.NewServiceError("failed to recover results from historical logs after timeout", err)
		}
		// Process historical results
		if err := s.recordLotteryResults(issueID, results, tx.Hash()); err != nil {
			return err
		}
		utils.Logger.Info("Lottery draw completed successfully using historical logs", "issue_id", issueID, "tx_hash", tx.Hash().Hex())
		return nil
	}
}

// subscribeToLotteryResults subscribes to the LotteryResults event and retries on failure
func (s *LotteryService) subscribeToLotteryResults(contract *lotteryBlockchain.LotteryManager) ([]*big.Int, error) {
	logs := make(chan *lotteryBlockchain.LotteryManagerLotteryResults)
	opts := &bind.WatchOpts{Context: context.Background()}

	for attempt := 1; attempt <= LotteryResultsRetries; attempt++ {
		utils.Logger.Info("Starting to subscribe to LotteryResults event", "attempt", attempt)
		sub, err := contract.WatchLotteryResults(opts, logs)
		if err != nil {
			utils.Logger.Warn("Failed to subscribe, retrying", "attempt", attempt, "error", err)
			if attempt == LotteryResultsRetries {
				return nil, fmt.Errorf("failed to subscribe to LotteryResults event after %d attempts: %v", LotteryResultsRetries, err)
			}
			time.Sleep(time.Second * time.Duration(attempt))
			continue
		}

		utils.Logger.Info("Successfully subscribed, waiting for event", "attempt", attempt)
		select {
		case event := <-logs:
			sub.Unsubscribe()
			if len(event.Results) != ExpectedResultCount {
				return nil, fmt.Errorf("expected %d results, got %d", ExpectedResultCount, len(event.Results))
			}
			// Validate results
			for i, result := range event.Results {
				if result == nil || result.Cmp(big.NewInt(0)) <= 0 {
					return nil, fmt.Errorf("invalid result at index %d: %v", i, result)
				}
			}
			utils.Logger.Info("Received LotteryResults event", "results", event.Results, "epoch", event.Epoch, "timestamp", event.Timestamp)
			return event.Results, nil
		case err := <-sub.Err():
			sub.Unsubscribe()
			utils.Logger.Warn("Subscription error, retrying", "attempt", attempt, "error", err)
			if attempt == LotteryResultsRetries {
				return nil, fmt.Errorf("subscription error: %v", err)
			}
		case <-time.After(LotteryResultsTimeout):
			sub.Unsubscribe()
			utils.Logger.Warn("Timeout waiting for LotteryResults event", "attempt", attempt)
			if attempt == LotteryResultsRetries {
				return nil, fmt.Errorf("timeout waiting for LotteryResults event after %d attempts", LotteryResultsRetries)
			}
		}
	}
	return nil, fmt.Errorf("failed to subscribe, exceeded retry attempts")
}

// queryHistoricalResults queries historical logs for LotteryResults events
func (s *LotteryService) queryHistoricalResults(contract *lotteryBlockchain.LotteryManager, txHash common.Hash) ([]*big.Int, error) {
	// Query logs from the block of the transaction
	_, _, err := s.client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transaction: %v", err)
	}
	receipt, err := s.client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transaction receipt: %v", err)
	}

	// Set filter options to query logs from the transaction block
	opts := &bind.FilterOpts{
		Context: context.Background(),
		Start:   receipt.BlockNumber.Uint64(),
		End:     nil, // Query up to the latest block
	}

	// Filter LotteryResults events
	iterator, err := contract.FilterLotteryResults(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to filter LotteryResults events: %v", err)
	}
	defer iterator.Close()

	for iterator.Next() {
		event := iterator.Event
		if len(event.Results) != ExpectedResultCount {
			continue
		}
		// Validate results
		for i, result := range event.Results {
			if result == nil || result.Cmp(big.NewInt(0)) <= 0 {
				return nil, fmt.Errorf("invalid historical result at index %d: %v", i, result)
			}
		}
		utils.Logger.Info("Found LotteryResults event in historical logs", "results", event.Results, "epoch", event.Epoch, "timestamp", event.Timestamp)
		return event.Results, nil
	}

	return nil, fmt.Errorf("no valid LotteryResults event found in historical logs")
}

// recordLotteryResults updates the issue and saves winners in a transaction
func (s *LotteryService) recordLotteryResults(issueID string, results []*big.Int, txHash common.Hash) error {
	// Begin transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil || tx.Error != nil {
			tx.Rollback()
		}
	}()

	// Fetch issue
	var issue models.LotteryIssue
	if err := tx.Where("issue_id = ?", issueID).First(&issue).Error; err != nil {
		utils.Logger.Error("Failed to find issue", "issue_id", issueID, "error", err)
		return utils.NewServiceError("failed to find issue", err)
	}

	// Update issue
	issue.WinningNumbers = fmt.Sprintf("%d,%d,%d", results[0], results[1], results[2])
	issue.DrawTxHash = txHash.Hex()
	issue.Status = models.IssueStatusDrawn
	issue.DrawTime = time.Now()
	issue.UpdatedAt = time.Now()
	if err := tx.Save(&issue).Error; err != nil {
		utils.Logger.Error("Failed to update issue", "issue_id", issueID, "error", err)
		return utils.NewServiceError("failed to update issue", err)
	}
	utils.Logger.Info("Issue updated successfully", "issue_id", issueID, "winning_numbers", issue.WinningNumbers)

	// Fetch and save winners
	winners, err := s.getWinnersFromChain(issueID, results)
	if err != nil {
		return utils.NewServiceError("failed to get winners from chain", err)
	}
	for _, winner := range winners {
		if err := tx.Create(&winner).Error; err != nil {
			utils.Logger.Error("Failed to save winner", "ticket_id", winner.TicketID, "error", err)
			return utils.NewServiceError("failed to save winner", err)
		}
	}
	utils.Logger.Info("Winners saved successfully", "issue_id", issueID, "winner_count", len(winners))

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		utils.Logger.Error("Failed to commit transaction", "issue_id", issueID, "error", err)
		return utils.NewServiceError("failed to commit transaction", err)
	}

	return nil
}

// getWinnersFromChain identifies winners by comparing ticket numbers with results
func (s *LotteryService) getWinnersFromChain(issueID string, results []*big.Int) ([]models.Winner, error) {
	var winners []models.Winner
	offset := 0

	for {
		var tickets []models.LotteryTicket
		if err := s.db.Where("issue_id = ?", issueID).Limit(TicketBatchSize).Offset(offset).Find(&tickets).Error; err != nil {
			utils.Logger.Error("Failed to fetch tickets", "issue_id", issueID, "error", err)
			return nil, utils.NewServiceError("failed to fetch tickets", err)
		}
		if len(tickets) == 0 {
			break
		}

		for _, ticket := range tickets {
			ticketNumbers := parseBetContent(ticket.BetContent)
			if len(ticketNumbers) == ExpectedResultCount &&
				ticketNumbers[0].Cmp(results[0]) == 0 &&
				ticketNumbers[1].Cmp(results[1]) == 0 &&
				ticketNumbers[2].Cmp(results[2]) == 0 {
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
				winners = append(winners, winner)
			}
		}
		offset += TicketBatchSize
	}

	return winners, nil
}

// GetDrawnLotteryByIssueID retrieves drawn lottery information by issue ID
func GetDrawnLotteryByIssueID(issueID string) (*models.LotteryIssue, error) {
	utils.Logger.Info("Fetching drawn lottery", "issue_id", issueID)
	var issue models.LotteryIssue
	if err := db.DB.Where("issue_id = ?", issueID).Where("status = ?", models.IssueStatusDrawn).First(&issue).Error; err != nil {
		utils.Logger.Warn("Drawn lottery not found", "issue_id", issueID)
		return nil, utils.NewServiceError("drawn lottery not found", err)
	}
	utils.Logger.Info("Successfully fetched drawn lottery", "issue_id", issueID)
	return &issue, nil
}
