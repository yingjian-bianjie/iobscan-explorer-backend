package contracts

import (
	"encoding/hex"
	"github.com/bianjieai/iobscan-explorer-backend/internal/common/constant"
	ddc_sdk "github.com/bianjieai/iobscan-explorer-backend/pkg/libs/ddc-sdk"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

func GetDDCSupportMethod(abiServe abi.ABI) (map[string]abi.Method, error) {
	methodMap := make(map[string]abi.Method)
	for _, method := range abiServe.Methods {
		methodMap[hex.EncodeToString(method.ID)] = method
	}
	return methodMap, nil
}

func GetDdcIdsByHash(msgEtheumTx DocMsgEthereumTx) ([]uint64, error) {
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
	default:
		return nil, constant.SkipErrmsgABITypeNoFound
	}
	return ddcIds, nil
}

func GetDdcUri(ddcId int64, msgEtheumTx *DocMsgEthereumTx) (ddcUri string, err error) {
	switch msgEtheumTx.DdcType {
	case ContractDDC721:
		if ddcUri, err = ddc_sdk.Client().GetDDC721Service().DdcURI(ddcId); err != nil {
			return "", err
		}
	case ContractDDC1155:
		if ddcUri, err = ddc_sdk.Client().GetDDC1155Service().DdcURI(ddcId); err != nil {
			return "", err
		}
	default:
		return "", constant.SkipErrmsgNoSupport

	}
	return
}

func GetDdcOwner(ddcId int64, msgEtheumTx *DocMsgEthereumTx) (owner string, err error) {
	switch msgEtheumTx.DdcType {
	case ContractDDC721:
		if owner, err = ddc_sdk.Client().GetDDC721Service().OwnerOf(ddcId); err != nil {
			return "", err
		}
	default:
		return "", constant.SkipErrmsgNoSupport

	}
	return
}

func GetDdcSymbol(msgEtheumTx *DocMsgEthereumTx) (owner string, err error) {
	switch msgEtheumTx.DdcType {
	case ContractDDC721:
		if owner, err = ddc_sdk.Client().GetDDC721Service().Symbol(); err != nil {
			return "", err
		}
	default:
		return "", constant.SkipErrmsgNoSupport

	}
	return
}

func GetDdcName(msgEtheumTx *DocMsgEthereumTx) (owner string, err error) {
	switch msgEtheumTx.DdcType {
	case ContractDDC721:
		if owner, err = ddc_sdk.Client().GetDDC721Service().Name(); err != nil {
			return "", err
		}
	default:
		return "", constant.SkipErrmsgNoSupport

	}
	return
}

func GetDdcAmount(owner string, ddcId int64, msgEtheumTx *DocMsgEthereumTx) (amount uint64, err error) {
	switch msgEtheumTx.DdcType {
	case ContractDDC1155:
		if amount, err = ddc_sdk.Client().GetDDC1155Service().BalanceOf(owner, ddcId); err != nil {
			return 0, err
		}
	default:
		return 0, constant.SkipErrmsgNoSupport

	}
	return
}
