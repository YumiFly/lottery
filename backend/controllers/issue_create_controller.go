package controllers

import (
	"net/http"
	"time"

	"backend/db"
	"backend/models"
	issueCreateService "backend/services/issue"
	"backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// CreateIssueRequest 定义创建期号的请求结构
type CreateIssueRequest struct {
	LotteryID      string    `json:"lottery_id" validate:"required,max=50"`
	IssueNumber    string    `json:"issue_number" validate:"required,max=50"`
	SaleEndTime    time.Time `json:"sale_end_time" validate:"required"`
	DrawTime       time.Time `json:"draw_time" validate:"required,gtfield=SaleEndTime"`
	Status         string    `json:"status" validate:"required,oneof= PENDING DRAWN"`
	PrizePool      float64   `json:"prize_pool" validate:"gte=0"`
	WinningNumbers string    `json:"winning_numbers" validate:"omitempty,max=100"`
	RandomSeed     string    `json:"random_seed" validate:"omitempty,max=100"`
	DrawTxHash     string    `json:"draw_tx_hash" validate:"omitempty,max=66"`
}

// CreateIssueResponse 定义创建期号的响应结构，包含交易哈希
type CreateIssueResponse struct {
	Issue  models.LotteryIssue `json:"issue"`
	TxHash string              `json:"tx_hash"`
}

// CreateIssue 处理 POST /lottery/issues 请求
// swagger:route POST /lottery/issues lottery createIssue
//
// 请求体:
//   - lottery_id: 彩票 ID（必填，最大 50 字符）
//   - issue_number: 期号编号（必填，最大 50 字符）
//   - sale_end_time: 销售截止时间（必填，ISO 8601 格式）
//   - draw_time: 开奖时间（必填，晚于 sale_end_time）
//   - status: 状态（必填，pending 或 drawn）
//   - prize_pool: 奖池金额（选填，非负）
//   - winning_numbers: 中奖号码（选填，最大 100 字符）
//   - draw_tx_hash: 开奖交易哈希（选填，最大 66 字符）
//   - random_seed: 随机种子（选填，最大 100 字符）
//
// 响应:
//   - 201: 成功，返回 Response{Message, Code, Data}，Data 为 CreateIssueResponse
//   - 400: 无效参数
//   - 500: 服务器错误
func NewLotteryIssue(c *gin.Context) {
	// 绑定请求体
	var req CreateIssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Warn("Failed to bind request body", "error", err)
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.NewBadRequestError("Invalid request body", err)))
		return
	}

	// 验证参数
	validate := validator.New()
	if err := validate.Struct(&req); err != nil {
		utils.Logger.Warn("Failed to validate request parameters", "error", err)
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.NewBadRequestError("Parameter validation failed", err)))
		return
	}

	// 调用 service 层
	issueService := issueCreateService.NewIssueCreateService(db.DB)
	issue, txhash, err := issueService.CreateIssue(c.Request.Context(), issueCreateService.CreateIssueParams{
		LotteryID:      req.LotteryID,
		IssueNumber:    req.IssueNumber,
		SaleEndTime:    req.SaleEndTime,
		DrawTime:       req.DrawTime,
		Status:         req.Status,
		PrizePool:      req.PrizePool,
		WinningNumbers: req.WinningNumbers,
		RandomSeed:     req.RandomSeed,
		DrawTxHash:     req.DrawTxHash,
	})
	if err != nil {
		utils.Logger.Error("Failed to create issue", "error", err)
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(err))
		return
	}

	// 记录日志
	utils.Logger.Info("Successfully created issue",
		"lottery_id", req.LotteryID,
		"issue_number", req.IssueNumber,
		"status", req.Status,
		"issue_id", issue.IssueID,
		"txhash", txhash.Hex())

	// 返回响应
	response := CreateIssueResponse{
		Issue:  *issue,
		TxHash: txhash.Hex(),
	}
	c.JSON(http.StatusCreated, utils.SuccessResponse("Successfully created issue", response))
}
