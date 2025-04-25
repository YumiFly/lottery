package controllers

import (
	"backend/blockchain"
	"backend/db"
	"backend/services/lottery"

	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type LotteryDrawRequest struct {
	IssueID string `json:"issue_id" validate:"required,max=36"`
}

func NewDrawLottery(c *gin.Context) {
	// 检查权限（确保调用者是 lottery_admin）
	// role, exists := c.Get("role")
	// if !exists || (role != "lottery_admin" && role != "admin") {
	// 	c.JSON(http.StatusForbidden, utils.ErrorResponse(utils.ErrCodeForbidden, "Insufficient permissions", nil))
	// 	return
	// }

	var req LotteryDrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Warn("Failed to bind request body", "error", err)
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.NewBadRequestError("Invalid request body", err)))
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		utils.Logger.Warn("Validation failed", "error", err)
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.NewBadRequestError("Validation failed", err)))
		return
	}
	service := lottery.NewLotteryDrawService(blockchain.Client, blockchain.Auth, db.DB)

	if err := service.DrawLotteryAsync(req.IssueID); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCodeInternalServer, "Failed to draw lottery", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Lottery draw initiated", nil))
}
