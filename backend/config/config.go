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

	BlockchainSyncInterval int
	MaxBlockchainRetries   int
	GasLimitIncreaseFactor float64

	RolloutContractAddress  string // Rollout 合约地址
	CallbackContractAddress string // Callback 合约地址

	// S3 配置
	Endpoint   string // S3 端点
	BucketName string // S3 存储桶名称
	Region     string // S3 区域
	AccessKey  string // S3 访问密钥
	SecretKey  string // S3 秘密密钥
}

var AppConfig AppConfigStruct

func getEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvFloat(key string, defaultValue float64) float64 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

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

<<<<<<< HEAD
		RolloutContractAddress: os.Getenv("ROLLOUT_CONTRACT_ADDRESS"),
		TokenContractAddress:   os.Getenv("TOKEN_CONTRACT_ADDRESS"),
		BlockchainSyncInterval: getEnvInt("BLOCKCHAIN_SYNC_INTERVAL", 60),
		MaxBlockchainRetries:   getEnvInt("MAX_BLOCKCHAIN_RETRIES", 3),
		GasLimitIncreaseFactor: getEnvFloat("GAS_LIMIT_INCREASE_FACTOR", 1.5),
=======
		RolloutContractAddress:  os.Getenv("ROLLOUT_CONTRACT_ADDRESS"),
		CallbackContractAddress: os.Getenv("CALLBACK_CONTRACT_ADDRESS"),

		// S3 配置
		Endpoint:   os.Getenv("S3_ENDPOINT"),
		BucketName: os.Getenv("S3_BUCKET_NAME"),
		Region:     os.Getenv("S3_REGION"),
		AccessKey:  os.Getenv("S3_ACCESS_KEY"),
		SecretKey:  os.Getenv("S3_SECRET_KEY"),
>>>>>>> 2add266 (上传文件增加支持S3存储)
	}
}
