package repository

import (
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/model"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	CollectionExTxType = "ex_tx_type"
)

func NewExTxTypeRepo(cli *qmgo.Client, database string) IExTxTypeRepo {
	return &exTxTypeRepo{coll: cli.Database(database).Collection(CollectionExTxType)}
}

type IExTxTypeRepo interface {
	QueryTxTypeList() ([]string, error)
}

type exTxTypeRepo struct {
	coll *qmgo.Collection
}

func (repo *exTxTypeRepo) QueryTxTypeList() ([]string, error) {
	var typeList []model.ExTxType
	err := repo.coll.Find(ctx, bson.M{}).All(&typeList)
	var txTypeList = make([]string, 0, len(typeList))
	for _, txType := range typeList {
		txTypeList = append(txTypeList, txType.TypeName)
	}
	return txTypeList, err
}
