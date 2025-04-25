package controllers

import (
	"backend/db"
	"backend/models"
	"backend/utils"
	"fmt"
	"net/http"
	"strings"

	ticketPurchaseService "backend/services/ticket"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// CreateTicketRequest defines the request body
type PurchaseTicketRequest struct {
	IssueID        string `json:"issue_id" validate:"required,max=36"`
	BuyerAddress   string `json:"buyer_address" validate:"required,hexadecimal,startsWith=0x,len=42"`
	PurchaseAmount uint64 `json:"purchase_amount" validate:"required,gt=0,lte=1000"`
	BetContent     string `json:"bet_content" validate:"required,max=100"`
}

// PurchaseTicketResponse defines the response structure for purchasing a ticket
type PurchaseTicketResponse struct {
	Ticket models.LotteryTicket `json:"ticket"`
	TxHash string               `json:"tx_hash"`
}

// PurchaseTicket handles POST /lottery/tickets requests
func NewPurchaseTicket(c *gin.Context) {
	var req PurchaseTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Warn("Failed to bind request body", "error", err)
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.NewBadRequestError("Invalid request body", err)))
		return
	}

	validate := validator.New()
	validate.RegisterValidation("hexadecimal", func(fl validator.FieldLevel) bool {
		return common.IsHexAddress(fl.Field().String())
	})
	validate.RegisterValidation("startsWith", func(fl validator.FieldLevel) bool {
		return strings.HasPrefix(fl.Field().String(), "0x")
	})
	if err := validate.Struct(&req); err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("Field %s: %s", err.Field(), err.Tag()))
		}
		utils.Logger.Warn("Failed to validate request body", "errors", errors, "request", req)
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.NewBadRequestError("Parameter validation failed: "+strings.Join(errors, ", "), err)))
		return
	}

	ticketService := ticketPurchaseService.NewTicketPurchaseService(db.DB)
	ticket, txHash, err := ticketService.PurchaseTicket(c.Request.Context(), ticketPurchaseService.PurchaseTicketParams{
		IssueID:        req.IssueID,
		BuyerAddress:   req.BuyerAddress,
		PurchaseAmount: req.PurchaseAmount,
		BetContent:     req.BetContent,
		TicketID:       uuid.NewString(),
	})
	if err != nil {
		utils.Logger.Error("Failed to buy ticket",
			"issue_id", req.IssueID,
			"buyer_address", req.BuyerAddress,
			"purchase_amount", req.PurchaseAmount,
			"txHash", txHash,
			"error", err)
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(err))
		return
	}

	utils.Logger.Info("Successfully purchased ticket",
		"ticket_id", ticket.TicketID,
		"issue_id", ticket.IssueID,
		"buyer_address", ticket.BuyerAddress)

	response := PurchaseTicketResponse{
		Ticket: *ticket,
		TxHash: txHash.Hex(),
	}
	c.JSON(http.StatusCreated, utils.SuccessResponse("Ticket purchased successfully", response))
}
