package controllers

import (
	"backend/db"
	"backend/models"
	issueListService "backend/services/issue"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// GetAllIssuesResponse 定义分页期号数据的响应结构
type GetAllIssuesResponse struct {
	Total    int64                 `json:"total"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"page_size"`
	Issues   []models.LotteryIssue `json:"issues"`
}

// GetAllIssuesQuery 定义查询期号的请求结构
type GetAllIssuesQuery struct {
	LotteryID   string `form:"lottery_id" validate:"omitempty,max=50"`
	Status      string `form:"status" validate:"omitempty,oneof=PENDING DRAWN DRAWING"`
	IssueNumber string `form:"issue_number" validate:"omitempty,max=50"`
	Page        int    `form:"page" validate:"min=1"`
	PageSize    int    `form:"page_size" validate:"min=1,max=100"`
}

// GetAllIssues 处理 GET /lottery/issues 请求，支持多参数查询彩票期号
//
// 路由: GET /lottery/issues
//
// 查询参数:
//   - lottery_id: 按彩票 ID 精确过滤（可选）
//   - status: 按期号状态精确过滤，如 "pending", "drawn"（可选）
//   - issue_number: 按期号编号模糊查询（可选）
//   - page: 分页页码，默认 1（可选）
//   - page_size: 每页记录数，默认 20，最大 100（可选）
//
// 响应:
//   - 200: 成功，返回 Response{Message, Code, Data}，Data 为 GetAllIssuesResponse
//   - 400: 无效参数
//   - 500: 服务器错误
func ListAllIssues(c *gin.Context) {
	// 绑定查询参数
	var query GetAllIssuesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.Logger.Warn("bind query failed", "error", err)
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.NewBadRequestError("invalid query parameters", err)))
		return
	}

	// Set default values
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}

	// 验证参数
	validate := validator.New()
	if err := validate.Struct(&query); err != nil {
		utils.Logger.Warn("invalid query parameters", "error", err)
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.NewBadRequestError("inivalid query parameters", err)))
		return
	}

	// 调用 service 层
	issueService := issueListService.NewIssueListService(db.DB)
	result, err := issueService.GetAllIssues(c.Request.Context(), issueListService.IssueQueryParams{
		LotteryID:   query.LotteryID,
		Status:      query.Status,
		IssueNumber: query.IssueNumber,
		Page:        query.Page,
		PageSize:    query.PageSize,
	})
	if err != nil {
		utils.Logger.Error("query all issues failed", "error", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// 记录日志
	utils.Logger.Info("query all issues success",
		"lottery_id", query.LotteryID,
		"status", query.Status,
		"issue_number", query.IssueNumber,
		"page", query.Page,
		"page_size", query.PageSize,
		"total", result.Total,
		"returned", len(result.Issues))

	// 返回响应
	c.JSON(http.StatusOK, utils.SuccessResponse("query all issues success", result))
}
