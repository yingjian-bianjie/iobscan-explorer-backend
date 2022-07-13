package repository

import (
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	CollectionNameExSyncNft = "ex_sync_nft"
)

func NewExSyncNftRepo(cli *qmgo.Client, database string) IExSyncNftRepo {
	return &exSyncNftRepo{coll: cli.Database(database).Collection(CollectionNameExSyncNft)}
}

type IExSyncNftRepo interface {
	QueryAssetCount() (int64, error)
}

type exSyncNftRepo struct {
	coll *qmgo.Collection
}

func (repo *exSyncNftRepo) QueryAssetCount() (int64, error) {
	count, err := repo.coll.Find(ctx, bson.M{}).Count()
	return count, err
}
