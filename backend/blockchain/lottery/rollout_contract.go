package lottery

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"reflect"

	"backend/blockchain"
	"backend/config"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// RolloutContract 封装 SimpleRollout 合约实例
type RolloutContract struct {
	Instance *SimpleRollout // 由 abigen 生成的合约绑定
}

// 全局 SimpleRollout 合约实例
var RolloutInstance *RolloutContract

// InitRollout 加载 SimpleRollout 合约
func InitRollout() {
	if blockchain.Client == nil {
		log.Fatalf("Blockchain client not initialized. Call InitClient first.")
	}

	// 加载 SimpleRollout 合约
	contractAddress := common.HexToAddress(config.AppConfig.RolloutContractAddress)
	instance, err := NewSimpleRollout(contractAddress, blockchain.Client)
	if err != nil {
		log.Fatalf("Failed to load SimpleRollout contract: %v", err)
	}

	RolloutInstance = &RolloutContract{
		Instance: instance,
	}
}

// RolloutCall 调用链上 rolloutCall 方法
func RolloutCall(callbackAddress common.Address) (uint64, error) {
	if RolloutInstance == nil {
		log.Fatalf("SimpleRollout contract not initialized. Call InitRollout first.")
	}

	// 调用 rolloutCall 方法
	tx, err := RolloutInstance.Instance.RolloutCall(blockchain.Auth, callbackAddress)
	if err != nil {
		return 0, err
	}
	// 等待交易确认
	receipt, err := bind.WaitMined(context.Background(), blockchain.Client, tx)
	if err != nil {
		return 0, err
	}

	// 从交易日志中提取 requestID
	for _, log := range receipt.Logs {
		event, err := RolloutInstance.Instance.ParseDiceRolled(*log)
		if err == nil {
			return event.RequestId.Uint64(), nil
		}
	}

	return 0, nil
}

// ListenForDiceLanded 监听 DiceLanded 事件，获取随机数结果
// ListenForDiceLanded 监听 DiceLanded 事件，获取随机数结果
func ListenForDiceLanded(requestID *big.Int, onResult func(results []*big.Int)) error {
	if RolloutInstance == nil {
		log.Fatalf("SimpleRollout contract not initialized. Call InitRollout first.")
	}

	// 通过反射获取底层的 BoundContract
	boundContractField := reflect.ValueOf(RolloutInstance.Instance).Elem().FieldByName("SimpleRollout")
	if !boundContractField.IsValid() {
		return fmt.Errorf("failed to get SimpleRollout field")
	}

	boundContract, ok := boundContractField.Interface().(*bind.BoundContract)
	if !ok {
		return fmt.Errorf("failed to cast to BoundContract")
	}

	// 获取 ABI
	abiField := reflect.ValueOf(boundContract).Elem().FieldByName("abi")
	if !abiField.IsValid() {
		return fmt.Errorf("failed to get abi field")
	}

	contractAbi, ok := abiField.Interface().(abi.ABI)
	if !ok {
		return fmt.Errorf("failed to cast to abi.ABI")
	}

	// 获取 DiceLanded 事件的主题
	diceLandedEvent, exists := contractAbi.Events["DiceLanded"]
	if !exists {
		return fmt.Errorf("DiceLanded event not found in ABI")
	}

	// 创建事件监听器
	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(config.AppConfig.RolloutContractAddress)},
		Topics:    [][]common.Hash{{diceLandedEvent.ID}, {common.BigToHash(requestID)}},
	}

	logs := make(chan types.Log)
	sub, err := blockchain.Client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return err
	}

	// 监听事件
	go func() {
		for {
			select {
			case err := <-sub.Err():
				log.Printf("Error in event subscription: %v", err)
				return
			case vLog := <-logs:
				event, err := RolloutInstance.Instance.ParseDiceLanded(vLog)
				if err != nil {
					log.Printf("Failed to parse DiceLanded event: %v", err)
					continue
				}
				// 调用回调函数，传递随机数结果
				onResult(event.Results)
				return
			}
		}
	}()

	return nil
}
