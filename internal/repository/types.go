package repository

import (
	"github.com/bianjieai/iobscan-explorer-backend/pkg/logger"
	"gopkg.in/mgo.v2"
)

var (
	docsShouldEnsureIndex = []DocV2{
		new(ExSyncDdc),
		new(ExSyncTxEvm),
		new(ExEvmContractsConfig),
		new(SyncTask),
		new(Tx),
	}
)

func EnsureIndexes() {
	if len(docsShouldEnsureIndex) == 0 {
		return
	}
	for _, v := range docsShouldEnsureIndex {
		v.EnsureIndexes()
	}
}

func ensureIndexes(collectionName string, indexes []mgo.Index) {
	session := GetSession()
	defer session.Close()
	c := session.DB(GetDataBase()).C(collectionName)
	if len(indexes) > 0 {
		for _, v := range indexes {
			if err := c.EnsureIndex(v); err != nil {
				logger.Warn("ensure index fail", logger.String("collectionName", collectionName),
					logger.String("err", err.Error()))
			}
		}
	}
}
