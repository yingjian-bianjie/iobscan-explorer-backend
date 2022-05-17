package repository

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	CollectionNameExSyncDdc = "ex_sync_ddc"
)

type ExSyncDdc struct {
	ID              bson.ObjectId `bson:"_id"`
	DdcId           int64         `bson:"ddc_id"`
	DdcType         int           `bson:"ddc_type"`
	DdcSymbl        string        `bson:"ddc_symbl"`
	DdcName         string        `bson:"ddc_name"`
	ContractAddress string        `bson:"contract_address"`
	DdcUri          string        `bson:"ddc_uri"`
	Owner           string        `bson:"owner"`
	Creator         string        `bson:"creator"`
	Amount          int           `bson:"amount"`
	DdcData         string        `bson:"ddc_data"`
	LatestTxTime    int64         `bson:"latest_tx_time"`
	LatestTxHeight  int64         `bson:"latest_tx_height"`
	IsDelete        bool          `bson:"is_delete"`
	IsFreeze        bool          `bson:"is_freeze"`
	CreateAt        int64         `bson:"create_at"`
	UpdateAt        int64         `bson:"update_at"`
}

func (d ExSyncDdc) Name() string {
	return CollectionNameExSyncDdc
}

func (d ExSyncDdc) PkKvPair() map[string]interface{} {
	return bson.M{"contract_address": d.ContractAddress, "ddc_id": d.DdcId}
}

func (d ExSyncDdc) EnsureIndexes() {
	var indexes []mgo.Index
	indexes = append(indexes,
		mgo.Index{
			Key:        []string{"-contract_address", "-ddc_id"},
			Unique:     true,
			Background: true,
		},
	)

	ensureIndexes(d.Name(), indexes)
}

func (d ExSyncDdc) findAll(q map[string]interface{}, sorts []string,
	skip, limit int) ([]ExSyncDdc, error) {

	var res []ExSyncDdc

	fn := func(c *mgo.Collection) error {
		return c.Find(q).Sort(sorts...).Skip(skip).Limit(limit).All(&res)
	}

	return res, ExecCollection(d.Name(), fn)
}

func (d ExSyncDdc) findOne(q map[string]interface{}, sorts []string) (ExSyncDdc, error) {
	var res ExSyncDdc
	fn := func(c *mgo.Collection) error {
		return c.Find(q).Sort(sorts...).One(&res)
	}

	return res, ExecCollection(d.Name(), fn)
}

func (d ExSyncDdc) count(q map[string]interface{}) (int, error) {
	var num int
	fn := func(c *mgo.Collection) error {
		n, err := c.Find(q).Count()
		num = n
		return err
	}

	return num, ExecCollection(d.Name(), fn)
}

func (d ExSyncDdc) Save(ddc ExSyncDdc) error {
	ddc.ID = bson.NewObjectId()
	ddc.CreateAt = time.Now().Unix()
	return Save(&ddc)
}

func (d ExSyncDdc) DdcLatest() (ExSyncDdc, error) {
	q := bson.M{}
	return d.findOne(q, []string{"-latest_tx_height"})
}

func (d ExSyncDdc) UpdateOwnerOrUri(contractAddr string, ddcId int64, owner, uri string) error {
	editData := bson.M{}
	if owner != "" {
		editData["owner"] = owner
	}
	if uri != "" {
		editData["ddc_uri"] = uri
	}
	fn := func(c *mgo.Collection) error {
		return c.Update(bson.M{
			"contract_address": contractAddr,
			"ddc_id":           ddcId,
		}, bson.M{
			"$set": editData,
		})
	}
	return ExecCollection(d.Name(), fn)
}

func (d ExSyncDdc) UpdateNameAndSymbol(contractAddr string, ddcId int64, name, symbol string) error {
	editData := bson.M{}
	if name != "" {
		editData["ddc_name"] = name
	}
	if symbol != "" {
		editData["ddc_symbol"] = symbol
	}
	fn := func(c *mgo.Collection) error {
		return c.Update(bson.M{
			"contract_address": contractAddr,
			"ddc_id":           ddcId,
		}, bson.M{
			"$set": editData,
		})
	}
	return ExecCollection(d.Name(), fn)
}

func (d ExSyncDdc) DeleteDdc(contractAddr string, ddcId int64) error {
	fn := func(c *mgo.Collection) error {
		return c.Update(bson.M{
			"contract_address": contractAddr,
			"ddc_id":           ddcId,
		}, bson.M{
			"$set": bson.M{"is_delete": true},
		})
	}
	return ExecCollection(d.Name(), fn)
}

func (d ExSyncDdc) UpdateDdcLatestTxHeight(contractAddr string, ddcId int64, latestTxHeight int64) error {
	fn := func(c *mgo.Collection) error {
		return c.Update(bson.M{
			"contract_address": contractAddr,
			"ddc_id":           ddcId,
		}, bson.M{
			"$set": bson.M{"latest_tx_height": latestTxHeight},
		})
	}
	return ExecCollection(d.Name(), fn)
}
