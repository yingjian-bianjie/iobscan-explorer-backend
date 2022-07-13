package repository

import (
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	CollectionNameExSyncValidator = "ex_sync_validator"
)

func NewExSyncValidatorRepo(cli *qmgo.Client, database string) IExSyncValidatorRepo {
	return &exSyncValidatorRepo{coll: cli.Database(database).Collection(CollectionNameExSyncValidator)}
}

type IExSyncValidatorRepo interface {
	QueryConsensusValidatorCount() (int64, error)
}

type exSyncValidatorRepo struct {
	coll *qmgo.Collection
}

func (repo *exSyncValidatorRepo) QueryConsensusValidatorCount() (int64, error) {
	count, err := repo.coll.Find(ctx, bson.M{}).Count()
	return count, err
}
