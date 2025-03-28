// models/lottery.go
package models

import (
	"time"
)

// LotteryType 表示彩票类型基础表
type LotteryType struct {
	// TypeID 是彩票类型的唯一标识符
	TypeID string `gorm:"column:type_id;type:varchar(50)" json:"type_id"`
	// TypeName 是彩票类型的名称
	TypeName string `gorm:"column:type_name;type:varchar(255);not null" json:"type_name"`
	// Desc 是彩票类型的描述
	Description string `gorm:"column:description;type:varchar(1000)" json:"description"`
	// CreatedAt 是记录的创建时间，格式为时间戳
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null" json:"created_at"`
	// UpdatedAt 是记录的最后更新时间，格式为时间戳
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;not null" json:"updated_at"`
}

// Lottery 表示彩票信息表
type Lottery struct {
	// LotteryID 是彩票的唯一标识符
	LotteryID string `gorm:"column:lottery_id;type:varchar(50)" json:"lottery_id"`
	// TypeID 是关联的彩票类型 ID，与 lottery_types 表的 type_id 对应
	TypeID string `gorm:"column:type_id;type:varchar(50);not null" json:"type_id"`
	// TicketName 是彩票的名称
	TicketName string `gorm:"column:ticket_name;type:varchar(255);not null" json:"ticket_name"`
	// TicketPrice 是每张彩票的票价，存储为字符串
	TicketPrice string `gorm:"column:ticket_price;type:varchar(50);not null" json:"ticket_price"`
	// BettingRules 是彩票的投注规则
	BettingRules string `gorm:"column:betting_rules;type:varchar(1000);not null" json:"betting_rules"`
	// PrizeStructure 是彩票的奖项设置
	PrizeStructure string `gorm:"column:prize_structure;type:varchar(1000);not null" json:"prize_structure"`
	// ContractAddress 是彩票智能合约的地址
	ContractAddress string `gorm:"column:contract_address;type:varchar(66);not null" json:"contract_address"`
	// CreatedAt 是记录的创建时间，格式为时间戳
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null" json:"created_at"`
	// UpdatedAt 是记录的最后更新时间，格式为时间戳
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;not null" json:"updated_at"`
}

// LotteryIssue 表示彩票期号表
type LotteryIssue struct {
	// IssueID 是彩票期号的唯一标识符
	IssueID string `gorm:"column:issue_id;type:varchar(50)" json:"issue_id"`
	// LotteryID 是关联的彩票 ID，与 lottery 表的 lottery_id 对应
	LotteryID string `gorm:"column:lottery_id;type:varchar(50);not null" json:"lottery_id"`
	// IssueNumber 是彩票的期号
	IssueNumber string `gorm:"column:issue_number;type:varchar(50);not null" json:"issue_number"`
	// SaleEndTime 是彩票销售的截止时间，格式为时间戳
	SaleEndTime time.Time `gorm:"column:sale_end_time;type:timestamp;not null" json:"sale_end_time"`
	// DrawTime 是彩票的开奖时间，格式为时间戳
	DrawTime time.Time `gorm:"column:draw_time;type:timestamp;not null" json:"draw_time"`
	// PrizePool 是当前期号的奖池金额，存储为字符串
	PrizePool string `gorm:"column:prize_pool;type:varchar(50);not null" json:"prize_pool"`
	// DrawStatus 是开奖状态，例如 "Pending", "Drawn", "Cancelled"
	DrawStatus string `gorm:"column:draw_status;type:varchar(20);not null" json:"draw_status"`
	// WinningNumbers 是中奖号码
	WinningNumbers string `gorm:"column:winning_numbers;type:varchar(100)" json:"winning_numbers"`
	// RandomSeed 是用于生成中奖号码的随机数种子
	RandomSeed string `gorm:"column:random_seed;type:varchar(100)" json:"random_seed"`
	// DrawTxHash 是开奖操作的区块链交易哈希
	DrawTxHash string `gorm:"column:draw_tx_hash;type:varchar(66)" json:"draw_tx_hash"`
	// CreatedAt 是记录的创建时间，格式为时间戳
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null" json:"created_at"`
	// UpdatedAt 是记录的最后更新时间，格式为时间戳
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;not null" json:"updated_at"`
}

// LotteryTicket 表示一张彩票票据
type LotteryTicket struct {
	// TicketID 是彩票票据的唯一标识符
	TicketID string `gorm:"column:ticket_id;type:varchar(50)" json:"ticket_id"`
	// IssueID 是关联的期号 ID，与 lottery_issues 表的 issue_id 对应
	IssueID string `gorm:"column:issue_id;type:varchar(50);not null" json:"issue_id"`
	// BuyerAddress 是购买者的区块链地址
	BuyerAddress string `gorm:"column:buyer_address;type:varchar(66);not null" json:"buyer_address"`
	// PurchaseTime 是彩票的购买时间，格式为时间戳
	PurchaseTime time.Time `gorm:"column:purchase_time;type:timestamp;not null" json:"purchase_time"`
	// BetContent 是用户的投注内容
	BetContent string `gorm:"column:bet_content;type:varchar(100);not null" json:"bet_content"`
	// PurchaseAmount 是购买彩票的总金额，存储为字符串
	PurchaseAmount string `gorm:"column:purchase_amount;type:varchar(50);not null" json:"purchase_amount"`
	// TransactionHash 是购买彩票的区块链交易哈希
	TransactionHash string `gorm:"column:transaction_hash;type:varchar(66)" json:"transaction_hash"`
	// ClaimStatus 是兑奖状态，例如 "Unclaimed", "Claimed", "Expired"
	ClaimStatus string `gorm:"column:claim_status;type:varchar(20)" json:"claim_status"`
	// ClaimTime 是兑奖时间，格式为时间戳
	ClaimTime time.Time `gorm:"column:claim_time;type:timestamp" json:"claim_time"`
	// ClaimTxHash 是兑奖操作的区块链交易哈希
	ClaimTxHash string `gorm:"column:claim_tx_hash;type:varchar(66)" json:"claim_tx_hash"`
	// CreatedAt 是记录的创建时间，格式为时间戳
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null" json:"created_at"`
	// UpdatedAt 是记录的最后更新时间，格式为时间戳
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;not null" json:"updated_at"`
}

// Winner 表示中奖者信息
type Winner struct {
	// WinnerID 是中奖记录的唯一标识符
	WinnerID string `gorm:"column:winner_id;type:varchar(50)" json:"winner_id"`
	// IssueID 是关联的期号 ID，与 lottery_issues 表的 issue_id 对应
	IssueID string `gorm:"column:issue_id;type:varchar(50);not null" json:"issue_id"`
	// TicketID 是关联的彩票票据 ID，与 lottery_tickets 表的 ticket_id 对应
	TicketID string `gorm:"column:ticket_id;type:varchar(50);not null" json:"ticket_id"`
	// Address 是中奖者的区块链地址
	Address string `gorm:"column:address;type:varchar(66);not null" json:"address"`
	// PrizeLevel 是奖项等级
	PrizeLevel string `gorm:"column:prize_level;type:varchar(50);not null" json:"prize_level"`
	// PrizeAmount 是奖金金额，存储为字符串
	PrizeAmount string `gorm:"column:prize_amount;type:varchar(50);not null" json:"prize_amount"`
	// CreatedAt 是记录的创建时间，格式为时间戳
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null" json:"created_at"`
	// UpdatedAt 是记录的最后更新时间，格式为时间戳
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;not null" json:"updated_at"`
}
