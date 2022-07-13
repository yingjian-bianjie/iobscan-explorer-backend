package repository

import (
	"context"
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/enum"
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/model"
	"github.com/qiniu/qmgo"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	CollectionNameSyncTx = "sync_tx"
)

func NewSyncTxRepo(cli *qmgo.Client, database string) ISyncTxRepo {
	return &syncTxRepo{coll: cli.Database(database).Collection(CollectionNameSyncTx)}
}

type ISyncTxRepo interface {
	QueryIncreTxCount(height int64) (int64, error)
	QueryTxCountStatistics() (int64, error)
	QueryTxCountWithHeight(height int64) (int64, error)
	QueryLatestHeight(height int64) (*model.SyncTx, error)
	QueryServiceCount() (int64, error)
}

type syncTxRepo struct {
	coll *qmgo.Collection
}

func (repo *syncTxRepo) QueryIncreTxCount(height int64) (int64, error) {
	list, err := ExTxTypeRepo.QueryTxTypeList()
	if err != nil {
		logrus.Errorf("QueryTxTypeList err:%s", err.Error())
	}
	q := bson.M{
		"msgs.type": bson.M{
			"$in": list,
		},
		"height": bson.M{
			"$gte": height,
		},
	}
	count, err := repo.coll.Find(context.Background(), q).Count()
	return count, err
}

func (repo *syncTxRepo) QueryTxCountStatistics() (int64, error) {
	list, err := ExTxTypeRepo.QueryTxTypeList()
	if err != nil {
		logrus.Errorf("QueryTxTypeList err:%s", err.Error())
	}
	m := bson.M{}
	m["msgs.type"] = bson.M{
		"$in": list,
	}
	count, err := repo.coll.Find(ctx, m).Count()
	return count, err
}

func (repo *syncTxRepo) QueryTxCountWithHeight(height int64) (int64, error) {
	q := bson.M{
		"height": height,
	}
	count, err := repo.coll.Find(ctx, q).Count()
	return count, err
}

func (repo *syncTxRepo) QueryLatestHeight(height int64) (*model.SyncTx, error) {
	var syncTx model.SyncTx
	q := bson.M{
		"height": bson.M{
			"$gte": height,
		},
	}
	err := repo.coll.Find(ctx, q).Sort("-height").One(&syncTx)
	return &syncTx, err
}

func (repo *syncTxRepo) QueryServiceCount() (int64, error) {
	q := bson.M{
		"msg.type": enum.DefineService,
		"status":   enum.Success,
	}
	count, err := repo.coll.Find(ctx, q).Count()
	return count, err
}
