package lttype

import (
	"backend/models"
	"backend/utils"
	"context"

	"gorm.io/gorm"
)

// LotteryService encapsulates lottery-related business logic
type TypeListService struct {
	db *gorm.DB
}

// NewLotteryService creates a new LotteryService instance
func NewTypeListService(db *gorm.DB) *TypeListService {
	return &TypeListService{db: db}
}

// GetAllLotteryTypes retrieves all lottery types
//
// Parameters:
//   - ctx: Request context
//
// Returns:
//   - []LotteryType: List of all lottery types
//   - error: Retrieval error
func (s *TypeListService) GetAllLotteryTypes(ctx context.Context) ([]models.LotteryType, error) {
	// Log retrieval attempt
	utils.Logger.Info("Fetching all lottery types")

	// Query lottery types
	var lotteryTypes []models.LotteryType
	if err := s.db.WithContext(ctx).Find(&lotteryTypes).Error; err != nil {
		utils.Logger.Error("Failed to fetch lottery types", "error", err)
		return nil, utils.NewInternalError("Failed to fetch lottery types", err)
	}

	// Log success
	utils.Logger.Info("Fetched lottery types successfully",
		"count", len(lotteryTypes))

	return lotteryTypes, nil
}
