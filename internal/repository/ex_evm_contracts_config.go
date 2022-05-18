package repository

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionNameExEvmContractsConfig = "ex_evm_contracts_config"
)

type ExEvmContractsConfig struct {
	Type       int    `bson:"type"`
	Name_      string `bson:"name"`
	Address    string `bson:"address"`
	AbiContent string `bson:"abi_content"`
}

func (d ExEvmContractsConfig) Name() string {
	return CollectionNameExEvmContractsConfig
}

func (d ExEvmContractsConfig) PkKvPair() map[string]interface{} {
	return bson.M{"address": d.Address}
}

func (d ExEvmContractsConfig) EnsureIndexes() {
	var indexes []mgo.Index
	indexes = append(indexes,
		mgo.Index{
			Key:        []string{"-address"},
			Unique:     true,
			Background: true,
		},
	)

	ensureIndexes(d.Name(), indexes)
}

func (d ExEvmContractsConfig) FindAll() ([]ExEvmContractsConfig, error) {

	var res []ExEvmContractsConfig

	fn := func(c *mgo.Collection) error {
		return c.Find(bson.M{}).All(&res)
	}

	return res, ExecCollection(d.Name(), fn)
}
