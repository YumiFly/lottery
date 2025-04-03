package db

import (
	"backend/config"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// 加载配置
	config.LoadConfig()

	// 构造 DSN
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBName,
		config.AppConfig.DB_SSLMODE,
		config.AppConfig.DB_TIMEZONE,
	)

	// 打印 DSN（隐藏密码）以便调试
	safeDSN := fmt.Sprintf("host=%s port=%s user=%s password=**** dbname=%s sslmode=%s TimeZone=%s",
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBUser,
		config.AppConfig.DBName,
		config.AppConfig.DB_SSLMODE,
		config.AppConfig.DB_TIMEZONE,
	)
	log.Printf("Attempting to connect to database with DSN: %s", safeDSN)

	// 尝试连接数据库
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		// 可选：启用 GORM 的详细日志
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 获取底层的 sql.DB 并配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大存活时间

	// 测试连接是否正常
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Database connected successfully")

	// 可选：自动迁移表结构（取消注释以启用）
	// err = DB.AutoMigrate(&models.Role{}, &models.RoleMenu{}, &models.Customer{},
	// 	&models.KYCData{}, &models.KYCVerificationHistory{},
	// 	&models.LotteryType{}, &models.Lottery{}, &models.LotteryIssue{},
	// 	&models.LotteryTicket{}, &models.Winner{})
	// if err != nil {
	// 	log.Fatalf("Failed to auto-migrate database: %v", err)
	// }
	// log.Println("Database schema migrated successfully")
}
