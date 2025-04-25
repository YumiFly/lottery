package issue

import (
	"context"
	"fmt"

	"backend/models"

	"gorm.io/gorm"
)

// IssueListService 封装期号列表查询的业务逻辑
type IssueListService struct {
	db *gorm.DB
}

// NewIssueListService 创建 IssueListService 实例
func NewIssueListService(db *gorm.DB) *IssueListService {
	return &IssueListService{db: db}
}

// IssueQueryParams 定义查询期号的参数结构
type IssueQueryParams struct {
	LotteryID   string // 彩票 ID
	Status      string // 期号状态
	IssueNumber string // 期号编号
	Page        int    // 页码
	PageSize    int    // 每页大小
}

// GetAllIssuesResponse 定义分页期号数据的响应结构
type GetAllIssuesResponse struct {
	Total    int64                 `json:"total"`     // 总记录数
	Page     int                   `json:"page"`      // 当前页码
	PageSize int                   `json:"page_size"` // 每页大小
	Issues   []models.LotteryIssue `json:"issues"`    // 期号列表
}

// validateIssueParams 验证查询参数
func (s *IssueListService) validateIssueParams(params IssueQueryParams) error {
	// 验证页码
	if params.Page < 1 {
		return fmt.Errorf("invalid page number: %d", params.Page)
	}

	// 验证每页大小
	if params.PageSize < 1 || params.PageSize > 100 {
		return fmt.Errorf("invalid page size: %d, should be between 1 and 100", params.PageSize)
	}

	// 验证状态值
	validStatuses := map[string]bool{
		models.IssueStatusPending: true,
		models.IssueStatusDrawn:   true,

		// 可根据 models.IssueStatusDrawn 等补充
	}
	if params.Status != "" && !validStatuses[params.Status] {
		return fmt.Errorf("invalid status: %s", params.Status)
	}

	// 验证期号编号长度
	if len(params.IssueNumber) > 50 {
		return fmt.Errorf("issue number too long: %s", params.IssueNumber)
	}

	return nil
}

// buildIssueQuery 根据参数构建 GORM 查询
func (s *IssueListService) buildIssueQuery(query *gorm.DB, params IssueQueryParams) *gorm.DB {
	if params.LotteryID != "" {
		query = query.Where("lottery_id = ?", params.LotteryID)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.IssueNumber != "" {
		// 使用 PostgreSQL 的 ILIKE 进行大小写不敏感模糊查询
		query = query.Where("issue_number ILIKE ?", "%"+params.IssueNumber+"%")
	}
	return query
}

// GetAllIssues 查询彩票期号列表，支持分页和多条件过滤
//
// 参数:
//   - ctx: 请求上下文，用于超时控制
//   - params: 查询参数，包括 lottery_id, status, issue_number, page, page_size
//
// 返回:
//   - *GetAllIssuesResponse: 分页结果，包含总记录数、页码、每页大小和期号列表
//   - error: 查询错误或参数无效
func (s *IssueListService) GetAllIssues(ctx context.Context, params IssueQueryParams) (*GetAllIssuesResponse, error) {
	// 验证参数
	if err := s.validateIssueParams(params); err != nil {
		return nil, err
	}

	// 构建查询
	query := s.db.WithContext(ctx).Model(&models.LotteryIssue{})
	query = s.buildIssueQuery(query, params)

	// 查询总记录数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("query total count failed: %w", err)
	}

	// 查询分页数据
	var issues []models.LotteryIssue
	offset := (params.Page - 1) * params.PageSize
	if err := query.
		Preload("Lottery"). // 预加载关联的 Lottery 数据
		Order("created_at DESC").
		Offset(offset).
		Limit(params.PageSize).
		Find(&issues).Error; err != nil {
		return nil, fmt.Errorf("query issues failed: %w", err)
	}

	// 构造响应
	return &GetAllIssuesResponse{
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
		Issues:   issues,
	}, nil
}
