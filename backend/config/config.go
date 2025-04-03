// config/config.go
package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// AppConfig 应用程序配置结构体
type AppConfigStruct struct {

	// 数据库配置
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DB_SSLMODE  string
	DB_TIMEZONE string
	JWTSecret   string

	//blockchain配置
	EthereumNodeURL string // 以太坊节点 URL（例如 Infura）
	AdminPrivateKey string // 管理员私钥（用于调用 verifyKYC）

	RolloutContractAddress string // Rollout 合约地址
	TokenContractAddress   string // Token 合约地址

	BlockchainSyncInterval int // 区块链同步间隔（以秒为单位）
}

var AppConfig AppConfigStruct

// LoadConfig 加载配置
func LoadConfig() {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found, using default environment variables: %v", err)
	}

	// 从环境变量加载配置
	dbPort := os.Getenv("DB_PORT")
	//dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("Invalid DB_PORT value: %v", err)
	}

	// 验证必要的环境变量是否存在
	requiredEnvVars := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "JWT_SECRET"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("Missing required environment variable: %s", envVar)
		}
	}
	syncInterval, err := strconv.Atoi(os.Getenv("BLOCKCHAIN_SYNC_INTERVAL"))
	if err != nil {
		log.Fatalf("Invalid BLOCKCHAIN_SYNC_INTERVAL value: %v", err)
	}

	AppConfig = AppConfigStruct{
		DBHost:      os.Getenv("DB_HOST"),
		DBPort:      dbPort,
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		DB_SSLMODE:  os.Getenv("DB_SSLMODE"),
		DB_TIMEZONE: os.Getenv("DB_TIMEZONE"),
		JWTSecret:   os.Getenv("JWT_SECRET"),

		EthereumNodeURL: os.Getenv("ETHEREUM_NODE_URL"),
		AdminPrivateKey: os.Getenv("ADMIN_PRIVATE_KEY"),

		RolloutContractAddress: os.Getenv("ROLLOUT_CONTRACT_ADDRESS"),
		TokenContractAddress:   os.Getenv("TOKEN_CONTRACT_ADDRESS"),
		BlockchainSyncInterval: syncInterval, // 默认为10秒
	}
}
