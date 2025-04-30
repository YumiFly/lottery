package lottery

import (
	"context"
	"math/big"
	"time"

	"backend/blockchain"
	lotteryBlockchain "backend/blockchain/lottery"
	"backend/config"
	"backend/models"
	"backend/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// CreateLotteryParams defines the parameters for creating a lottery
type CreateLotteryParams struct {
	TypeID                 string
	TicketName             string
	TicketSupply           int64
	TicketPrice            float64
	BettingRules           string
	PrizeStructure         string
	RegisteredAddr         string
	RolloutContractAddress string
}

// LotteryService encapsulates lottery creation business logic
type LotteryCreateService struct {
	db *gorm.DB
}

// NewLotteryService creates a new LotteryService instance
func NewLotteryCreateService(db *gorm.DB) *LotteryCreateService {
	return &LotteryCreateService{db: db}
}

// validateCreateLotteryParams validates the parameters for creating a lottery
func (s *LotteryCreateService) validateCreateLotteryParams(params CreateLotteryParams) error {
	// Validate type_id exists
	var lotteryType models.LotteryType
	if err := s.db.WithContext(context.Background()).
		Where("type_id = ?", params.TypeID).
		First(&lotteryType).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Logger.Warn("Lottery type not found", "type_id", params.TypeID)
			return utils.NewBadRequestError("Lottery type not found", nil)
		}
		return utils.NewInternalError("Failed to check lottery type ID", errors.Wrap(err, "database error"))
	}

	// Validate ticket supply and price
	if params.TicketSupply <= 0 {
		return utils.NewBadRequestError("Ticket supply must be positive", nil)
	}
	if params.TicketPrice <= 0 {
		return utils.NewBadRequestError("Ticket price must be positive", nil)
	}

	// Validate Ethereum addresses
	if !common.IsHexAddress(params.RegisteredAddr) {
		return utils.NewBadRequestError("Invalid registered address", nil)
	}
	if !common.IsHexAddress(params.RolloutContractAddress) {
		return utils.NewBadRequestError("Invalid rollout contract address", nil)
	}

	// Validate field lengths
	if len(params.TicketName) == 0 || len(params.TicketName) > 100 {
		return utils.NewBadRequestError("Ticket name must be between 1 and 100 characters", nil)
	}
	if len(params.TypeID) == 0 || len(params.TypeID) > 36 {
		return utils.NewBadRequestError("Type ID must be between 1 and 36 characters", nil)
	}

	return nil
}

// CreateLottery creates a new lottery
//
// Parameters:
//   - ctx: Request context
//   - params: Creation parameters, including type_id, ticket_name, etc.
//
// Returns:
//   - *Lottery: The created lottery record
//   - common.Hash: Blockchain transaction hash
//   - error: Creation error or invalid parameters
func (s *LotteryCreateService) CreateLottery(ctx context.Context, params CreateLotteryParams) (*models.Lottery, common.Hash, error) {
	// Validate parameters
	if err := s.validateCreateLotteryParams(params); err != nil {
		return nil, common.Hash{}, err
	}

	// Convert supply and price to big.Int
	supply, ok := new(big.Int).SetString(big.NewInt(params.TicketSupply).String(), 10)
	if !ok {
		utils.Logger.Error("Invalid ticket supply format", "supply", params.TicketSupply)
		return nil, common.Hash{}, utils.NewBadRequestError("Invalid ticket supply format", nil)
	}
	price, ok := new(big.Int).SetString(big.NewFloat(params.TicketPrice).Text('f', 0), 10)
	if !ok {
		utils.Logger.Error("Invalid ticket price format", "price", params.TicketPrice)
		return nil, common.Hash{}, utils.NewBadRequestError("Invalid ticket price format", nil)
	}

	// Construct lottery record
	lottery := models.Lottery{
		LotteryID:              uuid.NewString(),
		TypeID:                 params.TypeID,
		TicketName:             params.TicketName,
		TicketSupply:           params.TicketSupply,
		BettingRules:           params.BettingRules,
		PrizeStructure:         params.PrizeStructure,
		TicketPrice:            params.TicketPrice,
		RegisteredAddr:         params.RegisteredAddr,
		RolloutContractAddress: params.RolloutContractAddress,
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}
	lotteryType := models.LotteryType{}
	// Log creation attempt
	utils.Logger.Info("Creating lottery", "lottery_id", lottery.LotteryID, "type_id", lottery.TypeID)

	// Execute blockchain transaction
	data := []byte{} // Empty data, gas estimation handled by blockchain package
	executeTx := func() (common.Hash, error) {
		// Re-validate lottery type (for transaction consistency)
		if err := s.db.WithContext(ctx).
			Where("type_id = ?", lottery.TypeID).
			First(&lotteryType).Error; err != nil {
			utils.Logger.Warn("Lottery type not found", "type_id", lottery.TypeID)
			return common.Hash{}, utils.NewBadRequestError("Lottery type not found", err)
		}

		// Prepare blockchain parameters
		adminAddr := blockchain.Auth.From
		ownerAddr := common.HexToAddress(lottery.RegisteredAddr)
		rolloutContractAddr := common.HexToAddress(lottery.RolloutContractAddress)
		tokenContractAddr := common.HexToAddress(config.AppConfig.TokenContractAddress)

		// Deploy contract
		utils.Logger.Info("Deploying LotteryManager contract",
			"admin", adminAddr.Hex(),
			"owner", ownerAddr.Hex(),
			"nonce", blockchain.Auth.Nonce,
			"gas_limit", blockchain.Auth.GasLimit)
		contractAddr, tx, _, err := lotteryBlockchain.DeployLotteryManager(
			blockchain.Auth,
			blockchain.Client,
			adminAddr,
			ownerAddr,
			rolloutContractAddr,
			lottery.TicketName,
			supply,
			price,
			tokenContractAddr,
		)
		if err != nil {
			utils.Logger.Error("Failed to deploy LotteryManager contract", "error", err)
			if tx == nil {
				return common.Hash{}, utils.NewInternalError("Failed to deploy LotteryManager contract", err)
			}
			return tx.Hash(), utils.NewInternalError("Failed to deploy LotteryManager contract, transaction failed", err)
		}

		// Wait for transaction confirmation
		receipt, err := bind.WaitMined(ctx, blockchain.Client, tx)
		if err != nil {
			utils.Logger.Error("Failed to confirm contract deployment", "tx_hash", tx.Hash().Hex(), "error", err)
			return tx.Hash(), utils.NewInternalError("Failed to confirm contract deployment", err)
		}
		if receipt.Status != 1 {
			utils.Logger.Error("Contract deployment transaction failed", "tx_hash", tx.Hash().Hex(), "status", receipt.Status)
			return tx.Hash(), utils.NewInternalError("Contract deployment transaction failed", nil)
		}

		// Update gas history
		blockchain.BlockchainMgr.UpdateGasHistory(receipt)
		utils.Logger.Info("Transaction submitted successfully", "tx_hash", tx.Hash().Hex(), "gas_used", receipt.GasUsed)

		// Update contract address and save to database
		lottery.ContractAddress = contractAddr.Hex()
		if err := s.db.WithContext(ctx).Create(&lottery).Error; err != nil {
			utils.Logger.Error("Failed to save lottery to database", "error", err)
			return common.Hash{}, utils.NewInternalError("Failed to save lottery to database", err)
		}

		utils.Logger.Info("Lottery created successfully",
			"lottery_id", lottery.LotteryID,
			"contract_address", lottery.ContractAddress)
		return tx.Hash(), nil
	}
	// Execute blockchain transaction
	txhash, err := blockchain.WithBlockchain(ctx, data, executeTx)
	if err != nil {
		return nil, common.Hash{}, err
	}
	lottery.LotteryType = lotteryType

	return &lottery, txhash, nil
}
