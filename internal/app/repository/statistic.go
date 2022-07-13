package repository

import (
	"context"

	"github.com/bianjieai/iobscan-explorer-backend/internal/app/model"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	TxType   string
	TxStatus int
)

const (
	DefineService TxType   = "define_service"
	Success       TxStatus = 1
	Failed        TxStatus = 0
)

const (
	CollectionNameExStatistic = "ex_statistics"
)

var (
	ctx             = context.Background()
	validatorStatus = map[string]int{
		"Unbonded":  1,
		"Unbonding": 2,
		"bonded":    3,
	}
)

func NewStatisticRepo(cli *qmgo.Client, database string) IStatisticRepo {
	return &statisticRepo{coll: cli.Database(database).Collection(CollectionNameExStatistic)}
}

type IStatisticRepo interface {
	FindStatisticsRecord(statisticName string) (*model.StatisticsType, error)
	InsertStatisticsRecord(statisticsType model.StatisticsType) error
	UpdateStatisticsRecord(statisticsType model.StatisticsType) error
}

type statisticRepo struct {
	coll *qmgo.Collection
}

func (repo *statisticRepo) UpdateStatisticsRecord(statisticsType model.StatisticsType) error {
	q := bson.M{
		"statistics_name": statisticsType.StatisticsName,
	}
	update := bson.M{
		"$set": bson.M{
			"count":     statisticsType.Count,
			"update_at": statisticsType.UpdateAt,
			"data":      statisticsType.Data,
		},
	}
	err := repo.coll.UpdateOne(ctx, q, update)
	return err
}

func (repo *statisticRepo) InsertStatisticsRecord(statisticsType model.StatisticsType) error {
	_, err := repo.coll.InsertOne(ctx, statisticsType)
	return err
}

func (repo *statisticRepo) FindStatisticsRecord(statisticName string) (*model.StatisticsType, error) {
	var statistic model.StatisticsType
	q := bson.M{
		"statistics_name": statisticName,
	}
	err := repo.coll.Find(ctx, q).One(&statistic)
	return &statistic, err
}
