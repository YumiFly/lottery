package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"
	"time"

	"backend/config"
	"backend/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client 全局区块链客户端
var Client *ethclient.Client

// Auth 全局交易授权
var Auth *bind.TransactOpts

// blockchainManager 管理 Nonce、Gas 价格和 Gas 限制的结构体
type blockchainManager struct {
	mu               sync.Mutex    // 保护并发访问
	currentGas       *big.Int      // 当前 Gas 价格
	lastGasSync      time.Time     // 上次 Gas 价格同步时间
	gasHistory       []uint64      // Gas 使用历史，用于自适应管理
	currentGasLimit  uint64        // 当前 Gas 限制
	lastGasLimitSync time.Time     // 上次 Gas 限制同步时间
	syncInterval     time.Duration // 同步间隔，从配置读取
	nonceMutex       sync.Mutex    // 保护 Nonce 获取的互斥锁
}

// BlockchainMgr 全局区块链管理器
var BlockchainMgr = &blockchainManager{
	gasHistory: make([]uint64, 0, 10), // 初始化 Gas 历史，最大保留 10 次记录
}

// InitClient 初始化区块链客户端
func InitClient() {
	client, err := ethclient.Dial(config.AppConfig.EthereumNodeURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum node: %v", err)
	}

	privateKeyHex := config.AppConfig.AdminPrivateKey
	if strings.HasPrefix(privateKeyHex, "0x") {
		privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")
		utils.Logger.Info("Removed '0x' prefix from private key")
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("Failed to cast public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Failed to get initial nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest initial gas price: %v", err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("Failed to create transactor: %v", err)
	}
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(5000000) // 初始 GasLimit 设置为 500 万

	Client = client
	Auth = auth
	BlockchainMgr.currentGas = gasPrice
	BlockchainMgr.currentGasLimit = auth.GasLimit
	BlockchainMgr.lastGasSync = time.Now()
	BlockchainMgr.lastGasLimitSync = time.Now()
	BlockchainMgr.syncInterval = time.Duration(config.AppConfig.BlockchainSyncInterval) * time.Second

	utils.Logger.Info("Blockchain client initialized", "nonce", nonce, "gas_price", gasPrice.String(), "gas_limit", auth.GasLimit)
}

// EnsureInitialized 检查区块链客户端和授权是否初始化
func EnsureInitialized() error {
	if Client == nil || Auth == nil {
		utils.Logger.Error("Blockchain client or auth not initialized")
		return fmt.Errorf("blockchain client or auth not initialized")
	}
	return nil
}

// GetNextNonce 获取下一个可用 Nonce，每次实时获取
func (bm *blockchainManager) GetNextNonce(ctx context.Context) (uint64, error) {
	bm.nonceMutex.Lock()
	defer bm.nonceMutex.Unlock()

	currentNonce, err := Client.PendingNonceAt(ctx, Auth.From)
	if err != nil {
		return 0, fmt.Errorf("failed to get current nonce: %v", err)
	}

	return currentNonce, nil
}

// GetNextNonceForFunc 获取函数内部的下一个可用 Nonce
func (bm *blockchainManager) GetNextNonceForFunc(funcNonce uint64) uint64 {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	funcNonce++
	return funcNonce
}

// GetCurrentGasPrice 获取当前 Gas 价格，并在需要时同步
func (bm *blockchainManager) GetCurrentGasPrice(ctx context.Context) (*big.Int, error) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if time.Since(bm.lastGasSync) >= bm.syncInterval {
		gasPrice, err := Client.SuggestGasPrice(ctx)
		if err != nil {
			utils.Logger.Error("Failed to sync gas price", "error", err)
			return nil, fmt.Errorf("failed to sync gas price: %v", err)
		}
		bm.currentGas = gasPrice
		bm.lastGasSync = time.Now()
		utils.Logger.Debug("Synced gas price from blockchain", "gas_price", gasPrice.String())
	}

	Auth.GasPrice = bm.currentGas
	utils.Logger.Debug("Assigned current gas price", "gas_price", bm.currentGas.String())
	return bm.currentGas, nil
}

// UpdateGasHistory 更新 Gas 使用历史
func (bm *blockchainManager) UpdateGasHistory(receipt *types.Receipt) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	bm.gasHistory = append(bm.gasHistory, receipt.GasUsed)
	if len(bm.gasHistory) > 10 {
		bm.gasHistory = bm.gasHistory[1:]
	}
	utils.Logger.Debug("Updated gas history", "gas_used", receipt.GasUsed, "history_size", len(bm.gasHistory))
}

// GetCurrentGasLimit 自适应 Gas 管理
func (bm *blockchainManager) GetCurrentGasLimit(ctx context.Context, data []byte) (uint64, error) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if len(data) > 0 {
		msg := ethereum.CallMsg{
			From:     Auth.From,
			GasPrice: Auth.GasPrice,
			Data:     data,
		}
		gasLimit, err := Client.EstimateGas(ctx, msg)
		if err != nil {
			utils.Logger.Warn("Failed to estimate gas", "error", err)
			return 0, fmt.Errorf("failed to estimate gas: %v", err)
		}
		bm.currentGasLimit = uint64(float64(gasLimit) * config.AppConfig.GasLimitIncreaseFactor)
		if bm.currentGasLimit < 1000000 {
			bm.currentGasLimit = 1000000
		}
		bm.lastGasLimitSync = time.Now()
		utils.Logger.Info("Estimated gas limit from blockchain", "gas_limit", bm.currentGasLimit)
		Auth.GasLimit = bm.currentGasLimit
		return bm.currentGasLimit, nil
	}

	if len(bm.gasHistory) == 0 {
		bm.currentGasLimit = 5000000
	} else {
		var avgGas uint64
		for _, gas := range bm.gasHistory {
			avgGas += gas
		}
		avgGas /= uint64(len(bm.gasHistory))
		bm.currentGasLimit = uint64(float64(avgGas) * config.AppConfig.GasLimitIncreaseFactor)
		if bm.currentGasLimit < 1000000 {
			bm.currentGasLimit = 1000000
		}
	}

	bm.lastGasLimitSync = time.Now()
	Auth.GasLimit = bm.currentGasLimit
	utils.Logger.Info("Assigned gas limit from history or default", "gas_limit", bm.currentGasLimit)
	return bm.currentGasLimit, nil
}

// WithBlockchain 封装区块链操作，包含错误重试机制
func WithBlockchain(ctx context.Context, data []byte, fn func() (common.Hash, error)) (common.Hash, error) {
	if err := EnsureInitialized(); err != nil {
		return common.Hash{}, utils.NewServiceError("initialization check failed", err)
	}

	var lastErr error
	var txHash common.Hash
	for attempt := 0; attempt < config.AppConfig.MaxBlockchainRetries; attempt++ {
		// 实时获取 Nonce 和 Gas 参数
		nonce, err := BlockchainMgr.GetNextNonce(ctx)
		if err != nil {
			return common.Hash{}, utils.NewServiceError("failed to get next nonce", err)
		}
		Auth.Nonce = big.NewInt(int64(nonce))
		utils.Logger.Debug("Transaction attempt", "attempt", attempt+1, "nonce", nonce)

		if _, err := BlockchainMgr.GetCurrentGasPrice(ctx); err != nil {
			return common.Hash{}, utils.NewServiceError("failed to get current gas price", err)
		}
		if gasLimit, err := BlockchainMgr.GetCurrentGasLimit(ctx, data); err != nil {
			return common.Hash{}, utils.NewServiceError("failed to get current gas limit", err)
		} else {
			Auth.GasLimit = gasLimit
		}

		// 执行交易
		txHash, err = fn()
		if err != nil {
			if strings.Contains(err.Error(), "nonce too low") || strings.Contains(err.Error(), "nonce too high") {
				utils.Logger.Warn("Nonce issue, retrying immediately", "attempt", attempt+1, "nonce", nonce)
				continue
			}
			if strings.Contains(err.Error(), "ran out of gas") {
				utils.Logger.Warn("Transaction ran out of gas, retrying with increased limit", "attempt", attempt+1, "gas_limit", Auth.GasLimit)
				Auth.GasLimit = uint64(float64(Auth.GasLimit) * config.AppConfig.GasLimitIncreaseFactor)
				continue
			}
			utils.Logger.Error("Transaction failed", "tx_hash", txHash.Hex(), "error", err)
			lastErr = err
			continue
		}
		return txHash, nil
	}
	return common.Hash{}, utils.NewServiceError(fmt.Sprintf("max retries exceeded, last error: %v", lastErr), lastErr)
}
