package ddc_sdk

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type SignListener struct {
}

// SignEvent
// @Description: 自定义的签名方法
// @receiver s
// @param sender 调用者的账户地址
// @param tx 待签名的交易
// @return *types.Transaction 签名好的交易
// @return error
func (s *SignListener) SignEvent(sender common.Address, tx *types.Transaction) (*types.Transaction, error) {
	//no implement fot no use sign tx
	return tx, nil
}
