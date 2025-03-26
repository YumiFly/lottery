// blockchain/kyc_contract.go
package kyc

import (
	"context"
	"log"
	"math/big"

	"backend/blockchain"
	"backend/config"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// KYCContract 封装 KYC 合约实例
type KYCContract struct {
	Instance *KYC // 由 abigen 生成的合约绑定
}

// 全局 KYC 合约实例
var KYCInstance *KYCContract

// InitKYC 加载 KYC 合约
func InitKYC() {
	if blockchain.Client == nil {
		log.Fatalf("Blockchain client not initialized. Call InitClient first.")
	}

	// 加载 KYC 合约
	contractAddress := common.HexToAddress(config.AppConfig.KYCContractAddress)
	instance, err := NewKYC(contractAddress, blockchain.Client)
	if err != nil {
		log.Fatalf("Failed to load KYC contract: %v", err)
	}

	KYCInstance = &KYCContract{
		Instance: instance,
	}
}

// RegisterKYC 调用链上 register 方法
func RegisterKYC(customerAddress common.Address) error {
	if KYCInstance == nil {
		log.Fatalf("KYC contract not initialized. Call InitKYC first.")
	}

	// 使用管理员账户调用 register（实际中应由用户签名，这里简化处理）
	// 注意：这里假设 KYC.sol 已修改为允许管理员调用 register
	tx, err := KYCInstance.Instance.Register(blockchain.Auth)
	if err != nil {
		return err
	}

	// 等待交易确认
	_, err = bind.WaitMined(context.Background(), blockchain.Client, tx)
	if err != nil {
		return err
	}

	return nil
}

// VerifyKYC 调用链上 verifyKYC 方法
func VerifyKYC(customerAddress common.Address) error {
	if KYCInstance == nil {
		log.Fatalf("KYC contract not initialized. Call InitKYC first.")
	}

	tx, err := KYCInstance.Instance.VerifyKYC(blockchain.Auth, customerAddress)
	if err != nil {
		return err
	}

	// 等待交易确认
	_, err = bind.WaitMined(context.Background(), blockchain.Client, tx)
	if err != nil {
		return err
	}

	return nil
}

// GetKYCStatus 查询链上 KYC 状态
func GetKYCStatus(customerAddress common.Address) (bool, *big.Int, common.Address, error) {
	if KYCInstance == nil {
		log.Fatalf("KYC contract not initialized. Call InitKYC first.")
	}

	isVerified, verificationTime, verifier, err := KYCInstance.Instance.GetKYCStatus(&bind.CallOpts{}, customerAddress)
	if err != nil {
		return false, nil, common.Address{}, err
	}
	return isVerified, verificationTime, verifier, nil
}
