package controllers

import (
	"net/http"

	"backend/db"
	"backend/models"
	lotteryCreateService "backend/services/lottery"
	"backend/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// CreateLotteryRequest defines the request structure for creating a lottery
type CreateLotteryRequest struct {
	TypeID                 string  `json:"type_id" validate:"required,max=36"`
	TicketName             string  `json:"ticket_name" validate:"required,max=100"`
	TicketSupply           int64   `json:"ticket_supply" validate:"required,gt=0"`
	TicketPrice            float64 `json:"ticket_price" validate:"required,gt=0"`
	BettingRules           string  `json:"betting_rules" validate:"required"`
	PrizeStructure         string  `json:"prize_structure" validate:"required"`
	RegisteredAddr         string  `json:"registered_addr" validate:"required,len=42,eth_addr"`
	RolloutContractAddress string  `json:"rollout_contract_address" validate:"required,len=42,eth_addr"`
}

// CreateLotteryResponse defines the response structure, including the lottery and transaction hash
type CreateLotteryResponse struct {
	Lottery models.Lottery `json:"lottery"`
	TxHash  string         `json:"tx_hash"`
}

// CreateLottery handles POST /lottery/lottery requests
// swagger:route POST /lottery/lottery lottery createLottery
//
// Request body:
//   - type_id: Lottery type ID (required, max 36 characters)
//   - ticket_name: Ticket name (required, max 100 characters)
//   - ticket_supply: Total ticket supply (required, positive integer)
//   - ticket_price: Ticket price (required, positive float)
//   - registered_addr: Owner Ethereum address (required, 42-character hex)
//   - rollout_contract_address: Rollout contract address (required, 42-character hex)
//
// Responses:
//   - 201: Success, returns Response{Message, Code, Data}, Data is CreateLotteryResponse
//   - 400: Invalid input
//   - 403: Insufficient permissions
//   - 500: Server error
func NewLottery(c *gin.Context) {
	// 检查权限（确保调用者是 lottery_admin）
	// role, exists := c.Get("role")
	// if !exists || (role != "lottery_admin" && role != "admin") {
	// 	c.JSON(http.StatusForbidden, utils.ErrorResponse(utils.NewBadRequestError("Insufficient permissions", nil)))
	// 	return
	// }

	// Bind request body
	var req CreateLotteryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Warn("Failed to bind request body", "error", err)
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.NewBadRequestError("Invalid request body", err)))
		return
	}

	// Validate parameters
	validate := validator.New()
	validate.RegisterValidation("eth_addr", func(fl validator.FieldLevel) bool {
		addr := fl.Field().String()
		return common.IsHexAddress(addr)
	})
	if err := validate.Struct(&req); err != nil {
		utils.Logger.Warn("Failed to validate request parameters", "error", err)
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.NewBadRequestError("Parameter validation failed", err)))
		return
	}

	// Call service layer
	lCreateService := lotteryCreateService.NewLotteryCreateService(db.DB)
	lottery, txHash, err := lCreateService.CreateLottery(c.Request.Context(), lotteryCreateService.CreateLotteryParams{
		TypeID:                 req.TypeID,
		TicketName:             req.TicketName,
		TicketSupply:           req.TicketSupply,
		TicketPrice:            req.TicketPrice,
		BettingRules:           req.BettingRules,
		PrizeStructure:         req.PrizeStructure,
		RegisteredAddr:         req.RegisteredAddr,
		RolloutContractAddress: req.RolloutContractAddress,
	})
	if err != nil {
		utils.Logger.Error("Failed to create lottery", "error", err)
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(err))
		return
	}

	// Log success
	utils.Logger.Info("Successfully created lottery",
		"lottery_id", lottery.LotteryID,
		"type_id", lottery.TypeID,
		"contract_address", lottery.ContractAddress,
		"tx_hash", txHash)

	// Return response
	response := CreateLotteryResponse{
		Lottery: *lottery,
		TxHash:  txHash.Hex(),
	}
	c.JSON(http.StatusCreated, utils.SuccessResponse("Successfully created lottery", response))
}
