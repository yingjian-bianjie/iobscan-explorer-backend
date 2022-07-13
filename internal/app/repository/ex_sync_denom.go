package repository

import (
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	CollectionNameExSyncDenom = "ex_sync_denom"
)

func NewExSyncDenomRepo(cli *qmgo.Client, database string) IExSyncDenomRepo {
	return &exSyncDenomRepo{coll: cli.Database(database).Collection(CollectionNameExSyncDenom)}
}

type IExSyncDenomRepo interface {
	QueryDenomCount() (int64, error)
}

type exSyncDenomRepo struct {
	coll *qmgo.Collection
}

func (repo *exSyncDenomRepo) QueryDenomCount() (int64, error) {
	count, err := repo.coll.Find(ctx, bson.M{}).Count()
	return count, err
}
