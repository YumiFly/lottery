// models/lottery.go
package models

import (
	"time"
)

// LotteryType 表示彩票类型基础表
type LotteryType struct {
	TypeID      string    `gorm:"primaryKey;size:50" json:"type_id"`
	TypeName    string    `gorm:"size:255;not null" json:"type_name"`
	Description string    `gorm:"size:1000" json:"description"`
	CreatedAt   time.Time `gorm:"type:timestamptz;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:timestamptz;default:now()" json:"updated_at"`
}

// Lottery 表示彩票信息表
// Lottery 彩票表模型
type Lottery struct {
	LotteryID              string    `gorm:"primaryKey;size:50" json:"lottery_id"`
	TypeID                 string    `gorm:"size:50;not null" json:"type_id"`
	TicketName             string    `gorm:"size:255;not null" json:"ticket_name"`
	TicketPrice            float64   `gorm:"type:numeric;not null" json:"ticket_price"`
	TicketSupply           int64     `gorm:"type:numeric;not null" json:"ticket_supply"`
	BettingRules           string    `gorm:"size:1000;not null" json:"betting_rules"`
	PrizeStructure         string    `gorm:"size:1000;not null" json:"prize_structure"`
	RegisteredAddr         string    `gorm:"size:255;not null" json:"registered_addr"`
	RolloutContractAddress string    `gorm:"size:255;not null" json:"rollout_contract_address"`
	ContractAddress        string    `gorm:"size:255;not null" json:"contract_address"`
	CreatedAt              time.Time `gorm:"type:timestamptz;default:now()" json:"created_at"`
	UpdatedAt              time.Time `gorm:"type:timestamptz;default:now()" json:"updated_at"`
}

// LotteryIssue 彩票期号表模型
type LotteryIssue struct {
	IssueID        string    `gorm:"primaryKey;size:50" json:"issue_id"`
	LotteryID      string    `gorm:"size:50;not null" json:"lottery_id"`
	IssueNumber    string    `gorm:"size:50;not null" json:"issue_number"`
	SaleEndTime    time.Time `gorm:"type:timestamptz;not null" json:"sale_end_time"`
	DrawTime       time.Time `gorm:"type:timestamptz;not null" json:"draw_time"`
	Status         string    `gorm:"size:100;not null" json:"status"`
	PrizePool      float64   `gorm:"type:numeric;not null" json:"prize_pool"`
	WinningNumbers string    `gorm:"size:100" json:"winning_numbers"`
	RandomSeed     string    `gorm:"size:100" json:"random_seed"`
	DrawTxHash     string    `gorm:"size:66" json:"draw_tx_hash"`
	CreatedAt      time.Time `gorm:"type:timestamptz;default:now()" json:"created_at"`
	UpdatedAt      time.Time `gorm:"type:timestamptz;default:now()" json:"updated_at"`
}

// LotteryTicket 彩票票据表模型
type LotteryTicket struct {
	TicketID        string    `gorm:"primaryKey;size:50" json:"ticket_id"`
	IssueID         string    `gorm:"size:50;not null" json:"issue_id"`
	BuyerAddress    string    `gorm:"size:66;not null" json:"buyer_address"`
	PurchaseTime    time.Time `gorm:"type:timestamptz;not null" json:"purchase_time"`
	BetContent      string    `gorm:"size:100;not null" json:"bet_content"`
	PurchaseAmount  float64   `gorm:"type:numeric;not null" json:"purchase_amount"`
	TransactionHash string    `gorm:"size:66" json:"transaction_hash"`
	CreatedAt       time.Time `gorm:"type:timestamptz;default:now()" json:"created_at"`
	UpdatedAt       time.Time `gorm:"type:timestamptz;default:now()" json:"updated_at"`
}

// Winner 中奖者表模型
type Winner struct {
	WinnerID    string    `gorm:"primaryKey;size:50" json:"winner_id"`
	IssueID     string    `gorm:"size:50;not null" json:"issue_id"`
	TicketID    string    `gorm:"size:50;not null" json:"ticket_id"`
	Address     string    `gorm:"size:66;not null" json:"address"`
	PrizeLevel  string    `gorm:"size:50;not null" json:"prize_level"`
	PrizeAmount float64   `gorm:"type:numeric;not null" json:"prize_amount"`
	ClaimTxHash string    `gorm:"size:66" json:"claim_tx_hash"`
	CreatedAt   time.Time `gorm:"type:timestamptz;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:timestamptz;default:now()" json:"updated_at"`
}
