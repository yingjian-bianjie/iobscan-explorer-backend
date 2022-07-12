package irishub

import (
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/config"
	sdk "github.com/irisnet/irishub-sdk-go"
	"github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/types/store"
)

var irishubCli *sdk.IRISHUBClient

func InitClient(c *config.Irishub) {
	options := []types.Option{
		types.KeyDAOOption(store.NewMemory(nil)),
		types.TimeoutOption(10),
	}
	cfg, err := types.NewClientConfig(c.RpcAddr, c.GrpcAddr, c.ChainId, options...)
	if err != nil {
		panic(err)
	}
	cli := sdk.NewIRISHUBClient(cfg)
	irishubCli = &cli
}

func GetCli() *sdk.IRISHUBClient {
	return irishubCli
}
