package services

import (
	"backend/db"
	"backend/models"
	"backend/utils"
	"time"
)

func RemoveDuplicateString(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// GetAllPools 获取所有奖池总额
func GetAllPools() (int64, error) {
	utils.Logger.Info("Fetching all pools")
	var issues []models.LotteryIssue
	if err := db.DB.Where("sale_end_time > ?", time.Now()).Find(&issues).Error; err != nil {
		utils.Logger.Error("Failed to fetch issues for pools", "error", err)
		return 0, utils.NewServiceError("failed to fetch issues for pools", err)
	}

	totalPool := int64(0)
	for _, issue := range issues {
		issuePool := issue.PrizePool
		totalPool += int64(issuePool)
	}
	utils.Logger.Info("Fetched total pools", "total", totalPool)
	return totalPool, nil
}
