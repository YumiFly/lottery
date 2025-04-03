package services

import (
	"backend/db"
	"backend/models"
	"backend/utils"
	"time"
)

// CreateLotteryType 创建彩票类型
func CreateLotteryType(lotteryType *models.LotteryType) error {
	utils.Logger.Info("Creating lottery type", "type_id", lotteryType.TypeID, "type_name", lotteryType.TypeName)
	lotteryType.CreatedAt = time.Now()
	lotteryType.UpdatedAt = time.Now()
	if err := db.DB.Create(lotteryType).Error; err != nil {
		utils.Logger.Error("Failed to create lottery type", "error", err)
		return utils.NewServiceError("failed to create lottery type", err)
	}
	utils.Logger.Info("Lottery type created successfully", "type_id", lotteryType.TypeID)
	return nil
}

// GetAllLotteryTypes 获取所有彩票类型
func GetAllLotteryTypes() ([]models.LotteryType, error) {
	utils.Logger.Info("Fetching all lottery types")
	var lotteryTypes []models.LotteryType
	if err := db.DB.Find(&lotteryTypes).Error; err != nil {
		utils.Logger.Error("Failed to fetch lottery types", "error", err)
		return nil, utils.NewServiceError("failed to fetch lottery types", err)
	}
	utils.Logger.Info("Fetched lottery types", "count", len(lotteryTypes))
	return lotteryTypes, nil
}
