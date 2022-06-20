package contracts

import (
	"encoding/hex"
	ddc_sdk "github.com/bianjieai/iobscan-explorer-backend/pkg/libs/ddc-sdk"
	"github.com/bianjieai/iobscan-explorer-backend/pkg/logger"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"strings"
)

func GetDDCSupportMethod(abiServe abi.ABI) (map[string]abi.Method, error) {
	methodMap := make(map[string]abi.Method)
	for _, method := range abiServe.Methods {
		methodMap[hex.EncodeToString(method.ID)] = method
	}
	return methodMap, nil
}

func NeedRetryCallGetDdcIds(msgEtheumTx DocMsgEthereumTx) []uint64 {
	var (
		retryMaxTimes = 3
		ddcIds        []uint64
	)
	// retry call if not get ddcIds
	for len(ddcIds) == 0 {
		if retryMaxTimes == 0 {
			logger.Warn("ddc-sdk not get ddcIds when retry call 3 times",
				logger.String("evm_txHash", msgEtheumTx.Hash))
			return nil
		}
		retryMaxTimes--
		ddcIds = GetDdcIdsByHash(msgEtheumTx)
	}
	return ddcIds
}

func GetDdcIdsByHash(msgEtheumTx DocMsgEthereumTx) []uint64 {
	var (
		ddcIds []uint64
	)
	switch msgEtheumTx.DdcType {
	case ContractDDC721:
		ddcIds = ddc_sdk.Client().GetDDC721Service().DDCIdByHash(msgEtheumTx.Hash)
		break
	case ContractDDC1155:
		ddcIds = ddc_sdk.Client().GetDDC1155Service().DDCIdByHash(msgEtheumTx.Hash)
		break
	case ContractAuthority:
		ddcIds = ddc_sdk.Client().GetAuthorityService().DDCIdByHash(msgEtheumTx.Hash)
		break
	case ContractCharge:
		ddcIds = ddc_sdk.Client().GetChargeService().DDCIdByHash(msgEtheumTx.Hash)
		break
	}
	return ddcIds
}

func GetDdcUri(ddcId int64, msgEtheumTx *DocMsgEthereumTx) (ddcUri string, err error) {
	switch msgEtheumTx.DdcType {
	case ContractDDC721:
		if ddcUri, err = ddc_sdk.Client().GetDDC721Service().DdcURI(ddcId); err != nil {
			logger.Warn(err.Error(), logger.String("funcName", "DDC721.GetDdcUri"))
			return "", err
		}
	case ContractDDC1155:
		if ddcUri, err = ddc_sdk.Client().GetDDC1155Service().DdcURI(ddcId); err != nil {
			logger.Warn(err.Error(), logger.String("funcName", "DDC1155.GetDdcUri"))
			return "", err
		}
	}
	return
}

func GetDdcOwner(ddcId int64, msgEtheumTx *DocMsgEthereumTx) (owner string, err error) {
	switch msgEtheumTx.DdcType {
	case ContractDDC721:
		if owner, err = ddc_sdk.Client().GetDDC721Service().OwnerOf(ddcId); err != nil {
			logger.Warn(err.Error(), logger.String("funcName", "DDC721.DdcOwner"))
			return "", err
		}

	}
	return
}

func GetDdcSymbol(msgEtheumTx *DocMsgEthereumTx) (owner string, err error) {
	switch msgEtheumTx.DdcType {
	case ContractDDC721:
		if owner, err = ddc_sdk.Client().GetDDC721Service().Symbol(); err != nil {
			logger.Warn(err.Error(), logger.String("funcName", "DDC721.GetDdcSymbol"))
			return "", err
		}

	}
	return
}

func GetDdcName(msgEtheumTx *DocMsgEthereumTx) (owner string, err error) {
	switch msgEtheumTx.DdcType {
	case ContractDDC721:
		if owner, err = ddc_sdk.Client().GetDDC721Service().Name(); err != nil {
			logger.Warn(err.Error(), logger.String("funcName", "DDC721.GetDdcName"))
			return "", err
		}

	}
	return
}

func GetDdcAmount(owner string, ddcId int64, msgEtheumTx *DocMsgEthereumTx) (amount uint64, err error) {
	switch msgEtheumTx.DdcType {
	case ContractDDC1155:
		if amount, err = ddc_sdk.Client().GetDDC1155Service().BalanceOf(owner, ddcId); err != nil {
			logger.Warn(err.Error(), logger.String("funcName", "DDC1155.GetDdcName"))
			return 0, err
		}

	}
	return
}

func NeedRetryCallGet(ddcId int64, msgEtheumTx *DocMsgEthereumTx, call func(int64, *DocMsgEthereumTx) (string, error)) (string, error) {
	var (
		ret           string
		err           error
		retryMaxTimes = 3
	)
	// retry call if not get uri or owner
	for ret == "" {
		if retryMaxTimes == 0 {
			logger.Warn("ddc-sdk not get uri or owner when retry call 3 times",
				logger.String("evm_txHash", msgEtheumTx.Hash))
			return ret, err
		}
		retryMaxTimes--
		ret, err = call(ddcId, msgEtheumTx)
	}
	return ret, err
}

// GetTxReceipt
// @Description: 运营方或平台方根据交易哈希对交易回执信息进行查询。
// @receiver t
// @param txHash: 交易哈希
// @return string： 交易回执
// @return error
func GetTxReceipt(txHash string) (*types.Receipt, error) {
	return ddc_sdk.Client().GetTxReceipt(txHash)
}

func ParseArrStr(arrStr string) []string {
	start := strings.Index(arrStr, "[")
	end := strings.Index(arrStr, "]")
	if start < end && start >= 0 {
		data := arrStr[start+1 : end]
		urisArr := strings.Split(data, ",")
		return urisArr
	}
	return nil
}

func NeedRetryCallGetTxReceipt(txHash string) (*types.Receipt, error) {
	var (
		receipt       *types.Receipt
		err           error
		retryMaxTimes = 3
	)
	// retry call if not get receipt
	for receipt == nil {
		if retryMaxTimes == 0 {
			logger.Warn("ddc-sdk not get receipt when retry call 3 times",
				logger.String("evm_txHash", txHash))
			return receipt, err
		}
		retryMaxTimes--
		receipt, err = GetTxReceipt(txHash)
	}
	return receipt, err
}
