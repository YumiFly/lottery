package ticket

import (
	"backend/models"
	"backend/utils"
	"context"

	"gorm.io/gorm"
)

// TicketService encapsulates ticket purchasing business logic
type TicketListService struct {
	db *gorm.DB
}

// NewTicketService creates a new TicketService instance
func NewTicketListService(db *gorm.DB) *TicketListService {
	return &TicketListService{db: db}
}

// TicketQueryParams defines the query parameters for fetching tickets
type TicketQueryParams struct {
	IssueID      string
	BuyerAddress string
	TicketID     string
	Page         int
	PageSize     int
}

// TicketListResult defines the result structure for ticket queries
type TicketListResult struct {
	Total    int64
	Page     int
	PageSize int
	Tickets  []models.LotteryTicket
}

// GetAllTickets retrieves all purchased tickets with optional filters and pagination
func (s *TicketListService) GetAllTickets(ctx context.Context, params TicketQueryParams) (*TicketListResult, error) {
	var tickets []models.LotteryTicket
	query := s.db.WithContext(ctx).Model(&models.LotteryTicket{})

	// Apply filters
	if params.IssueID != "" {
		query = query.Where("issue_id = ?", params.IssueID)
	}
	if params.BuyerAddress != "" {
		query = query.Where("buyer_address = ?", params.BuyerAddress)
	}
	if params.TicketID != "" {
		query = query.Where("ticket_id = ?", params.TicketID)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		utils.Logger.Error("Failed to count tickets", "error", err)
		return nil, utils.NewInternalError("Failed to count tickets", err)
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

	// Fetch tickets
	if err := query.Preload("LotteryIssue").Preload("LotteryIssue.Lottery").Find(&tickets).Error; err != nil {
		utils.Logger.Error("Failed to fetch tickets", "error", err)
		return nil, utils.NewInternalError("Failed to fetch tickets", err)
	}

	return &TicketListResult{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Tickets:  tickets,
	}, nil
}
