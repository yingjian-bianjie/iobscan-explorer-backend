package task

import (
	"github.com/bianjieai/iobscan-explorer-backend/internal/common/contracts"
	"github.com/bianjieai/iobscan-explorer-backend/pkg/logger"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"strings"
	"testing"

	"github.com/bianjieai/iobscan-explorer-backend/internal/common/configs"
	"github.com/bianjieai/iobscan-explorer-backend/internal/repository"
	ddc_sdk "github.com/bianjieai/iobscan-explorer-backend/pkg/libs/ddc-sdk"
)

func TestMain(m *testing.M) {
	confSrv := &configs.Config{
		DataBaseConf: configs.DataBaseConf{
			//Database: "bifrost-sync",
			//Addrs:    "localhost:27018",
			//User:     "iris",
			//Passwd:   "irispassword",
			Database: "iobscan_wcchain",
			Addrs:    "192.168.150.60:27017",
			User:     "wcchain",
			Passwd:   "wcchainPassword",
		},
		DdcClient: configs.DdcClientConf{
			//GatewayURL:       "http://192.168.150.42:8545",
			//AuthorityAddress: "0xBcE9AA1924D7197C9C945e43638Bf589f91bcB71",
			//ChargeAddress:    "0xF41b6185bFB22E2EFC5fB8395Fa3B952951E2d0b",
			//DDC721Address:    "0x74b6114d011891Ac21FD1d586bc7F3407c63c216",
			//DDC1155Address:   "0x9f7388e114DfDFAbAF8e4b881894E4C7e1b52C17",
			GatewayURL:       "http://192.168.150.60:8545",
			AuthorityAddress: "0xA1750ca7016E49bC3e3A3669cAd33F9a582D58ef",
			ChargeAddress:    "0x661Cf1B1E2f52b8AAe468A858E8aA0C04447D9c0",
			DDC721Address:    "0x42e82933Dfa426995ff53835aD96C7b4aCcdB4C4",
			DDC1155Address:   "0x8f62501A49865dB905414DFC9e8D8E02d7B3bb26",
			LogFilepath:      "./log.log",
		},
		Server: configs.Server{
			Prometheus:        "9092",
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
	evmContractCfgData, err := d.evmCfgModel.FindAll()
	if err != nil {
		logger.Fatal("failed to get data from " + d.evmCfgModel.Name() + err.Error())
	}
	if len(evmContractCfgData) == 0 {
		logger.Fatal(d.evmCfgModel.Name() + " data should config.")
	}
	d.contractABIsMap = make(map[string]abi.ABI, len(evmContractCfgData))
	d.contractTypeNamesMap = make(map[string]string, len(evmContractCfgData))
	for _, val := range evmContractCfgData {
		abiServer, err := abi.JSON(strings.NewReader(val.AbiContent))
		if err != nil {
			logger.Fatal(err.Error())
		}
		d.contractABIsMap[val.Address] = abiServer
		d.contractTypeNamesMap[val.Address] = contracts.DdcTypeName[val.Type]
	}
	txs := d.getDdcTxsWithScope(599035, 600187)
	data, err := d.handleDdcTx(&txs[0])
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(data)
}
