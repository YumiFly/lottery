package controllers

import (
	"net/http"

	"backend/db"
	"backend/models"
	lotteryListService "backend/services/lottery"
	"backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// GetAllLotteryResponse 定义彩票列表的响应结构
type GetAllLotteryResponse struct {
	Lotteries []models.Lottery `json:"lotteries"`
	Total     int64            `json:"total"`
}

// GetAllLotteryQuery 定义查询彩票的请求结构
type GetAllLotteryQuery struct {
	TypeID     string `form:"type_id" validate:"omitempty,max=50"`
	TicketName string `form:"ticket_name" validate:"omitempty,max=255"`
}

// GetAllLottery 处理 GET /lottery/lottery 请求
// swagger:route GET /lottery/lottery lottery getAllLottery
// ...
func ListAllLotteries(c *gin.Context) {
	// 绑定查询参数
	var query GetAllLotteryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.Logger.Warn("binding query failed", "error", err)
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.NewBadRequestError("invalid bind parameters", err)))
		return
	}

	// 验证参数
	validate := validator.New()
	if err := validate.Struct(&query); err != nil {
		utils.Logger.Warn("invalid query parameters", "error", err)
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.NewBadRequestError("invalid query parameters", err)))
		return
	}

	// 调用 service 层
	lListService := lotteryListService.NewLotteryListService(db.DB)
	result, err := lListService.GetAllLotteries(c.Request.Context(), lotteryListService.LotteryQueryParams{
		TypeID:     query.TypeID,
		TicketName: query.TicketName,
	})
	if err != nil {
		utils.Logger.Error("get lottery list failed", "error", err)
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(err))
		return
	}

	// 记录日志
	utils.Logger.Info("get lottery list success",
		"type_id", query.TypeID,
		"ticket_name", query.TicketName,
		"total", result.Total,
		"returned", len(result.Lotteries))

	// 返回响应
	c.JSON(http.StatusOK, utils.SuccessResponse("get lottery list success", result))
}
