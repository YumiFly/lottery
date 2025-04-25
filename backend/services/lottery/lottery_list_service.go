package lottery

import (
	"context"

	"backend/models"
	"backend/services/common"
	"backend/utils"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// LotteryQueryParams 定义查询彩票的参数结构
type LotteryQueryParams struct {
	TypeID     string
	TicketName string
}

// GetAllLotteryResponse 定义彩票列表的响应结构
type GetAllLotteryResponse struct {
	Lotteries []models.Lottery `json:"lotteries"`
	Total     int64            `json:"total"`
}

// LotteryService 封装彩票相关业务逻辑
type LotteryListService struct {
	db *gorm.DB
}

// NewLotteryService 创建 LotteryService 实例
func NewLotteryListService(db *gorm.DB) *LotteryListService {
	return &LotteryListService{db: db}
}

// buildLotteryQueryOptions 构建彩票查询选项
func (s *LotteryListService) buildLotteryQueryOptions(params LotteryQueryParams) ([]common.QueryOption, error) {
	var options []common.QueryOption

	if params.TypeID != "" {
		if len(params.TypeID) > 50 {
			return nil, utils.NewBadRequestError("lottery type ID length cannot exceed 50 characters", nil)
		}
		options = append(options, func(q *gorm.DB) *gorm.DB {
			return q.Where("type_id = ?", params.TypeID)
		})
	}

	if params.TicketName != "" {
		if len(params.TicketName) > 255 {
			return nil, utils.NewBadRequestError("ticket name length cannot exceed 255 characters", nil)
		}
		options = append(options, func(q *gorm.DB) *gorm.DB {
			return q.Where("ticket_name ILIKE ?", "%"+params.TicketName+"%")
		})
	}

	return options, nil
}

// GetAllLottery 查询所有彩票信息
func (s *LotteryListService) GetAllLotteries(ctx context.Context, params LotteryQueryParams) (*GetAllLotteryResponse, error) {
	// 构建查询
	options, err := s.buildLotteryQueryOptions(params)
	if err != nil {
		return nil, err
	}
	query := s.db.WithContext(ctx).Model(&models.Lottery{})
	query = common.ApplyQueryOptions(query, options)

	// 查询总记录数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, utils.NewInternalError("count lottery failed", errors.Wrap(err, "database error"))
	}

	// 查询数据
	var lotteries []models.Lottery
	if err := query.
		Preload("LotteryType").
		Order("created_at DESC").
		Find(&lotteries).Error; err != nil {
		return nil, utils.NewInternalError("get lottery list failed", errors.Wrap(err, "database error"))
	}

	return &GetAllLotteryResponse{
		Lotteries: lotteries,
		Total:     total,
	}, nil
}
