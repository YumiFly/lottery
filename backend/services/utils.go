package services

import (
	"backend/blockchain"
	lotteryBlockchain "backend/blockchain/lottery"
	"backend/db"
	"backend/models"
	"backend/utils"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// connectLotteryContract 连接到 LotteryManager 合约
func connectLotteryContract(contractAddress string) (*lotteryBlockchain.LotteryManager, error) {
	contractAddr := common.HexToAddress(contractAddress)
	contract, err := lotteryBlockchain.NewLotteryManager(contractAddr, blockchain.Client)
	if err != nil {
		utils.Logger.Error("Failed to connect to lottery contract", "address", contractAddress, "error", err)
		return nil, utils.NewServiceError("failed to connect to lottery contract", err)
	}
	return contract, nil
}

// parseBetContent 解析投注内容
func parseBetContent(content string) []*big.Int {
	parts := strings.Split(content, ",")
	result := make([]*big.Int, 0, len(parts))
	for _, part := range parts {
		num, ok := new(big.Int).SetString(strings.TrimSpace(part), 10)
		if ok {
			result = append(result, num)
		}
	}
	return result
}

// parseBetContentV2bigIntSlice 解析投注内容（重命名版本）
func parseBetContentV2bigIntSlice(content string) []*big.Int {
	return parseBetContent(content)
}

// GetAllPools 获取所有奖池总额
func GetAllPools() (int64, error) {
	utils.Logger.Info("Fetching all pools")
	var issues []models.LotteryIssue
	if err := db.DB.Where("sale_end_time > ?", time.Now()).Find(&issues).Error; err != nil {
		utils.Logger.Error("Failed to fetch issues for pools", "error", err)
		return 0, utils.NewServiceError("failed to fetch issues for pools", err)
	}

	totalPool := int64(0)
	for _, issue := range issues {
		issuePool := issue.PrizePool
		totalPool += int64(issuePool)
	}
	utils.Logger.Info("Fetched total pools", "total", totalPool)
	return totalPool, nil
}
