package controllers

import (
	"backend/db"
	"backend/models"
	ticketListSrevice "backend/services/ticket"
	"backend/utils"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// GetAllTicketsResponse defines the response structure for paginated ticket data
type GetAllTicketsResponse struct {
	Total    int64                  `json:"total"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
	Tickets  []models.LotteryTicket `json:"tickets"`
}

// GetAllTicketsQuery defines the query parameters for fetching tickets
type GetAllTicketsQuery struct {
	IssueID      string `form:"issue_id" validate:"omitempty,max=50"`
	BuyerAddress string `form:"buyer_address" validate:"omitempty,len=42,hexadecimal,startsWith=0x"`
	TicketID     string `form:"ticket_id" validate:"omitempty,max=50"`
	Page         int    `form:"page" validate:"omitempty,min=1"`
	PageSize     int    `form:"page_size" validate:"omitempty,min=1,max=100"`
}

// ListPurchasedTickets handles GET /lottery/tickets requests
// swagger:route GET /lottery/tickets lottery listPurchasedTickets
//
// Query parameters:
//   - issue_id: Lottery issue ID (optional, max 50 characters)
//   - buyer_address: Buyer's Ethereum address (optional, 42-character hex starting with 0x)
//   - ticket_id: Ticket ID (optional, max 50 characters)
//   - page: Page number, default 1 (optional)
//   - page_size: Records per page, default 20, max 100 (optional)
//
// Responses:
//   - 200: Success, returns Response{Message, Code, Data}, Data is GetAllTicketsResponse
//   - 400: Invalid query parameters
//   - 500: Server error
func ListPurchasedTickets(c *gin.Context) {
	// Bind query parameters
	var query GetAllTicketsQuery
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
	ticketService := ticketListSrevice.NewTicketListService(db.DB)
	result, err := ticketService.GetAllTickets(c.Request.Context(), ticketListSrevice.TicketQueryParams{
		IssueID:      query.IssueID,
		BuyerAddress: query.BuyerAddress,
		TicketID:     query.TicketID,
		Page:         query.Page,
		PageSize:     query.PageSize,
	})
	if err != nil {
		utils.Logger.Error("Failed to query purchased tickets", "error", err)
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(err))
		return
	}

	// Log success
	utils.Logger.Info("Successfully queried purchased tickets",
		"issue_id", query.IssueID,
		"buyer_address", query.BuyerAddress,
		"ticket_id", query.TicketID,
		"page", query.Page,
		"page_size", query.PageSize,
		"total", result.Total,
		"returned", len(result.Tickets))

	// Return response
	c.JSON(http.StatusOK, utils.SuccessResponse("Successfully queried purchased tickets", result))
}
