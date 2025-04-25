package issue

import (
	"backend/models"
	"backend/utils"
	"context"
	"time"

	"gorm.io/gorm"
)

// IssuePoolService 奖池服务
type IssuePoolService struct {
	db *gorm.DB
}

func NewIssuePoolService(db *gorm.DB) *IssuePoolService {
	return &IssuePoolService{db: db}
}

// GetAllPools 获取所有奖池总额
func (s *IssuePoolService) CountIssuePools(ctx context.Context) (int64, error) {
	utils.Logger.Info("Fetching all pools")
	var issues []models.LotteryIssue
	if err := s.db.WithContext(ctx).Where("sale_end_time > ?", time.Now()).Find(&issues).Error; err != nil {
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
