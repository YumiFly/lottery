package winner

import (
	"backend/models"
	"backend/utils"
	"context"

	"gorm.io/gorm"
)

// WinnerListService encapsulates winner list business logic
type WinnerListService struct {
	db *gorm.DB
}

// NewWinnerListService creates a new WinnerListService instance
func NewWinnerListService(db *gorm.DB) *WinnerListService {
	return &WinnerListService{db: db}
}

// WinnerQueryParams defines the query parameters for fetching winners
type WinnerQueryParams struct {
	IssueID    string
	Address    string
	PrizeLevel string
	Page       int
	PageSize   int
}

// WinnerListResult defines the result structure for winner queries
type WinnerListResult struct {
	Total    int64
	Page     int
	PageSize int
	Winners  []models.Winner
}

// GetAllWinners retrieves all winners with optional filters and pagination
func (s *WinnerListService) GetAllWinners(ctx context.Context, params WinnerQueryParams) (*WinnerListResult, error) {
	var winners []models.Winner
	query := s.db.WithContext(ctx).Model(&models.Winner{})

	// Apply filters
	if params.IssueID != "" {
		query = query.Where("issue_id = ?", params.IssueID)
	}
	if params.Address != "" {
		query = query.Where("address = ?", params.Address)
	}
	if params.PrizeLevel != "" {
		query = query.Where("prize_level = ?", params.PrizeLevel)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		utils.Logger.Error("Failed to count winners", "error", err)
		return nil, utils.NewInternalError("Failed to count winners", err)
	}

	// Apply pagination
	page := params.Page
	if page < 1 {
		page = 1
	}
	pageSize := params.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Fetch winners
	if err := query.Preload("LotteryIssue").Preload("LotteryTicket").Order("created_at desc").Find(&winners).Error; err != nil {
		utils.Logger.Error("Failed to fetch winners", "error", err)
		return nil, utils.NewInternalError("Failed to fetch winners", err)
	}

	return &WinnerListResult{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Winners:  winners,
	}, nil
}
