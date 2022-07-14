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

type Tokens struct {
	Denom       string `bson:"denom"`
	IsMainToken bool   `bson:"is_main_token"`
	Chain       string `bson:"chain"`
}

type SyncBlock struct {
	Height int64 `json:"height"`
}

type TxMsgsCount struct {
	Count int64 `bson:"count"`
}

type TxMsgsInfo struct {
	RecordHeight       int64 `json:"record_height"`
	RecordHeightTxMsgs int64 `json:"record_height_tx_msgs"`
}
