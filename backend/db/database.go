// db/database.go
package db

import (
	"backend/config"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	config.LoadConfig() // 调用加载配置，无返回值
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		config.AppConfig.DBHost, config.AppConfig.DBPort, config.AppConfig.DBUser,
		config.AppConfig.DBPassword, config.AppConfig.DBName, config.AppConfig.DB_SSLMODE,
		config.AppConfig.DB_TIMEZONE)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	DB = gormDB
}
