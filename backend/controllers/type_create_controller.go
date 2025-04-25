package controllers

import (
	"backend/db"
	"backend/models"
	typeCreateService "backend/services/lttype"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// CreateLotteryTypeRequest defines the request structure for creating a lottery type
type CreateLotteryTypeRequest struct {
	TypeName    string `json:"type_name" validate:"required,max=100"`
	Description string `json:"description" validate:"omitempty,max=100"`
}

// CreateLotteryTypeResponse defines the response structure for creating a lottery type
type CreateLotteryTypeResponse struct {
	LotteryType models.LotteryType `json:"lottery_type"`
}

// CreateLotteryType handles POST /lottery/types requests
// swagger:route POST /lottery/types lottery createLotteryType
//
// Request body:
//   - name: Lottery type name (required, max 100 characters)
//
// Responses:
//   - 201: Success, returns Response{Message, Code, Data}, Data is CreateLotteryTypeResponse
//   - 400: Invalid input
//   - 403: Insufficient permissions
//   - 500: Server error
func NewLotteryType(c *gin.Context) {
	// Check permissions (ensure caller is lottery_admin)
	// role, exists := c.Get("role")
	// if !exists || role != "lottery_admin" {
	// 	utils.Logger.Warn("Insufficient permissions", "role", role)
	// 	c.JSON(http.StatusForbidden, utils.NewErrorResponse(utils.NewBadRequestError("Insufficient permissions", nil)))
	// 	return
	// }

	// Bind request body
	var req CreateLotteryTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Warn("Failed to bind request body", "error", err)
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.NewBadRequestError("Invalid request body", err)))
		return
	}

	// Validate parameters
	validate := validator.New()
	if err := validate.Struct(&req); err != nil {
		utils.Logger.Warn("Failed to validate request body", "error", err)
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.NewBadRequestError("Invalid request parameters", err)))
		return
	}

	// Call service layer
	lotteryService := typeCreateService.NewTypeCreateService(db.DB)
	lotteryType, err := lotteryService.CreateLotteryType(c.Request.Context(), typeCreateService.CreateLotteryTypeParams{
		TypeName:    req.TypeName,
		Description: req.Description,
	})
	if err != nil {
		utils.Logger.Error("Failed to create lottery type", "error", err)
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(err))
		return
	}

	// Log success
	utils.Logger.Info("Successfully created lottery type",
		"type_id", lotteryType.TypeID,
		"name", lotteryType.TypeName)

	// Return response
	response := CreateLotteryTypeResponse{
		LotteryType: *lotteryType,
	}
	c.JSON(http.StatusCreated, utils.SuccessResponse("Lottery type created successfully", response))
}
