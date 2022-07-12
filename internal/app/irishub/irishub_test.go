package irishub

import (
	"context"
	"testing"

	"github.com/bianjieai/iobscan-explorer-backend/internal/app/config"
)

func TestMain(m *testing.M) {
	cfg := &config.Irishub{
		RpcAddr:  "http://seed-2.mainnet.irisnet.org:26657",
		GrpcAddr: "seed-2.mainnet.irisnet.org:9090",
		ChainId:  "irishub-1",
	}

	InitClient(cfg)
	m.Run()
}

func TestStatus(t *testing.T) {
	status, err := GetCli().Status(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(status)
}
