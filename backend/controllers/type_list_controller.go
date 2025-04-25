package controllers

import (
	"backend/db"
	"backend/utils"
	"net/http"

	typeListService "backend/services/lttype"

	"github.com/gin-gonic/gin"
)

// GetAllLotteryTypes handles GET /lottery/types requests
// swagger:route GET /lottery/types lottery getAllLotteryTypes
//
// Responses:
//   - 200: Success, returns Response{Message, Code, Data}, Data is []LotteryType
//   - 500: Server error
func ListLotteryTypes(c *gin.Context) {
	// Call service layer
	lotteryService := typeListService.NewTypeListService(db.DB)
	lotteryTypes, err := lotteryService.GetAllLotteryTypes(c.Request.Context())
	if err != nil {
		utils.Logger.Error("Failed to get lottery types", "error", err)
		c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(err))
		return
	}

	// Log success
	utils.Logger.Info("Successfully retrieved lottery types",
		"count", len(lotteryTypes))

	// Return response
	c.JSON(http.StatusOK, utils.SuccessResponse("Lottery types retrieved successfully", lotteryTypes))
}
