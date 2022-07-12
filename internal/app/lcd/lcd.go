package lcd

import (
	"encoding/json"
	"fmt"

	"github.com/bianjieai/iobscan-explorer-backend/internal/app/config"
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/utils"
)

var (
	backend        = "http://34.77.68.145:1317/"
	bondTokensUrl  = "/cosmos/staking/v1beta1/pool"
	totalSupplyUrl = "/cosmos/bank/v1beta1/supply"
)

func Init(cfg *config.Lcd) {
	backend = cfg.Backend
	bondTokensUrl = cfg.BondTokensUrl
	totalSupplyUrl = cfg.TotalSupplyUrl
}

func GetTotalSupply() (*Supply, error) {
	var totalSupply Supply
	bytes, err := utils.HttpGet(fmt.Sprintf(fmt.Sprintf("%s%s", backend, totalSupplyUrl)))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &totalSupply)
	return &totalSupply, err
}

func GetBondedTokens() (*BondTokens, error) {
	var bondTokens BondTokens
	bytes, err := utils.HttpGet(fmt.Sprintf(fmt.Sprintf("%s%s", backend, bondTokensUrl)))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &bondTokens)
	return &bondTokens, err
}
