package task

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
	"strings"

	"github.com/bianjieai/iobscan-explorer-backend/internal/common/constant"
	"github.com/bianjieai/iobscan-explorer-backend/internal/common/contracts"
	"github.com/bianjieai/iobscan-explorer-backend/internal/common/util"
	"github.com/bianjieai/iobscan-explorer-backend/internal/repository"
	ddc_sdk "github.com/bianjieai/iobscan-explorer-backend/pkg/libs/ddc-sdk"
	"github.com/bianjieai/iobscan-explorer-backend/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/mgo.v2"
)

type SyncDdcTask struct {
	ddcTxInfoModel       repository.ExSyncTxEvm
	syncDdcModel         repository.ExSyncDdc
	syncTask             repository.SyncTask
	txModel              repository.Tx
	evmCfgModel          repository.ExEvmContractsConfig
	contractABIsMap      map[string]abi.ABI
	contractTypeNamesMap map[string]string
}

func (d *SyncDdcTask) Name() string {
	return constant.SyncDdcTaskName
}

func (d *SyncDdcTask) Cron() int {
	return constant.CronTimeSyncDdcTask
}

func (d *SyncDdcTask) DoTask(fn func(string) chan bool) error {
	return nil
}
func (d *SyncDdcTask) loadEvmConfig() {
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
	logger.Info("load contract config  data success.")
}
func (d *SyncDdcTask) Start() {
	d.loadEvmConfig()

	handleDdc := func() {
		//logger.Debug("sync ddc cron task start...")
		defer func() {
			if err := recover(); err != nil {
				logger.Error("occur error", logger.Any("err", err))
			}
			//logger.Debug("sync ddc cron task exit...")
		}()

		follow, err := d.syncTask.QueryValidFollowTasks()
		if err != nil {
			logger.Error("failed to get ValidFollowTasks " + err.Error())
			return
		}
		if !follow {
			logger.Warn("waiting sync working is follow......")
			return
		}
		ddcLatestHeight, err := d.getDdcLatestHeight()
		if err != nil {
			logger.Error("failed to get DdcLatestHeight " + err.Error())
			return
		}
		maxHeight, err := d.getMaxHeight()
		if err != nil {
			logger.Error("failed to get MaxHeight " + err.Error())
			return
		}

		if maxHeight > 0 && ddcLatestHeight < maxHeight {
			txs := d.getDdcTxsWithScope(ddcLatestHeight, maxHeight)
			if err := d.handleTxs(txs); err != nil {
				if err != mgo.ErrNotFound && err != constant.ErrDbExist {
					logger.Error("failed to handle Txs " + err.Error())
				}
				return
			}
		}
	}

	util.RunTimer(d.Cron(), util.Sec, func() {
		handleDdc()
	})

}
func (d *SyncDdcTask) handleTxs(txs []repository.Tx) error {
	evmTxs := make([]*repository.ExSyncTxEvm, 0, len(txs))
	var latestTxHeight int64
	for _, tx := range txs {
		if latestTxHeight < tx.Height {
			latestTxHeight = tx.Height
		}
		evmTx, err := d.handleDdcTx(&tx)
		if err != nil {
			return err
		}
		evmTxs = append(evmTxs, evmTx)
	}

	if len(evmTxs) > 0 {
		if err := d.bitchSaveWithTxn(evmTxs); err != nil {
			return err
		}
	}

	//update latest ddc latest_tx_height
	ddc, err := d.syncDdcModel.DdcLatest()
	if err != nil {
		return err
	}
	if ddc.LatestTxHeight < latestTxHeight {
		return d.syncDdcModel.UpdateDdcLatestTxHeight(ddc.ContractAddress, ddc.DdcId, latestTxHeight)
	}

	return nil
}
func (d *SyncDdcTask) parseContractsInput(inputDataStr string, doctx *contracts.DocMsgEthereumTx) error {
	ddcMethodId := inputDataStr[:8]
	abiServe, ok := d.contractABIsMap[doctx.ContractAddr]
	if !ok {
		return constant.SkipErrmsgNoSupportContract
	}
	methodMap, err := contracts.GetDDCSupportMethod(abiServe)
	if err != nil {
		return err
	}
	if val, ok := methodMap[ddcMethodId]; ok {
		doctx.Method = val.Name
		inputData, err := hex.DecodeString(inputDataStr[8:])
		if err != nil {
			return err
		}

		inputs, err := val.Inputs.Unpack(inputData)
		if err != nil {
			return err
		}

		for _, val := range inputs {
			doctx.Inputs = append(doctx.Inputs, util.MarshalJsonIgnoreErr(val))
		}

		return nil
	}

	return constant.SkipErrmsgABIMethodNoFound
}

func (d *SyncDdcTask) handleOneMsg(msg repository.TxMsg, tx *repository.Tx) ([]repository.DdcInfo, []repository.EvmData, error) {

	var ddcsInfo []repository.DdcInfo
	var evmDatas []repository.EvmData

	var msgEtheumTx contracts.DocMsgEthereumTx
	var txData contracts.LegacyTx
	bytesData, err := json.Marshal(msg.Msg)
	if err != nil {
		return ddcsInfo, evmDatas, err
	}
	if err := json.Unmarshal(bytesData, &msgEtheumTx); err != nil {
		return ddcsInfo, evmDatas, err
	}

	if err := json.Unmarshal([]byte(msgEtheumTx.Data), &txData); err != nil {
		return ddcsInfo, evmDatas, err
	}
	msgEtheumTx.ContractAddr = txData.To
	msgEtheumTx.DdcType = d.contractTypeNamesMap[msgEtheumTx.ContractAddr]
	msgEtheumTx.EvmType = contracts.EvmDdcType

	inputDataStr := hex.EncodeToString(common.CopyBytes(txData.Data))
	if err := d.parseContractsInput(inputDataStr, &msgEtheumTx); err != nil {
		if err != constant.SkipErrmsgNoSupportContract && err != constant.SkipErrmsgABIMethodNoFound {
			return ddcsInfo, evmDatas, err
		}
		// skip msg when no support
		logger.Warn(err.Error()+fmt.Sprint(" height: ", tx.Height, " txHash: ", tx.TxHash),
			logger.String("contract_address", txData.To))
	}

	//save txHeight,txTime,Signer to msgEtheumTx
	msgEtheumTx.TxHeight = tx.Height
	msgEtheumTx.TxTime = tx.Time
	if len(tx.Signers) > 0 {
		msgEtheumTx.Signer = tx.Signers[0]
	}

	evmData := repository.EvmData{
		EvmTxHash:       msgEtheumTx.Hash,
		EvmMethod:       msgEtheumTx.Method,
		EvmInputs:       msgEtheumTx.Inputs,
		EvmOutputs:      msgEtheumTx.Outputs,
		DataType:        msgEtheumTx.EvmType,
		ContractAddress: msgEtheumTx.ContractAddr,
	}
	evmDatas = append(evmDatas, evmData)

	ddcIds := contracts.GetDdcIdsByHash(msgEtheumTx)
	ddcMap := make(map[int64]repository.ExSyncDdc, len(ddcIds))
	if len(ddcIds) > 0 {
		ddcs, err := d.syncDdcModel.FindDdcsByDdcIds(msgEtheumTx.ContractAddr, ddcIds)
		if err != nil {
			logger.Warn(err.Error(), logger.String("funcName", "FindDdcsByDdcIds"))
		}
		for _, val := range ddcs {
			ddcMap[val.DdcId] = val
		}
	}

	//opts handle
	ddcOpt, _ := contracts.DdcMethod[msgEtheumTx.Method]
	for i, ddcId := range ddcIds {
		//save evm tx ddc info
		ddcInfo := repository.DdcInfo{
			DdcId: int64(ddcId),
			//DdcName:   ddcName,
			DdcType: msgEtheumTx.DdcType,
			//DdcUri:    ddcUri,
			EvmTxHash: msgEtheumTx.Hash,
			Sender:    msgEtheumTx.Signer,
		}
		if tx.Status == repository.TxStatusSuccess {
			switch ddcOpt {
			case contracts.BurnDdc:
				if dbDdc, ok := ddcMap[int64(ddcId)]; ok {
					ddcInfo.DdcName = dbDdc.DdcName
					ddcInfo.DdcUri = dbDdc.DdcUri
					ddcInfo.DdcSymbol = dbDdc.DdcSymbol
				}
				if err := d.deleteDdcs(int64(ddcId), msgEtheumTx.ContractAddr); err != nil {
					return ddcsInfo, evmDatas, err
				}
				break
			case contracts.TransferDdc:
				if dbDdc, ok := ddcMap[int64(ddcId)]; ok {
					ddcInfo.DdcName = dbDdc.DdcName
					ddcInfo.DdcUri = dbDdc.DdcUri
					ddcInfo.DdcSymbol = dbDdc.DdcSymbol
				}
				if len(msgEtheumTx.Inputs) > 2 {
					// 0:from 1:to 2:ddcId or ddcIds 3:amount
					msgEtheumTx.Inputs[1] = strings.ReplaceAll(msgEtheumTx.Inputs[1], "\"", "")
					if msgEtheumTx.Inputs[1] != "" {
						toAddr, _ := ddc_sdk.Client().HexToBech32(msgEtheumTx.Inputs[1])
						ddcInfo.Recipient = toAddr
					}
				}
				if ddcInfo.Recipient != "" {
					if err := d.syncDdcModel.UpdateOwnerOrUri(msgEtheumTx.ContractAddr, int64(ddcId), ddcInfo.Recipient, ""); err != nil {
						return ddcsInfo, evmDatas, err
					}
				}
				break
			case contracts.MintDdc, contracts.SafeMintDdc, contracts.MintBatchDdc:
				if len(msgEtheumTx.Inputs) > 1 {
					//ddc721 mint,safeMint 0:to 1:ddcURI_
					//ddc1155 safeMint 0:to 1:amount 2:_ddcURI
					//ddc1155 safeMintBatch 0:to 1:amounts 2:_ddcURIs
					msgEtheumTx.Inputs[0] = strings.ReplaceAll(msgEtheumTx.Inputs[0], "\"", "")
					if msgEtheumTx.Inputs[0] != "" {
						toAddr, _ := ddc_sdk.Client().HexToBech32(msgEtheumTx.Inputs[0])
						ddcInfo.Recipient = toAddr
					}

					switch msgEtheumTx.DdcType {
					case contracts.ContractDDC721:
						if len(msgEtheumTx.Inputs) >= 2 {
							msgEtheumTx.Inputs[1] = strings.ReplaceAll(msgEtheumTx.Inputs[1], "\"", "")
							ddcInfo.DdcUri = msgEtheumTx.Inputs[1]
						}
					case contracts.ContractDDC1155:
						if len(msgEtheumTx.Inputs) >= 3 {
							msgEtheumTx.Inputs[2] = strings.ReplaceAll(msgEtheumTx.Inputs[2], "\"", "")
							if ddcOpt == contracts.MintBatchDdc {
								_ddcURIs := contracts.ParseArrStr(msgEtheumTx.Inputs[2])
								if len(_ddcURIs) > i {
									ddcInfo.DdcUri = _ddcURIs[i]
								}
							} else {
								ddcInfo.DdcUri = msgEtheumTx.Inputs[2]
							}
						}

					}
				}

				ddcid := int64(ddcId)
				ddcDoc := d.createExSyncDdcByDdcId(ddcid, &msgEtheumTx)
				ddcOwner, _ := contracts.GetDdcOwner(ddcid, &msgEtheumTx)
				ddcDoc.Owner, _ = ddc_sdk.Client().HexToBech32(ddcOwner)
				ddcDoc.DdcSymbol, _ = contracts.GetDdcSymbol(&msgEtheumTx)
				ddcDoc.DdcData = util.MarshalJsonIgnoreErr(evmData.EvmInputs)
				//handle burn ddc tx owner
				if ddcOwner == "" {
					ddcDoc.Owner = ddcInfo.Recipient
				}
				if ddcInfo.DdcUri == "" {
					ddcDoc.DdcUri = ddcInfo.DdcUri
				}
				ddcInfo.DdcName = ddcDoc.DdcName
				ddcInfo.DdcSymbol = ddcDoc.DdcSymbol

				if tx.Status == repository.TxStatusSuccess {
					if err := d.syncDdcModel.Save(*ddcDoc); err != nil && err != constant.ErrDbExist {
						//return err when failed to save ddc data to ex_sync_ddc
						return ddcsInfo, evmDatas, err
					}
				}

				break
			case contracts.EditDdc:
				if dbDdc, ok := ddcMap[int64(ddcId)]; ok {
					ddcInfo.DdcUri = dbDdc.DdcUri
				}
				if len(msgEtheumTx.Inputs) > 1 {
					// 0:name 1:symbol
					msgEtheumTx.Inputs[0] = strings.ReplaceAll(msgEtheumTx.Inputs[0], "\"", "")
					msgEtheumTx.Inputs[1] = strings.ReplaceAll(msgEtheumTx.Inputs[1], "\"", "")
					ddcInfo.DdcName = msgEtheumTx.Inputs[0]
					ddcInfo.DdcSymbol = msgEtheumTx.Inputs[1]
				}
				if ddcInfo.DdcName != "" || ddcInfo.DdcSymbol != "" {
					if err := d.syncDdcModel.UpdateNameAndSymbol(msgEtheumTx.ContractAddr, int64(ddcId), ddcInfo.DdcName, ddcInfo.DdcSymbol); err != nil {
						return ddcsInfo, evmDatas, err
					}
				}
				break
			case contracts.SetURI:
				if dbDdc, ok := ddcMap[int64(ddcId)]; ok {
					ddcInfo.DdcName = dbDdc.DdcName
					ddcInfo.DdcSymbol = dbDdc.DdcSymbol
				}
				if len(msgEtheumTx.Inputs) > 1 {
					switch msgEtheumTx.DdcType {
					case contracts.ContractDDC721:
						// 0:ddcId 1:ddcUri
						if len(msgEtheumTx.Inputs) >= 2 {
							msgEtheumTx.Inputs[1] = strings.ReplaceAll(msgEtheumTx.Inputs[1], "\"", "")
							ddcInfo.DdcUri = msgEtheumTx.Inputs[1]
						}
					case contracts.ContractDDC1155:
						// 0:owner 1:ddcId 2:ddcUri
						if len(msgEtheumTx.Inputs) >= 3 {
							msgEtheumTx.Inputs[2] = strings.ReplaceAll(msgEtheumTx.Inputs[2], "\"", "")
							ddcInfo.DdcUri = msgEtheumTx.Inputs[2]
						}

					}
				}
				if ddcInfo.DdcUri != "" {
					if err := d.syncDdcModel.UpdateOwnerOrUri(msgEtheumTx.ContractAddr, int64(ddcId), "", ddcInfo.DdcUri); err != nil {
						return ddcsInfo, evmDatas, err
					}
				}
				break
				//case contracts.FreezeDdc, contracts.UnFreezeDdc:
				//	//err := d.syncDdcModel.Update(contractAddr, ddcId,
				//	//	bson.M{"is_freeze": ddcOpt == contracts.FreezeDdc})
				//	//if err != nil {
				//	//	return nil, err
				//	//}
				//	break
			}
		}
		ddcsInfo = append(ddcsInfo, ddcInfo)
	}
	return ddcsInfo, evmDatas, nil
}

func (d *SyncDdcTask) handleDdcTx(tx *repository.Tx) (*repository.ExSyncTxEvm, error) {
	var ddcsInfo []repository.DdcInfo
	var evmDatas []repository.EvmData
	for _, msg := range tx.DocTxMsgs {
		if msg.Type != repository.EthereumTxType {
			continue
		}
		ddcinfos, evmdatas, err := d.handleOneMsg(msg, tx)
		if err != nil {
			return nil, err
		}
		ddcsInfo = append(ddcsInfo, ddcinfos...)
		evmDatas = append(evmDatas, evmdatas...)

	}

	txInfo := &repository.ExSyncTxEvm{
		Time:       tx.Time,
		Height:     tx.Height,
		TxHash:     tx.TxHash,
		Status:     tx.Status,
		Fee:        tx.Fee,
		Types:      tx.Types,
		Signers:    tx.Signers,
		EvmDatas:   evmDatas,
		ExDdcInfos: ddcsInfo,
	}
	return txInfo, nil
}

func (d *SyncDdcTask) createExSyncDdcByDdcId(ddcId int64, msgEtheumTx *contracts.DocMsgEthereumTx) *repository.ExSyncDdc {
	ddcName, _ := contracts.GetDdcName(msgEtheumTx)
	ddcUri, _ := contracts.GetDdcUri(ddcId, msgEtheumTx)
	data := &repository.ExSyncDdc{
		DdcId:   ddcId,
		DdcType: contracts.DdcType[msgEtheumTx.DdcType],
		DdcName: ddcName,
		//DdcSymbl:        ddcSymbol,
		ContractAddress: msgEtheumTx.ContractAddr,
		//DdcData:         ddcData,
		Creator: msgEtheumTx.Signer,
		//Owner:           ddcOwner,
		//Amount:          amount,
		DdcUri:         ddcUri,
		LatestTxHeight: msgEtheumTx.TxHeight,
		LatestTxTime:   msgEtheumTx.TxTime,
	}
	return data
}

func (d *SyncDdcTask) getDdcTxsWithScope(latestHeight, maxHeight int64) []repository.Tx {
	txs, err := d.txModel.FindDdcTx(latestHeight)
	if err != nil {
		logger.Error(err.Error(), logger.String("funcName", "getTxsWithScope"))
		return []repository.Tx{}
	}
	latestHeight += repository.GetSrvConf().IncreHeight
	if latestHeight < maxHeight && len(txs) < repository.GetSrvConf().MaxOperateTxCount {
		retTxs := d.getDdcTxsWithScope(latestHeight, maxHeight)
		txs = append(txs, retTxs...)
	}
	return txs
}

func (d *SyncDdcTask) deleteDdcs(ddcId int64, contractAddr string) error {
	if err := d.syncDdcModel.DeleteDdc(contractAddr, ddcId); err != nil && err != mgo.ErrNotFound {
		return err
	}
	return nil
}

func (d *SyncDdcTask) bitchSaveWithTxn(evmTxs []*repository.ExSyncTxEvm) error {

	insertOps := make([]txn.Op, 0, repository.GetSrvConf().InsertBatchLimit)
	for _, val := range evmTxs {
		op := txn.Op{
			C:      d.ddcTxInfoModel.Name(),
			Id:     bson.NewObjectId(),
			Insert: val,
		}
		insertOps = append(insertOps, op)
		if len(insertOps) >= repository.GetSrvConf().InsertBatchLimit {
			if err := repository.Txn(insertOps); err != nil {
				return err
			}
			insertOps = make([]txn.Op, 0, repository.GetSrvConf().InsertBatchLimit)
		}
	}
	if len(insertOps) > 0 {
		return repository.Txn(insertOps)
	}
	return nil
}

func (d *SyncDdcTask) getMaxHeight() (int64, error) {
	latestTx, err := d.txModel.FindLatestTx()
	if err != nil && err != mgo.ErrNotFound {
		return 0, err
	}
	return latestTx.Height, nil
}

func (d *SyncDdcTask) getDdcLatestHeight() (int64, error) {
	ddc, err := d.syncDdcModel.DdcLatest()
	if err != nil && err != mgo.ErrNotFound {
		return 0, err
	}
	//if get latest tx height from tx_evm when ex_sync_ddc data is empty
	if ddc.LatestTxHeight == 0 {
		txEvm, err := d.ddcTxInfoModel.TxEvmLatest()
		if err != nil && err != mgo.ErrNotFound {
			return 0, err
		}
		return txEvm.Height, nil
	}

	return ddc.LatestTxHeight, nil
}
