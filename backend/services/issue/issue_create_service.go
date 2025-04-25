package issue

import (
	"context"
	"time"

	"backend/blockchain"
	"backend/models"
	"backend/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// CreateIssueParams 定义创建期号的参数结构
type CreateIssueParams struct {
	LotteryID      string
	IssueNumber    string
	SaleEndTime    time.Time
	DrawTime       time.Time
	Status         string
	PrizePool      float64
	WinningNumbers string
	RandomSeed     string
	DrawTxHash     string
}

// IssueCreateService 封装期号创建的业务逻辑
type IssueCreateService struct {
	db *gorm.DB
}

// NewIssueCreateService 创建 IssueCreateService 实例
func NewIssueCreateService(db *gorm.DB) *IssueCreateService {
	return &IssueCreateService{db: db}
}

// validateCreateIssueParams 验证创建期号的参数
func (s *IssueCreateService) validateCreateIssueParams(params CreateIssueParams) error {
	// 验证 lottery_id 存在
	var lottery models.Lottery
	if err := s.db.WithContext(context.Background()).
		Where("lottery_id = ?", params.LotteryID).
		First(&lottery).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Logger.Warn("Lottery not found", "lottery_id", params.LotteryID)
			return utils.NewBadRequestError("Lottery not found", nil)
		}
		return utils.NewInternalError("Failed to check lottery ID", errors.Wrap(err, "database error"))
	}

	// 验证 issue_number 唯一性
	var existingIssue models.LotteryIssue
	if err := s.db.WithContext(context.Background()).
		Where("lottery_id = ? AND issue_number = ?", params.LotteryID, params.IssueNumber).
		First(&existingIssue).Error; err == nil {
		utils.Logger.Warn("Issue number already exists", "issue_number", params.IssueNumber)
		return utils.NewBadRequestError("Issue number already exists", nil)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.NewInternalError("Failed to check issue number uniqueness", errors.Wrap(err, "database error"))
	}

	// 验证状态
	if !map[string]bool{models.IssueStatusPending: true, models.IssueStatusDrawn: true}[params.Status] {
		return utils.NewBadRequestError("Invalid status value", nil)
	}

	// 验证时间
	if params.SaleEndTime.After(params.DrawTime) {
		return utils.NewBadRequestError("Sale end time cannot be later than draw time", nil)
	}
	if params.SaleEndTime.Before(time.Now()) {
		return utils.NewBadRequestError("Sale end time cannot be earlier than current time", nil)
	}

	// 验证字段长度
	if len(params.WinningNumbers) > 100 || len(params.RandomSeed) > 100 || len(params.DrawTxHash) > 66 {
		return utils.NewBadRequestError("Optional field length exceeded", nil)
	}

	return nil
}

// CreateIssue 创建新的彩票期号
//
// 参数:
//   - ctx: 请求上下文
//   - params: 创建参数，包括 lottery_id, issue_number 等
//
// 返回:
//   - *LotteryIssue: 创建的期号记录
//   - common.Hash: 区块链交易哈希
//   - error: 创建错误或参数无效
func (s *IssueCreateService) CreateIssue(ctx context.Context, params CreateIssueParams) (*models.LotteryIssue, common.Hash, error) {
	// 验证参数
	if err := s.validateCreateIssueParams(params); err != nil {
		return nil, common.Hash{}, err
	}

	// 构造期号记录
	issue := models.LotteryIssue{
		IssueID:        generateUUID(),
		LotteryID:      params.LotteryID,
		IssueNumber:    params.IssueNumber,
		SaleEndTime:    params.SaleEndTime,
		DrawTime:       params.DrawTime,
		Status:         params.Status,
		PrizePool:      0, // 按原始代码，初始奖池为 0
		WinningNumbers: params.WinningNumbers,
		RandomSeed:     params.RandomSeed,
		DrawTxHash:     params.DrawTxHash,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	var lottery models.Lottery

	// 日志记录
	utils.Logger.Info("Creating issue", "issue_id", issue.IssueID, "lottery_id", issue.LotteryID)

	// 执行区块链交易
	executeTx := func() (common.Hash, error) {
		if err := s.db.WithContext(ctx).Preload("LotteryType").
			Where("lottery_id = ?", issue.LotteryID).
			First(&lottery).Error; err != nil {
			utils.Logger.Warn("Lottery not found", "lottery_id", issue.LotteryID)
			return common.Hash{}, utils.NewBadRequestError("Lottery not found", err)
		}

		// 再次确认期号唯一性（防止并发）
		var existingIssue models.LotteryIssue
		if err := s.db.WithContext(ctx).
			Where("lottery_id = ? AND issue_number = ?", issue.LotteryID, issue.IssueNumber).
			First(&existingIssue).Error; err == nil {
			utils.Logger.Warn("Issue number already exists", "issue_number", issue.IssueNumber)
			return common.Hash{}, utils.NewBadRequestError("Issue number already exists", nil)
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return common.Hash{}, utils.NewInternalError("Failed to check issue number uniqueness", errors.Wrap(err, "database error"))
		}

		// 连接彩票合约
		contract, err := blockchain.ConnectLotteryContract(lottery.ContractAddress)
		if err != nil {
			return common.Hash{}, utils.NewInternalError("Failed to connect to lottery contract", errors.Wrap(err, "contract connection error"))
		}

		// 获取当前合约状态
		currentState, err := contract.GetState(nil)
		if err != nil {
			utils.Logger.Error("Failed to get contract state", "error", err)
			return common.Hash{}, utils.NewInternalError("Failed to get contract state", errors.Wrap(err, "contract state error"))
		}
		utils.Logger.Info("Current contract state", "state", currentState)

		// 设置合约状态为 Distribute
		tx, err := contract.TransState(blockchain.Auth, uint8(1))
		if err != nil {
			utils.Logger.Error("Failed to set state to Distribute", "error", err)
			if tx != nil {
				return tx.Hash(), utils.NewInternalError("Failed to set state to Distribute", errors.Wrap(err, "transaction error"))
			}
			return common.Hash{}, utils.NewInternalError("Failed to set state to Distribute", errors.Wrap(err, "transaction error"))
		}

		// 等待交易确认
		receipt, err := bind.WaitMined(ctx, blockchain.Client, tx)
		if err != nil {
			utils.Logger.Error("Transaction failed", "tx_hash", tx.Hash().Hex(), "error", err)
			return tx.Hash(), utils.NewInternalError("Transaction failed", errors.Wrap(err, "transaction mining error"))
		}
		if receipt.Status != 1 {
			utils.Logger.Error("Transaction failed", "tx_hash", tx.Hash().Hex(), "status", receipt.Status)
			return tx.Hash(), utils.NewInternalError("Transaction failed", nil)
		}

		// 保存到数据库
		if err := s.db.WithContext(ctx).Create(&issue).Error; err != nil {
			utils.Logger.Error("Failed to save issue to database", "error", err)
			return tx.Hash(), utils.NewInternalError("failed to save issue to database", err)
		}

		utils.Logger.Info("Issue created successfully", "issue_id", issue.IssueID)
		return tx.Hash(), nil
	}

	// 执行区块链交易
	data := []byte{}
	txhash, err := blockchain.WithBlockchain(ctx, data, executeTx)
	if err != nil {
		return nil, common.Hash{}, err
	}
	issue.Lottery = lottery
	return &issue, txhash, nil
}

// generateUUID 生成唯一 ID（示例）
func generateUUID() string {
	return "issue-" + time.Now().Format("20060102150405")
}
