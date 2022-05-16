package contracts

import (
	"encoding/hex"
	ABI "github.com/ethereum/go-ethereum/accounts/abi"
	"strings"
)

var contractABIsMap = map[string]string{
	ContractDDC1155:   _dDC1155MetaData.ABI,
	ContractDDC721:    _dDC721MetaData.ABI,
	ContractAuthority: _authorityMetaData.ABI,
	ContractCharge:    _chargeMetaData.ABI,
}

func GetDDCSupportMethod() (map[string]MethodData, error) {
	methodMap := make(map[string]MethodData)
	for key, ddcABI := range contractABIsMap {
		abiServe, err := ABI.JSON(strings.NewReader(ddcABI))
		if err != nil {
			return nil, err
		}
		for _, method := range abiServe.Methods {
			methodMap[hex.EncodeToString(method.ID)] = MethodData{
				Method:    method,
				Contracts: key,
			}
		}
	}

	return methodMap, nil
}
