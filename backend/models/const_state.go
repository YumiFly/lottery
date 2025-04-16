// models/contract_state.go (或您存储 ContractState 相关代码的文件)

package models

// ContractState 合约状态
type ContractState uint8

const (
	ContractStateReady      ContractState = 0
	ContractStateDistribute ContractState = 1
	ContractStateRollout    ContractState = 2 // 添加这一行
	ContractStateTerminal   ContractState = 3
)

const (
	//IssueStatusPending 待开奖
	IssueStatusPending = "PENDING"
	//IssueStatusDrawing 开奖中
	IssueStatusDrawing = "DRAWING"
	//IssueStatusDrawn 已开奖
	IssueStatusDrawn = "DRAWN"
)
