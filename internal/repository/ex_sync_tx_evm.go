package repository

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	CollectionNameExSyncTxEvm = "ex_sync_tx_evm"
)

type (
	ExSyncTxEvm struct {
		Time       int64     `bson:"time"`
		Height     int64     `bson:"height"`
		TxHash     string    `bson:"tx_hash"`
		Status     uint32    `bson:"status"`
		Fee        *Fee      `bson:"fee"`
		Types      []string  `bson:"types"`
		Signers    []string  `bson:"signers"`
		EvmDatas   []EvmData `bson:"evm_datas"`
		ExDdcInfos []DdcInfo `bson:"ex_ddc_infos"`
		CreateAt   int64     `bson:"create_at"`
		UpdateAt   int64     `bson:"update_at"`
	}
	EvmData struct {
		EvmTxHash       string    `bson:"evm_tx_hash"`
		EvmMethod       string    `bson:"evm_method"`
		TxReceipt       TxReceipt `json:"tx_receipt"`
		EvmInputs       []string  `bson:"evm_inputs"`
		EvmOutputs      []string  `bson:"evm_outputs"`
		DataType        string    `bson:"data_type"`
		ContractAddress string    `bson:"contract_address"`
	}
	TxReceipt struct {
		Status int64    `bson:"status"`
		Logs   []string `bson:"logs"`
	}
	DdcInfo struct {
		DdcId     int64  `bson:"ddc_id"`
		DdcName   string `bson:"ddc_name"`
		DdcSymbol string `bson:"ddc_symbol"`
		DdcType   string `bson:"ddc_type"`
		DdcUri    string `bson:"ddc_uri"`
		EvmTxHash string `bson:"evm_tx_hash"`
		Sender    string `bson:"sender"`
		Recipient string `bson:"recipient"`
	}
)

func (d ExSyncTxEvm) Name() string {
	return CollectionNameExSyncTxEvm
}

func (d ExSyncTxEvm) EnsureIndexes() {
	var indexes []mgo.Index
	indexes = append(indexes, mgo.Index{
		Key:        []string{"-height", "-tx_hash"},
		Unique:     true,
		Background: true,
	})

	ensureIndexes(d.Name(), indexes)

}

func (d ExSyncTxEvm) PkKvPair() map[string]interface{} {
	return bson.M{"height": d.Height, "tx_hash": d.TxHash}
}

func (d ExSyncTxEvm) Save(ddcTxInfo ExSyncTxEvm) error {
	ddcTxInfo.CreateAt = time.Now().Unix()
	return Save(&ddcTxInfo)
}

func (d ExSyncTxEvm) TxEvmLatest() (ExSyncTxEvm, error) {
	var res ExSyncTxEvm
	q := bson.M{}
	fn := func(c *mgo.Collection) error {
		return c.Find(q).Sort([]string{"-height"}...).One(&res)
	}
	return res, ExecCollection(d.Name(), fn)
}
