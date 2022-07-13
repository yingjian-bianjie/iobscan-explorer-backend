package repository

import (
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	CollectionNameExStakingValidator = "ex_staking_validator"
)

func NewExStakingValidatorRepo(cli *qmgo.Client, database string) IExStakingValidatorRepo {
	return &exStakingValidatorRepo{coll: cli.Database(database).Collection(CollectionNameExStakingValidator)}
}

type IExStakingValidatorRepo interface {
	QueryValidatorNumCount() (int64, error)
}

type exStakingValidatorRepo struct {
	coll *qmgo.Collection
}

func (repo *exStakingValidatorRepo) QueryValidatorNumCount() (int64, error) {
	q := bson.M{
		"status": validatorStatus["bonded"],
		"jailed": false,
	}
	count, err := repo.coll.Find(ctx, q).Count()
	return count, err
}
