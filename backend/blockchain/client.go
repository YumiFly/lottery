// blockchain/client.go
package blockchain

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"

	"backend/config"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client 全局区块链客户端
var Client *ethclient.Client

// Auth 全局交易授权
var Auth *bind.TransactOpts

// InitClient 初始化区块链客户端
func InitClient() {
	// 连接以太坊节点（例如通过 Infura）
	client, err := ethclient.Dial(config.AppConfig.EthereumNodeURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum node: %v", err)
	}

	// 加载管理员私钥（用于调用合约方法）
	privateKey, err := crypto.HexToECDSA(config.AppConfig.AdminPrivateKey)
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
		log.Fatalf("Failed to get nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("Failed to create transactor: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // 不发送以太币
	auth.GasLimit = uint64(300000) // 设置 Gas 限制
	auth.GasPrice = gasPrice

	// 设置全局变量
	Client = client
	Auth = auth
}
