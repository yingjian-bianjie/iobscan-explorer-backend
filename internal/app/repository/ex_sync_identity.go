package repository

import (
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	CollectionNameExSyncIdentity = "ex_sync_identity"
)

func NewExSyncIdentityRepo(cli *qmgo.Client, database string) IExSyncIdentityRepo {
	return &exSyncIdentityRepo{coll: cli.Database(database).Collection(CollectionNameExSyncIdentity)}
}

type IExSyncIdentityRepo interface {
	QueryIdentityCount() (int64, error)
}

type exSyncIdentityRepo struct {
	coll *qmgo.Collection
}

func (repo *exSyncIdentityRepo) QueryIdentityCount() (int64, error) {
	count, err := repo.coll.Find(ctx, bson.M{}).Count()
	return count, err
}
