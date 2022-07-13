package model

type StatisticsType struct {
	StatisticsName string `bson:"statistics_name"`
	Count          int64  `bson:"count"`
	Data           string `bson:"data"`
	StatisticsInfo string `bson:"statistics_info"`
	CreateAt       int64  `bson:"create_at"`
	UpdateAt       int64  `bson:"update_at"`
}

type SyncTx struct {
	Time    int64     `bson:"time"`
	Height  int64     `bson:"height"`
	TxHash  string    `bson:"tx_hash"`
	Type    string    `bson:"type"`
	Memo    string    `bson:"memo"`
	Status  int64     `bson:"status"`
	Log     string    `bson:"log"`
	Fee     FeeAmount `bson:"fee"`
	Singers []string  `bson:"singers"`
}

type FeeAmount struct {
	Amount []DenomAmount `bson:"amount"`
	Gas    int64         `bson:"gas"`
}

type DenomAmount struct {
	Denom  string `bson:"denom"`
	Amount string `bson:"amount"`
}

type AllTxStatisticsInfoType struct {
	RecordHeight         int64 `json:"record_height"`
	RecordHeightBlockTxs int64 `json:"record_height_block_txs"`
}

type ExTxType struct {
	TypeName string `bson:"type_name"`
	TypeCn   string `bson:"type_cn"`
	TypeEn   string `bson:"type_en"`
}

type ExTxTypes []ExTxType

func (types ExTxTypes) GetExTypeList() []string {
	res := make([]string, 0, len(types))
	for _, v := range types {
		res = append(res, v.TypeName)
	}
	return res
}
