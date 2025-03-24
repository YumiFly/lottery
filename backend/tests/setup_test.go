// tests/setup_test.go
package tests

import (
	"backend/config"
	"backend/db"
	"backend/models"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type TestSuite struct {
	DB *gorm.DB
}

func SetupTestDB() *TestSuite {
	// 加载 .env 文件
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// 验证必要的环境变量是否存在
	requiredEnvVars := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "JWT_SECRET"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("Missing required environment variable: %s", envVar)
		}
	}

	// 设置测试数据库名称
	testDBName := os.Getenv("DB_NAME") + "_test"
	os.Setenv("DB_NAME", testDBName)

	// 使用 postgres 数据库检查并创建测试数据库
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	// 连接到 postgres 数据库
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable", dbHost, dbPort, dbUser, dbPassword)
	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to postgres database: %v", err)
	}
	defer sqlDB.Close()

	// 测试连接
	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping postgres database: %v", err)
	}

	// 检查测试数据库是否存在
	var dbExists bool
	err = sqlDB.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", testDBName).Scan(&dbExists)
	if err != nil {
		log.Fatalf("Failed to check if database exists: %v", err)
	}

	// 如果数据库不存在，则创建
	if !dbExists {
		_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", testDBName))
		if err != nil {
			log.Fatalf("Failed to create test database %s: %v", testDBName, err)
		}
		log.Printf("Test database %s created successfully", testDBName)
	} else {
		log.Printf("Test database %s already exists", testDBName)
	}

	// 加载配置并初始化数据库
	config.LoadConfig()
	db.InitDB()

	// 清理数据库
	db.DB.Exec("DROP TABLE IF EXISTS kyc_verification_history CASCADE;")
	db.DB.Exec("DROP TABLE IF EXISTS kyc_data CASCADE;")
	db.DB.Exec("DROP TABLE IF EXISTS customers CASCADE;")
	db.DB.Exec("DROP TABLE IF EXISTS role_menus CASCADE;")
	db.DB.Exec("DROP TABLE IF EXISTS roles CASCADE;")

	// 自动迁移
	db.DB.AutoMigrate(&models.Role{}, &models.RoleMenu{}, &models.Customer{}, &models.KYCData{}, &models.KYCVerificationHistory{})

	// 插入初始数据
	role := models.Role{
		RoleName:    "lottery_admin",
		RoleType:    "admin",
		Description: "Administrator for lottery management",
	}
	db.DB.Create(&role)

	roleMenu := models.RoleMenu{
		RoleID:   role.RoleID,
		MenuName: "lottery_management",
		MenuPath: "/lottery/manage",
	}
	db.DB.Create(&roleMenu)

	customer := models.Customer{
		CustomerAddress:  "0xTestAddress123",
		RoleID:           role.RoleID,
		RegistrationTime: time.Now(),
		AssignedDate:     time.Now(),
		KYCData: models.KYCData{
			CustomerAddress: "0xTestAddress123",
			Email:           "test@example.com",
			Name:            "Test User",
		},
	}
	db.DB.Create(&customer)

	return &TestSuite{DB: db.DB}
}

func (suite *TestSuite) TearDown() {
	// 清理数据库
	suite.DB.Exec("DROP TABLE IF EXISTS kyc_verification_history CASCADE;")
	suite.DB.Exec("DROP TABLE IF EXISTS kyc_data CASCADE;")
	suite.DB.Exec("DROP TABLE IF EXISTS customers CASCADE;")
	suite.DB.Exec("DROP TABLE IF EXISTS role_menus CASCADE;")
	suite.DB.Exec("DROP TABLE IF EXISTS roles CASCADE;")

	// 删除测试数据库
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	testDBName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable", dbHost, dbPort, dbUser, dbPassword)
	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to postgres database for cleanup: %v", err)
	}
	defer sqlDB.Close()

	_, err = sqlDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBName))
	if err != nil {
		log.Fatalf("Failed to drop test database %s: %v", testDBName, err)
	}
}
