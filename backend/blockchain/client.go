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

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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
	currentNonce     uint64        // 当前本地 Nonce
	lastNonceSync    time.Time     // 上次 Nonce 同步时间
	currentGas       *big.Int      // 当前 Gas 价格
	lastGasSync      time.Time     // 上次 Gas 价格同步时间
	currentGasLimit  uint64        // 当前 Gas 限制
	lastGasLimitSync time.Time     // 上次 Gas 限制同步时间
	syncInterval     time.Duration // 同步间隔，从配置读取
}

// BlockchainMgr 全局区块链管理器
var BlockchainMgr = &blockchainManager{}

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
	auth.GasLimit = uint64(50000000)

	Client = client
	Auth = auth
	BlockchainMgr.currentNonce = nonce
	BlockchainMgr.currentGas = gasPrice
	BlockchainMgr.currentGasLimit = auth.GasLimit
	BlockchainMgr.lastNonceSync = time.Now()
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

// syncNonceWithBlockchain 从区块链同步 Nonce
func (bm *blockchainManager) syncNonceWithBlockchain(ctx context.Context) error {
	nonce, err := Client.PendingNonceAt(ctx, Auth.From)
	if err != nil {
		utils.Logger.Error("Failed to sync nonce from blockchain", "error", err)
		return fmt.Errorf("failed to sync nonce from blockchain: %v", err)
	}
	if nonce > bm.currentNonce {
		bm.currentNonce = nonce
	}
	bm.lastNonceSync = time.Now()
	utils.Logger.Debug("Synced nonce from blockchain", "nonce", nonce)
	return nil
}

// syncGasPriceWithBlockchain 从区块链同步 Gas 价格
func (bm *blockchainManager) syncGasPriceWithBlockchain(ctx context.Context) error {
	gasPrice, err := Client.SuggestGasPrice(ctx)
	if err != nil {
		utils.Logger.Error("Failed to sync gas price from blockchain", "error", err)
		return fmt.Errorf("failed to sync gas price from blockchain: %v", err)
	}
	bm.currentGas = gasPrice
	bm.lastGasSync = time.Now()
	utils.Logger.Debug("Synced gas price from blockchain", "gas_price", gasPrice.String())
	return nil
}

// GetNextNonce 获取下一个可用 Nonce，每次都与区块链同步
func (bm *blockchainManager) GetNextNonce(ctx context.Context) (uint64, error) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	nonce, err := Client.PendingNonceAt(ctx, Auth.From)
	if err != nil {
		utils.Logger.Error("Failed to sync nonce from blockchain", "error", err)
		return 0, fmt.Errorf("failed to sync nonce from blockchain: %v", err)
	}
	Auth.Nonce = big.NewInt(int64(nonce))
	utils.Logger.Debug("Assigned next nonce", "nonce", nonce, "from_address", Auth.From.Hex())
	return nonce, nil
}

// GetCurrentGasPrice 获取当前 Gas 价格，并在需要时同步
func (bm *blockchainManager) GetCurrentGasPrice(ctx context.Context) (*big.Int, error) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if time.Since(bm.lastGasSync) >= bm.syncInterval {
		if err := bm.syncGasPriceWithBlockchain(ctx); err != nil {
			return nil, err
		}
	}

	Auth.GasPrice = bm.currentGas
	utils.Logger.Debug("Assigned current gas price", "gas_price", bm.currentGas.String())
	return bm.currentGas, nil
}

// GetCurrentGasLimit 获取当前 Gas 限制，动态调整需要结合具体交易
func (bm *blockchainManager) GetCurrentGasLimit(ctx context.Context, call func(*bind.TransactOpts) error) (uint64, error) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	// 如果提供了 call 函数，仍然支持现有逻辑
	if call != nil {
		opts := *Auth
		opts.GasLimit = 0
		if err := call(&opts); err != nil {
			utils.Logger.Warn("Failed to estimate gas limit with call, using default", "error", err)
			bm.currentGasLimit = 5000000 // 默认值 500 万
		} else if opts.GasLimit > 0 {
			estimatedGas := opts.GasLimit * 15 / 10
			if estimatedGas < 1000000 {
				estimatedGas = 1000000
			}
			bm.currentGasLimit = estimatedGas
			utils.Logger.Info("Estimated gas limit from call", "gas_limit", bm.currentGasLimit)
		} else {
			bm.currentGasLimit = 5000000 // 如果 call 未更新 GasLimit，使用默认值
			utils.Logger.Warn("Gas estimation returned 0, using default", "gas_limit", bm.currentGasLimit)
		}
	} else {
		// 如果没有提供 call，使用默认值
		bm.currentGasLimit = 5000000
		utils.Logger.Info("No gas estimation provided, using default", "gas_limit", bm.currentGasLimit)
	}

	Auth.GasLimit = bm.currentGasLimit
	bm.lastGasLimitSync = time.Now()
	utils.Logger.Debug("Assigned current gas limit", "gas_limit", bm.currentGasLimit)
	return bm.currentGasLimit, nil
}

// WithBlockchain 封装区块链操作
func WithBlockchain(ctx context.Context, estimateGas func(*bind.TransactOpts) error, fn func() (common.Hash, error)) error {
	if err := EnsureInitialized(); err != nil {
		return utils.NewServiceError("initialization check failed", err)
	}

	if _, err := BlockchainMgr.GetNextNonce(ctx); err != nil {
		return utils.NewServiceError("failed to get next nonce", err)
	}
	if _, err := BlockchainMgr.GetCurrentGasPrice(ctx); err != nil {
		return utils.NewServiceError("failed to get current gas price", err)
	}
	if _, err := BlockchainMgr.GetCurrentGasLimit(ctx, estimateGas); err != nil {
		return utils.NewServiceError("failed to get current gas limit", err)
	}

	txHash, err := fn()
	if err != nil {
		utils.Logger.Warn("Transaction failed, not rolling back nonce due to automining", "tx_hash", txHash.Hex())
		return err
	}
	return nil
}
