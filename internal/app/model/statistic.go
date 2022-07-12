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
