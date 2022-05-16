package repository

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionNameTx = "sync_tx"
	EthereumTxType   = "ethereum_tx"
	TxStatusSuccess  = 1
)

type (
	Tx struct {
		Time      int64       `bson:"time"`
		Height    int64       `bson:"height"`
		TxHash    string      `bson:"tx_hash"`
		Type      string      `bson:"type"` // parse from first msg
		Memo      string      `bson:"memo"`
		Status    uint32      `bson:"status"`
		Log       string      `bson:"log"`
		Fee       *Fee        `bson:"fee"`
		Types     []string    `bson:"types"`
		Events    []Event     `bson:"events"`
		EventsNew []EventNew  `bson:"events_new"`
		Signers   []string    `bson:"signers"`
		DocTxMsgs []TxMsg     `bson:"msgs"`
		Addrs     []string    `bson:"addrs"`
		TxIndex   uint32      `bson:"tx_index"`
		Ext       interface{} `bson:"ext"`
	}

	Event struct {
		Type       string   `bson:"type"`
		Attributes []KvPair `bson:"attributes"`
	}

	KvPair struct {
		Key   string `bson:"key"`
		Value string `bson:"value"`
	}

	EventNew struct {
		MsgIndex uint32  `bson:"msg_index" json:"msg_index"`
		Events   []Event `bson:"events"`
	}
	TxMsg struct {
		Type string      `bson:"type"`
		Msg  interface{} `bson:"msg"`
	}
)

type Fee struct {
	Amount []Coin `bson:"amount"`
	Gas    int64  `bson:"gas"`
}

type Coin struct {
	Denom  string `bson:"denom" json:"denom"`
	Amount string `bson:"amount" json:"amount"`
}

func (d Tx) Name() string {
	return CollectionNameTx
}

func (d Tx) EnsureIndexes() {
	var indexes []mgo.Index
	indexes = append(indexes, mgo.Index{
		Key:        []string{"-height", "-tx_hash"},
		Unique:     true,
		Background: true,
	})

	ensureIndexes(d.Name(), indexes)

}

func (d Tx) PkKvPair() map[string]interface{} {
	return bson.M{"height": d.Height, "tx_hash": d.TxHash}
}

func (d Tx) FindDdcTx(latestHeight int64) ([]Tx, error) {
	var res []Tx
	q := bson.M{
		//"status":    TxStatusSuccess,
		"msgs.type": EthereumTxType,
		"height": bson.M{
			"$gt":  latestHeight,
			"$lte": latestHeight + GetSrvConf().IncreHeight,
		},
	}
	fn := func(c *mgo.Collection) error {
		return c.Find(q).Sort([]string{"+height"}...).All(&res)
	}

	return res, ExecCollection(d.Name(), fn)
}

func (d Tx) FindLatestTx() (Tx, error) {
	var res Tx
	q := bson.M{
		//"status":    TxStatusSuccess,
		"msgs.type": EthereumTxType,
	}
	fn := func(c *mgo.Collection) error {
		return c.Find(q).Sort([]string{"-height"}...).One(&res)
	}

	return res, ExecCollection(d.Name(), fn)
}

func (d Tx) findAll(q map[string]interface{}, sorts []string, skip, limit int) ([]Tx, error) {

	var res []Tx

	fn := func(c *mgo.Collection) error {
		return c.Find(q).Sort(sorts...).Skip(skip).Limit(limit).All(&res)
	}

	return res, ExecCollection(d.Name(), fn)
}
