package controllers

import (
	"net/http"

	"backend/db"
	"backend/models"
	winnerListService "backend/services/winner"
	"backend/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// GetAllWinnersResponse defines the response structure for paginated winner data
type GetAllWinnersResponse struct {
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
	Winners  []models.Winner `json:"winners"`
}

// GetAllWinnersQuery defines the query parameters for fetching winners
type GetAllWinnersQuery struct {
	IssueID    string `form:"issue_id" validate:"omitempty,max=50"`
	Address    string `form:"address" validate:"omitempty,len=42,hexadecimal,startsWith=0x"`
	PrizeLevel string `form:"prize_level" validate:"omitempty,max=50"`
	Page       int    `form:"page" validate:"omitempty,min=1"`
	PageSize   int    `form:"page_size" validate:"omitempty,min=1,max=100"`
}

// ListWinners handles GET /lottery/winners requests
// swagger:route GET /lottery/winners lottery getAllWinners
//
// Query parameters:
//   - issue_id: Lottery issue ID (optional, max 50 characters)
//   - address: Winner's Ethereum address (optional, 42-character hex starting with 0x)
//   - prize_level: Prize level (optional, max 50 characters)
//   - page: Page number, default 1 (optional)
//   - page_size: Records per page, default 20, max 100 (optional)
//
// Responses:
//   - 200: Success, returns Response{Message, Code, Data}, Data is GetAllWinnersResponse
//   - 400: Invalid query parameters
//   - 500: Server error
func ListWinners(c *gin.Context) {
	// Bind query parameters
	var query GetAllWinnersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.Logger.Warn("Failed to bind query parameters", "error", err)
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.NewBadRequestError("Invalid query parameters", err)))
		return
	}

	// Set default values
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}

	// Validate parameters
	validate := validator.New()
	validate.RegisterValidation("hexadecimal", func(fl validator.FieldLevel) bool {
		return common.IsHexAddress(fl.Field().String())
	})
	validate.RegisterValidation("startsWith", func(fl validator.FieldLevel) bool {
		return fl.Field().String()[:2] == "0x"
	})
	if err := validate.Struct(&query); err != nil {
		utils.Logger.Warn("Failed to validate query parameters", "error", err)
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.NewBadRequestError("Invalid query parameters", err)))
		return
	}

	// Call service layer
	winnerService := winnerListService.NewWinnerListService(db.DB)
	result, err := winnerService.GetAllWinners(c.Request.Context(), winnerListService.WinnerQueryParams{
		IssueID:    query.IssueID,
		Address:    query.Address,
		PrizeLevel: query.PrizeLevel,
		Page:       query.Page,
		PageSize:   query.PageSize,
	})
	if err != nil {
		utils.Logger.Error("Failed to query winners", "error", err)
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(err))
		return
	}

	// Log success
	utils.Logger.Info("Successfully queried winners",
		"issue_id", query.IssueID,
		"address", query.Address,
		"prize_level", query.PrizeLevel,
		"page", query.Page,
		"page_size", query.PageSize,
		"total", result.Total,
		"returned", len(result.Winners))

	// Return response
	c.JSON(http.StatusOK, utils.SuccessResponse("Successfully queried winners", result))
}
