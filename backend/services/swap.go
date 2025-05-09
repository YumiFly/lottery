package services

import (
	"backend/blockchain"
	"backend/config"
	"backend/models"
	"backend/utils"
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// 添加修改稳定币
func SetStableCoin(stbCoin *models.LotterySTBCoin) error {
	utils.Logger.Info("Set Stable Coin", "name:", stbCoin.STBCoinName, "STBCoinAddr:", stbCoin.STBCoinAddr, "STBCoinRate:", stbCoin.STB2LOTRate, "STBRecvAddr:", stbCoin.STBReceiverAddr)

	executeTx := func() (common.Hash, error) {
		// 获取 LOTToken 合约实例
		tokenContract, err := blockchain.ConnectTokenContract(config.AppConfig.TokenContractAddress)
		if err != nil {
			utils.Logger.Error("Failed to connect to LOTToken contract", "error", err)
			return common.Hash{}, utils.NewServiceError("failed to connect to LOTToken contract", err)
		}

		// 调用 LOTToken 合约的 setStableCoin 函数
		tx, err := tokenContract.SetStablecoin(blockchain.Auth,
			common.HexToAddress(stbCoin.STBCoinAddr),
			stbCoin.STBCoinName,
			big.NewInt(stbCoin.STB2LOTRate),
			common.HexToAddress(stbCoin.STBReceiverAddr))
		if err != nil {
			utils.Logger.Error("Failed to set stable coin", "error", err)
			if tx != nil {
				return tx.Hash(), utils.NewServiceError("failed to set stable coin", err)
			}
			return common.Hash{}, utils.NewServiceError("failed to set stable coin", err)
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

		utils.Logger.Info("Set Stable Coin successfully", "tx_hash", tx.Hash().Hex())
		return tx.Hash(), nil
	}

	// 使用空的 data 调用 WithBlockchain，Gas 估算依赖内部逻辑
	data := []byte{}
	_, err := blockchain.WithBlockchain(context.Background(), data, executeTx)
	return err
}

// 移除稳定币
func RemoveStableCoin(stbCoin *models.LotterySTBCoin) error {
	utils.Logger.Info("Remove Stable Coin", "name:", stbCoin.STBCoinName, "STBCoinAddr:", stbCoin.STBCoinAddr, "STBCoinRate:", stbCoin.STB2LOTRate, "STBRecvAddr:", stbCoin.STBReceiverAddr)

	executeTx := func() (common.Hash, error) {
		// 获取 LOTToken 合约实例
		tokenContract, err := blockchain.ConnectTokenContract(config.AppConfig.TokenContractAddress)
		if err != nil {
			utils.Logger.Error("Failed to connect to LOTToken contract", "error", err)
			return common.Hash{}, utils.NewServiceError("failed to connect to LOTToken contract", err)
		}

		// 调用 LOTToken 合约的 setStableCoin 函数
		tx, err := tokenContract.RemoveStablecoin(blockchain.Auth, common.HexToAddress(stbCoin.STBCoinAddr))
		if err != nil {
			utils.Logger.Error("Failed to remove stable coin", "error", err)
			if tx != nil {
				return tx.Hash(), utils.NewServiceError("failed to remove stable coin", err)
			}
			return common.Hash{}, utils.NewServiceError("failed to remove stable coin", err)
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

		utils.Logger.Info("Remove Stable Coin successfully", "tx_hash", tx.Hash().Hex())
		return tx.Hash(), nil
	}

	// 使用空的 data 调用 WithBlockchain，Gas 估算依赖内部逻辑
	data := []byte{}
	_, err := blockchain.WithBlockchain(context.Background(), data, executeTx)
	return err
}
