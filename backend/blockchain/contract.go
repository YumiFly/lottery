package blockchain

import (
	lotteryBlockchain "backend/blockchain/lottery"
	"backend/utils"

	"github.com/ethereum/go-ethereum/common"
)

// connectLotteryContract 连接到 LotteryManager 合约
func ConnectLotteryContract(contractAddress string) (*lotteryBlockchain.LotteryManager, error) {
	contractAddr := common.HexToAddress(contractAddress)
	contract, err := lotteryBlockchain.NewLotteryManager(contractAddr, Client)
	if err != nil {
		utils.Logger.Error("Failed to connect to lottery contract", "address", contractAddress, "error", err)
		return nil, utils.NewServiceError("failed to connect to lottery contract", err)
	}
	return contract, nil
}

// connectTokenContract 连接到 LOTToken 合约
func ConnectTokenContract(contractAddress string) (*lotteryBlockchain.LOTToken, error) {
	contractAddr := common.HexToAddress(contractAddress)
	contract, err := lotteryBlockchain.NewLOTToken(contractAddr, Client)
	if err != nil {
		utils.Logger.Error("Failed to connect to LOTToken contract", "address", contractAddress, "error", err)
		return nil, utils.NewServiceError("failed to connect to LOTToken contract", err)
	}
	return contract, nil
}
