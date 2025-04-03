package services

import (
	"backend/blockchain"
	lotteryBlockchain "backend/blockchain/lottery"
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

// CreateLottery 创建彩票
func CreateLottery(lottery *models.Lottery) error {
	utils.Logger.Info("Creating lottery", "lottery_id", lottery.LotteryID, "type_id", lottery.TypeID)

	estimateGas := func(opts *bind.TransactOpts) error {
		adminAddr := blockchain.Auth.From
		ownerAddr := common.HexToAddress(lottery.RegisteredAddr)
		rolloutContractAddr := common.HexToAddress(lottery.RolloutContractAddress)
		tokenContractAddr := common.HexToAddress(config.AppConfig.TokenContractAddress)
		supply, _ := new(big.Int).SetString(big.NewInt(lottery.TicketSupply).String(), 10)
		price, _ := new(big.Int).SetString(big.NewFloat(lottery.TicketPrice).Text('f', 0), 10)
		_, _, _, err := lotteryBlockchain.DeployLotteryManager(
			opts, blockchain.Client, adminAddr, ownerAddr, rolloutContractAddr,
			lottery.TicketName, supply, price, tokenContractAddr,
		)
		if err != nil {
			utils.Logger.Warn("Gas estimation failed", "error", err)
		} else {
			utils.Logger.Info("Gas estimation succeeded", "estimated_gas", opts.GasLimit)
		}
		return err
	}

	executeTx := func() (common.Hash, error) {
		var lotteryType models.LotteryType
		if err := db.DB.Where("type_id = ?", lottery.TypeID).First(&lotteryType).Error; err != nil {
			utils.Logger.Warn("Lottery type not found", "type_id", lottery.TypeID)
			return common.Hash{}, utils.NewServiceError("lottery type not found", err)
		}

		lottery.CreatedAt = time.Now()
		lottery.UpdatedAt = time.Now()

		supply, ok := new(big.Int).SetString(big.NewInt(lottery.TicketSupply).String(), 10)
		if !ok {
			utils.Logger.Error("Invalid ticket supply format", "supply", lottery.TicketSupply)
			return common.Hash{}, utils.NewServiceError("invalid ticket supply format", nil)
		}
		price, ok := new(big.Int).SetString(big.NewFloat(lottery.TicketPrice).Text('f', 0), 10)
		if !ok {
			utils.Logger.Error("Invalid ticket price format", "price", lottery.TicketPrice)
			return common.Hash{}, utils.NewServiceError("invalid ticket price format", nil)
		}

		adminAddr := blockchain.Auth.From
		ownerAddr := common.HexToAddress(lottery.RegisteredAddr)
		rolloutContractAddr := common.HexToAddress(lottery.RolloutContractAddress)
		tokenContractAddr := common.HexToAddress(config.AppConfig.TokenContractAddress)

		nonce, err := blockchain.BlockchainMgr.GetNextNonce(context.Background())
		if err != nil {
			return common.Hash{}, utils.NewServiceError("failed to get next nonce before deployment", err)
		}
		blockchain.Auth.Nonce = big.NewInt(int64(nonce))
		utils.Logger.Info("Deploying LotteryManager contract", "admin", adminAddr.Hex(), "owner", ownerAddr.Hex(), "nonce", nonce, "gas_limit", blockchain.Auth.GasLimit)

		contractAddr, tx, _, err := lotteryBlockchain.DeployLotteryManager(
			blockchain.Auth, blockchain.Client, adminAddr, ownerAddr, rolloutContractAddr,
			lottery.TicketName, supply, price, tokenContractAddr,
		)
		if err != nil {
			utils.Logger.Error("Failed to deploy LotteryManager contract", "error", err)
			if tx == nil {
				return common.Hash{}, utils.NewServiceError("failed to deploy LotteryManager contract", err)
			}
			return tx.Hash(), utils.NewServiceError("failed to deploy LotteryManager contract", err)
		}

		// 等待交易被挖出并获取收据
		receipt, err := bind.WaitMined(context.Background(), blockchain.Client, tx)
		if err != nil {
			utils.Logger.Error("Failed to confirm contract deployment", "tx_hash", tx.Hash().Hex(), "error", err)
			return tx.Hash(), utils.NewServiceError("failed to confirm contract deployment", err)
		}
		if receipt.Status != 1 {
			utils.Logger.Error("Contract deployment transaction failed", "tx_hash", tx.Hash().Hex(), "status", receipt.Status)
			return tx.Hash(), utils.NewServiceError("contract deployment transaction failed", nil)
		}

		// 记录成功的交易信息，包括 Gas 使用量
		utils.Logger.Info("Transaction submitted successfully", "tx_hash", tx.Hash().Hex(), "gas_used", receipt.GasUsed)

		lottery.ContractAddress = contractAddr.Hex()
		if err := db.DB.Create(lottery).Error; err != nil {
			utils.Logger.Error("Failed to save lottery to database", "error", err)
			return common.Hash{}, utils.NewServiceError("failed to save lottery to database", err)
		}

		utils.Logger.Info("Lottery created successfully", "lottery_id", lottery.LotteryID, "contract_address", lottery.ContractAddress)
		return tx.Hash(), nil
	}

	return blockchain.WithBlockchain(context.Background(), estimateGas, executeTx)
}

// GetAllLotteries 获取所有彩票
func GetAllLotteries() ([]models.Lottery, error) {
	utils.Logger.Info("Fetching all lotteries")
	var lotteries []models.Lottery
	if err := db.DB.Find(&lotteries).Error; err != nil {
		utils.Logger.Error("Failed to fetch lotteries", "error", err)
		return nil, utils.NewServiceError("failed to fetch lotteries", err)
	}
	utils.Logger.Info("Fetched lotteries", "count", len(lotteries))
	return lotteries, nil
}

// GetLotteryByTypeID 根据类型ID获取彩票
func GetLotteryByTypeID(typeID string) (*models.Lottery, error) {
	utils.Logger.Info("Fetching lottery by type ID", "type_id", typeID)
	var lottery models.Lottery
	if err := db.DB.Where("type_id = ?", typeID).First(&lottery).Error; err != nil {
		utils.Logger.Warn("Lottery not found for type ID", "type_id", typeID)
		return nil, utils.NewServiceError("lottery not found", err)
	}
	utils.Logger.Info("Fetched lottery", "lottery_id", lottery.LotteryID)
	return &lottery, nil
}
