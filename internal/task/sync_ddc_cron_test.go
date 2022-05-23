package task

import (
	"testing"

	"github.com/bianjieai/iobscan-explorer-backend/internal/common/configs"
	"github.com/bianjieai/iobscan-explorer-backend/internal/repository"
	ddc_sdk "github.com/bianjieai/iobscan-explorer-backend/pkg/libs/ddc-sdk"
)

func TestMain(m *testing.M) {
	confSrv := &configs.Config{
		DataBaseConf: configs.DataBaseConf{
			Database: "bifrost-sync",
			Addrs:    "localhost:27018",
			User:     "iris",
			Passwd:   "irispassword",
		},
		DdcClient: configs.DdcClientConf{
			GatewayURL:       "http://192.168.150.42:8545",
			AuthorityAddress: "0xBcE9AA1924D7197C9C945e43638Bf589f91bcB71",
			ChargeAddress:    "0xF41b6185bFB22E2EFC5fB8395Fa3B952951E2d0b",
			DDC721Address:    "0x74b6114d011891Ac21FD1d586bc7F3407c63c216",
			DDC1155Address:   "0x9f7388e114DfDFAbAF8e4b881894E4C7e1b52C17",
			LogFilepath:      "./log.log",
		},
		Server: configs.Server{
			Prometheus:        "",
			IncreHeight:       1000,
			InsertBatchLimit:  100,
			MaxOperateTxCount: 100,
		},
	}
	ddc_sdk.InitDDCSDKClient(confSrv.DdcClient)
	repository.Start(confSrv)
	//repository.EnsureIndexes()
	defer repository.Stop()
	m.Run()
}

func TestStart(t *testing.T) {
	Start()
}

func TestSyncDdcTask_getMaxHeight(t *testing.T) {
	maxHeight, err := new(SyncDdcTask).getMaxHeight()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(maxHeight)
}

func TestSyncDdcTask_getDdcLatestHeight(t *testing.T) {
	latestHeight, err := new(SyncDdcTask).getDdcLatestHeight()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(latestHeight)
}

func TestSyncDdcTask_getTxsWithScope(t *testing.T) {
	txs := new(SyncDdcTask).getDdcTxsWithScope(0, 763920)
	t.Log(len(txs))
	t.Log(txs)
}

func TestSyncDdcTask_Start(t *testing.T) {
	new(SyncDdcTask).Start()
}

func TestSyncDdcTask_handleDdcTx(t *testing.T) {
	d := new(SyncDdcTask)
	d.loadEvmConfig()
	txs := d.getDdcTxsWithScope(697326, 698600)
	data, err := d.handleDdcTx(&txs[0])
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(data)
}
