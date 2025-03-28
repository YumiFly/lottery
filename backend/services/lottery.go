// services/lottery.go
package services

import (
	lotteryBlockchain "backend/blockchain/lottery"
	"backend/db"
	"backend/models"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// CreateLotteryType 创建彩票类型
// 该方法用于在数据库中创建一条彩票类型记录，并自动设置创建和更新时间。
func CreateLotteryType(lotteryType *models.LotteryType) error {
	// 自动设置创建和更新时间
	lotteryType.CreatedAt = time.Now()
	lotteryType.UpdatedAt = time.Now()

	// 插入数据库
	return db.DB.Create(lotteryType).Error
}

// CreateLottery 创建彩票
// 该方法用于在数据库中创建一条彩票记录，并自动设置创建和更新时间。
func CreateLottery(lottery *models.Lottery) error {
	// 手动检查 type_id 是否存在
	var lotteryType models.LotteryType
	if err := db.DB.Where("type_id = ?", lottery.TypeID).First(&lotteryType).Error; err != nil {
		return errors.New("lottery type not found")
	}

	// 自动设置创建和更新时间
	lottery.CreatedAt = time.Now()
	lottery.UpdatedAt = time.Now()

	// 插入数据库
	return db.DB.Create(lottery).Error
}

// CreateIssue 创建彩票期号
// 该方法用于在数据库中创建一条彩票期号记录，并自动设置创建和更新时间。
func CreateIssue(issue *models.LotteryIssue) error {
	// 手动检查 lottery_id 是否存在
	var lottery models.Lottery
	if err := db.DB.Where("lottery_id = ?", issue.LotteryID).First(&lottery).Error; err != nil {
		return errors.New("lottery not found")
	}

	// 手动检查 issue_number 是否重复
	var existingIssue models.LotteryIssue
	if err := db.DB.Where("issue_number = ?", issue.IssueNumber).First(&existingIssue).Error; err == nil {
		return errors.New("issue number already exists")
	}

	// 设置默认值
	issue.DrawStatus = "Pending"
	issue.CreatedAt = time.Now()
	issue.UpdatedAt = time.Now()

	return db.DB.Create(issue).Error
}

// PurchaseTicket 购买彩票
// 该方法用于在数据库中创建一条彩票票据记录，并自动设置创建和更新时间。
func PurchaseTicket(ticket *models.LotteryTicket) error {
	// 手动检查 issue_id 是否存在
	var issue models.LotteryIssue
	if err := db.DB.Where("issue_id = ?", ticket.IssueID).First(&issue).Error; err != nil {
		return errors.New("issue not found")
	}

	// 检查销售是否已截止,
	if time.Now().After(issue.SaleEndTime) {
		return errors.New("sale has ended")
	}

	//设置开奖状态
	ticket.ClaimStatus = "Unclaimed"
	// 记录购买时间
	ticket.PurchaseTime = time.Now()
	ticket.CreatedAt = time.Now()
	ticket.UpdatedAt = time.Now()

	return db.DB.Create(ticket).Error
}

// DrawLottery 开奖
// 该方法用于处理彩票的开奖逻辑，更新期号状态并记录中奖号码。
func DrawLottery(issueID string) error {
	// 查询期号
	var issue models.LotteryIssue
	if err := db.DB.Where("issue_id = ?", issueID).First(&issue).Error; err != nil {
		return errors.New("issue not found")
	}

	// 检查是否已开奖
	if issue.DrawStatus != "Pending" {
		return errors.New("lottery already drawn or cancelled")
	}

	// 检查是否到达开奖时间
	if time.Now().Before(issue.DrawTime) {
		return errors.New("draw time not reached")
	}

	// 查询彩票信息以获取合约地址
	var lottery models.Lottery
	if err := db.DB.Where("lottery_id = ?", issue.LotteryID).First(&lottery).Error; err != nil {
		return errors.New("lottery not found")
	}

	// 调用 SimpleRollout 合约的 rolloutCall 方法
	// 将字符串形式的地址转换为 common.Address
	contractAddr := common.HexToAddress(lottery.ContractAddress)
	requestID, err := lotteryBlockchain.RolloutCall(contractAddr)
	if err != nil {
		return fmt.Errorf("failed to call rollout: %v", err)
	}

	// 监听 DiceLanded 事件，获取随机数结果
	err = lotteryBlockchain.ListenForDiceLanded(big.NewInt(int64(requestID)), func(results []*big.Int) {
		// 将随机数结果转换为字符串
		winningNumbers := fmt.Sprintf("%d,%d,%d", results[0], results[1], results[2])

		// 更新期号状态
		issue.DrawStatus = "Drawn"
		issue.WinningNumbers = winningNumbers
		issue.RandomSeed = fmt.Sprintf("RequestID: %d", requestID)

		// TODO：获取交易哈希

		// 记录交易哈希
		issue.DrawTxHash = ""
		issue.UpdatedAt = time.Now()

		// 保存期号
		if err := db.DB.Save(&issue).Error; err != nil {
			fmt.Printf("Failed to update issue: %v", err)
			return
		}
		fmt.Printf("Lottery drawn: %s\n", winningNumbers)

		// 根据issueID查询所有购买的彩票信息
		var tickets []models.LotteryTicket
		if err := db.DB.Where("issue_id = ?", issueID).Find(&tickets).Error; err != nil {
			fmt.Printf("Failed to get tickets: %v", err)
			return
		}
		// 根据中奖号码计算中奖者
		for _, ticket := range tickets {
			// 这里简化处理，实际中需要比较投注内容和中奖号码
			if ticket.BetContent == winningNumbers {
				ticket.ClaimStatus = "Claimed"
			} else {
				ticket.ClaimStatus = "Unclaimed"
			}
			if err := db.DB.Save(&ticket).Error; err != nil {
				fmt.Printf("Failed to update ticket: %v", err)
			}
			//将中奖者地址和中奖号码发送给前端
			fmt.Printf("Ticket %s is %s\n", ticket.TicketID, ticket.ClaimStatus)

		}

		// 这里简化处理，实际中需要比较投注内容和中奖号码
	})
	if err != nil {
		return fmt.Errorf("failed to listen for DiceLanded: %v", err)
	}

	return nil
}

func GetAllLotteryTypes() ([]models.LotteryType, error) {
	var lotteryTypes []models.LotteryType
	if err := db.DB.Find(&lotteryTypes).Error; err != nil {
		return nil, err
	}
	return lotteryTypes, nil
}

func GetAllLotteries() ([]models.Lottery, error) {
	var lotteries []models.Lottery
	if err := db.DB.Find(&lotteries).Error; err != nil {
		return nil, err
	}
	return lotteries, nil
}

func GetLotteryByTypeID(typeID string) (*models.Lottery, error) {
	var lottery models.Lottery
	if err := db.DB.Where("type_id = ?", typeID).First(&lottery).Error; err != nil {
		return nil, err
	}
	return &lottery, nil
}

func GetLatestIssueByLotteryID(lotteryID string) (*models.LotteryIssue, error) {
	var issue models.LotteryIssue
	if err := db.DB.Where("lottery_id = ?", lotteryID).Order("issue_id desc").First(&issue).Error; err != nil {
		return nil, err
	}
	return &issue, nil
}

func GetExpiringIssues() ([]models.LotteryIssue, error) {
	var issues []models.LotteryIssue
	if err := db.DB.Where("draw_status = ?", "Pending").Order("sale_end_time asc").Find(&issues).Error; err != nil {
		return nil, err
	}
	return issues, nil
}

func GetPurchasedTicketsByCustomerAddress(customerAddress string) ([]models.LotteryTicket, error) {
	var tickets []models.LotteryTicket
	if err := db.DB.Where("buyer_address = ?", customerAddress).Find(&tickets).Error; err != nil {
		return nil, err
	}
	return tickets, nil
}

func GetDrawnLotteryByIssueID(issueID string) (*models.LotteryIssue, error) {
	var issue models.LotteryIssue
	if err := db.DB.Where("issue_id = ?", issueID).First(&issue).Error; err != nil {
		return nil, err
	}
	if issue.DrawStatus != "Drawn" {
		return nil, errors.New("issue not drawn")
	}
	return &issue, nil
}
