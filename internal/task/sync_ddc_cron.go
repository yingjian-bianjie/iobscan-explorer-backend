package task

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bianjieai/iobscan-explorer-backend/internal/common/constant"
	"github.com/bianjieai/iobscan-explorer-backend/internal/common/contracts"
	"github.com/bianjieai/iobscan-explorer-backend/internal/common/util"
	"github.com/bianjieai/iobscan-explorer-backend/internal/repository"
	ddc_sdk "github.com/bianjieai/iobscan-explorer-backend/pkg/libs/ddc-sdk"
	"github.com/bianjieai/iobscan-explorer-backend/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
)

type SyncDdcTask struct {
	ddcTxInfoModel repository.ExSyncTxEvm
	syncDdcModel   repository.ExSyncDdc
	syncTask       repository.SyncTask
	txModel        repository.Tx
}

func (d SyncDdcTask) Start() {
	logger.Debug("sync ddc cron task start...")
	defer func() {
		if err := recover(); err != nil {
			logger.Error("occur error", logger.Any("err", err))
		}
		logger.Debug("sync ddc cron task exit...")
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
		if err != mgo.ErrNotFound {
			logger.Error("failed to get DdcLatestHeight " + err.Error())
		}
		return
	}
	maxHeight, err := d.getMaxHeight()
	if err != nil {
		if err != mgo.ErrNotFound {
			logger.Error("failed to get MaxHeight " + err.Error())
		}
		return
	}

	if maxHeight > 0 && ddcLatestHeight < maxHeight {
		txs := d.getDdcTxsWithScope(ddcLatestHeight, maxHeight)
		if err := d.handleTxs(txs); err != nil {
			if err != mgo.ErrNotFound {
				logger.Error("failed to handle Txs " + err.Error())
			}
			return
		}
	}

}
func (d SyncDdcTask) handleTxs(txs []repository.Tx) error {
	var latest_tx_height int64
	for _, tx := range txs {
		if latest_tx_height < tx.Height {
			latest_tx_height = tx.Height
		}
		err := d.handleDdcTx(&tx)
		if err != nil && err != constant.ErrDbExist {
			return err
		}
	}

	//update latest ddc latest_tx_height
	ddc, err := d.syncDdcModel.DdcLatest()
	if err != nil {
		return err
	}
	if ddc.LatestTxHeight < latest_tx_height {
		return d.syncDdcModel.Update(ddc.ContractAddress, ddc.DdcId, bson.M{"latest_tx_height": latest_tx_height})
	}

	return nil
}
func parseContractsInput(inputDataStr string, doctx *contracts.DocMsgEthereumTx) error {
	ddcMethodId := inputDataStr[:8]
	methodMap, err := contracts.GetDDCSupportMethod()
	if err != nil {
		return err
	}
	if val, ok := methodMap[ddcMethodId]; ok {
		doctx.EvmType = contracts.EvmDdcType
		doctx.DdcType = val.Contracts
		doctx.Method = val.Method.Name
		inputData, err := hex.DecodeString(inputDataStr[8:])
		if err != nil {
			return err
		}

		inputs, err := val.Method.Inputs.Unpack(inputData)
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

func (d SyncDdcTask) handleOneMsg(msg repository.TxMsg, tx *repository.Tx) ([]repository.DdcInfo, []repository.EvmData, error) {

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

	inputDataStr := hex.EncodeToString(common.CopyBytes(txData.Data))
	if err := parseContractsInput(inputDataStr, &msgEtheumTx); err != nil {
		return ddcsInfo, evmDatas, err
	}

	//save txHeight,txTime,Signer to msgEtheumTx
	msgEtheumTx.TxHeight = tx.Height
	msgEtheumTx.TxTime = tx.Time
	if len(tx.Signers) > 0 {
		msgEtheumTx.Signer = tx.Signers[0]
	}

	msgEtheumTx.ContractAddr = txData.To
	evmData := repository.EvmData{
		EvmTxHash:       msgEtheumTx.Hash,
		EvmMethod:       msgEtheumTx.Method,
		EvmInputs:       msgEtheumTx.Inputs,
		EvmOutputs:      msgEtheumTx.Outputs,
		DataType:        msgEtheumTx.EvmType,
		ContractAddress: msgEtheumTx.ContractAddr,
	}
	evmDatas = append(evmDatas, evmData)

	//todo ddcIds may be no found for burned ddc/Nft
	ddcIds, err := contracts.GetDdcIdsByHash(msgEtheumTx)
	if err != nil {
		return ddcsInfo, evmDatas, err
	}

	//opts handle
	ddcOpt, exist := contracts.DdcMethod[msgEtheumTx.Method]
	if !exist {
		//todo
		logger.Error("no support ddcMethod: " + msgEtheumTx.Method)
		return ddcsInfo, evmDatas, constant.SkipErrmsgNoSupport
	}
	for _, ddcId := range ddcIds {
		ddcName, _ := contracts.GetDdcName(&msgEtheumTx)
		ddcUri, _ := contracts.GetDdcUri(int64(ddcId), &msgEtheumTx)
		ddcInfo := repository.DdcInfo{
			DdcId:     int64(ddcId),
			DdcName:   ddcName,
			DdcType:   msgEtheumTx.DdcType,
			DdcUri:    ddcUri,
			EvmTxHash: msgEtheumTx.Hash,
		}
		ddcsInfo = append(ddcsInfo, ddcInfo)
		if tx.Status == repository.TxStatusSuccess {
			switch ddcOpt {
			case contracts.BurnDdc:
				if err := d.deleteDdcs(int64(ddcId), msgEtheumTx.ContractAddr); err != nil {
					return ddcsInfo, evmDatas, err
				}
				break
			case contracts.TransferDdc:
				ddcOwner, err := contracts.GetDdcOwner(int64(ddcId), &msgEtheumTx)
				if err != nil {
					//return err when failed to get ddc data
					return ddcsInfo, evmDatas, err
				}
				owner, _ := ddc_sdk.Client().HexToBech32(ddcOwner)
				if err := d.editDdcs(ddcId, msgEtheumTx.ContractAddr, bson.M{"owner": owner}); err != nil {
					return ddcsInfo, evmDatas, err
				}

				break
			case contracts.MintDdc, contracts.SafeMintDdc, contracts.MintBatchDdc:
				ddcid := int64(ddcId)
				ddcDoc := d.createExSyncDdcByDdcId(ddcid, &msgEtheumTx)
				ddcDoc.DdcUri = ddcUri
				ddcOwner, _ := contracts.GetDdcOwner(ddcid, &msgEtheumTx)
				ddcDoc.Owner, _ = ddc_sdk.Client().HexToBech32(ddcOwner)
				ddcDoc.DdcSymbl, _ = contracts.GetDdcSymbol(&msgEtheumTx)
				ddcDoc.DdcName = ddcName
				//if err != nil {
				//	//return err when failed to get ddc data
				//	return ddcsInfo, evmDatas, err
				//}
				ddcDoc.DdcData = util.MarshalJsonIgnoreErr(evmData.EvmInputs)

				if tx.Status == repository.TxStatusSuccess {
					if err := d.syncDdcModel.Save(*ddcDoc); err != nil {
						//return err when failed to save ddc data to ex_sync_ddc
						return ddcsInfo, evmDatas, err
					}
				}

				break
			case contracts.EditDdc:
				name, err := contracts.GetDdcName(&msgEtheumTx)
				if err != nil {
					logger.Warn("failed to get ddcName " + err.Error())
				}
				symbol, err := contracts.GetDdcSymbol(&msgEtheumTx)
				if err != nil {
					logger.Warn("failed to get ddcSymbol " + err.Error())
				}
				if name != "" || symbol != "" {
					editData := bson.M{}
					if name != "" {
						editData["ddc_name"] = name
					}
					if symbol != "" {
						editData["ddc_symbol"] = symbol
					}
					if err := d.editDdcs(ddcId, msgEtheumTx.ContractAddr, editData); err != nil {
						return ddcsInfo, evmDatas, err
					}
				}
				break
			case contracts.SetURI:
				if err := d.editDdcs(ddcId, msgEtheumTx.ContractAddr, bson.M{"ddc_uri": ddcUri}); err != nil {
					return ddcsInfo, evmDatas, err
				}
				break
			//case contracts.FreezeDdc, contracts.UnFreezeDdc:
			//	//err := d.syncDdcModel.Update(contractAddr, ddcId,
			//	//	bson.M{"is_freeze": ddcOpt == contracts.FreezeDdc})
			//	//if err != nil {
			//	//	return nil, err
			//	//}
			//	break
			default:
				//todo
				return ddcsInfo, evmDatas, constant.SkipErrmsgNoSupport
			}
		}
	}
	return ddcsInfo, evmDatas, nil
}

func (d SyncDdcTask) handleDdcTx(tx *repository.Tx) error {
	if len(tx.DocTxMsgs) == 0 {
		return errors.New("empty msg")
	}
	var ddcsInfo []repository.DdcInfo
	var evmDatas []repository.EvmData
	for _, msg := range tx.DocTxMsgs {
		if msg.Type != repository.EthereumTxType {
			continue
		}
		ddcinfos, evmdatas, err := d.handleOneMsg(msg, tx)
		if err != nil {
			if err == constant.SkipErrmsgABIMethodNoFound {
				logger.Warn("skip tx msg for " + err.Error() + fmt.Sprint(" height: ", tx.Height, " txHash: ", tx.TxHash))
				continue
			}
			if err == constant.SkipErrmsgABITypeNoFound {
				logger.Warn("skip tx msg for " + err.Error() + fmt.Sprint(" height: ", tx.Height, " txHash: ", tx.TxHash))
				continue
			}
			if err == constant.SkipErrmsgNoSupport {
				logger.Warn("skip tx msg for " + err.Error() + fmt.Sprint(" height: ", tx.Height, " txHash: ", tx.TxHash))
				continue
			}
			return err
		}
		ddcsInfo = append(ddcsInfo, ddcinfos...)
		evmDatas = append(evmDatas, evmdatas...)

	}

	txInfo := repository.ExSyncTxEvm{
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
	return d.ddcTxInfoModel.Save(txInfo)
}

//func (d SyncDdcTask) getDdc721InfoByDdcId(ddcId int64, msgEtheumTx *contracts.DocMsgEthereumTx) (*repository.ExSyncDdc, error) {
//	var (
//		ddcUri    string
//		ddcName   string
//		ddcSymbol string
//		ddcOwner  string
//		ddcData   string
//		err       error
//	)
//	switch msgEtheumTx.DdcType {
//	case contracts.ContractDDC721:
//		if ddcSymbol, err = ddc_sdk.Client().GetDDC721Service().Symbol(); err != nil {
//			return nil, err
//		}
//		if ddcName, err = ddc_sdk.Client().GetDDC721Service().Name(); err != nil {
//			return nil, err
//		}
//		if ddcUri, err = ddc_sdk.Client().GetDDC721Service().DdcURI(ddcId); err != nil {
//			return nil, err
//		}
//		if owner, err1 := ddc_sdk.Client().GetDDC721Service().OwnerOf(ddcId); err1 != nil {
//			return nil, err1
//		} else {
//			if ddcOwner, err = ddc_sdk.Client().GetDDC721Service().HexToBech32(owner); err != nil {
//				return nil, err
//			}
//		}
//
//		if len(msgEtheumTx.Inputs) > 1 {
//			ddcData = msgEtheumTx.Inputs[1]
//		}
//
//	}
//	data := &repository.ExSyncDdc{
//		DdcId:           ddcId,
//		DdcType:         contracts.DdcType[msgEtheumTx.DdcType],
//		DdcName:         ddcName,
//		DdcSymbl:        ddcSymbol,
//		ContractAddress: msgEtheumTx.ContractAddr,
//		DdcData:         ddcData,
//		Creator:         msgEtheumTx.Signer,
//		Owner:           ddcOwner,
//		//Amount:          amount,
//		DdcUri:         ddcUri,
//		LatestTxHeight: msgEtheumTx.TxHeight,
//		LatestTxTime:   msgEtheumTx.TxTime,
//	}
//	return data, nil
//}

func (d SyncDdcTask) createExSyncDdcByDdcId(ddcId int64, msgEtheumTx *contracts.DocMsgEthereumTx) *repository.ExSyncDdc {
	data := &repository.ExSyncDdc{
		DdcId:   ddcId,
		DdcType: contracts.DdcType[msgEtheumTx.DdcType],
		//DdcName:         ddcName,
		//DdcSymbl:        ddcSymbol,
		ContractAddress: msgEtheumTx.ContractAddr,
		//DdcData:         ddcData,
		Creator: msgEtheumTx.Signer,
		//Owner:           ddcOwner,
		//Amount:          amount,
		//DdcUri:         ddcUri,
		LatestTxHeight: msgEtheumTx.TxHeight,
		LatestTxTime:   msgEtheumTx.TxTime,
	}
	return data
}

func (d SyncDdcTask) getDdcTxsWithScope(latestHeight, maxHeight int64) []repository.Tx {
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

func (d SyncDdcTask) deleteDdcs(ddcId int64, contractAddr string) error {
	if err := d.syncDdcModel.Update(contractAddr, ddcId, bson.M{"is_delete": true}); err != nil && err != mgo.ErrNotFound {
		return err
	}
	return nil
}

func (d SyncDdcTask) editDdcs(ddcId uint64, contractAddr string, editData bson.M) error {
	if err := d.syncDdcModel.Update(contractAddr, int64(ddcId), editData); err != nil && err != mgo.ErrNotFound {
		return err
	}
	return nil
}

func (d SyncDdcTask) bitchSave(ddcs []*repository.ExSyncDdc) error {

	insertOps := make([]txn.Op, 0, repository.GetSrvConf().InsertBatchLimit)
	for _, val := range ddcs {
		op := txn.Op{
			C:      d.syncDdcModel.Name(),
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

func (d SyncDdcTask) getMaxHeight() (int64, error) {
	latestTx, err := d.txModel.FindLatestTx()
	if err != nil {
		return 0, err
	}
	return latestTx.Height, nil
}

func (d SyncDdcTask) getDdcLatestHeight() (int64, error) {
	ddc, err := d.syncDdcModel.DdcLatest()
	if err != nil && err != mgo.ErrNotFound {
		return 0, err
	}

	return ddc.LatestTxHeight, nil
}
