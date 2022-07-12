package repository

import (
	"context"

	"github.com/bianjieai/iobscan-explorer-backend/internal/app/enum"
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/model"
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
	CollectionNameSyncTx             = "sync_tx"
	CollectionNameExSyncNft          = "ex_sync_nft"
	CollectionNameExSyncValidator    = "ex_sync_validator"
	CollectionNameExSyncIdentity     = "ex_sync_identity"
	CollectionNameExSyncDenom        = "ex_sync_denom"
	CollectionNameExStakingValidator = "ex_staking_validator"
	CollectionNameExTokens           = "ex_tokens"
	CollectionNameExStatistic        = "ex_statistics"
	CollectionNameSyncTask           = "sync_task"
)

var (
	ctx             = context.Background()
	validatorStatus = map[string]int{
		"Unbonded":  1,
		"Unbonding": 2,
		"bonded":    3,
	}
)

func NewStatisticRepo() IStatisticRepo {
	return &statisticRepo{}
}

type IStatisticRepo interface {
	QueryDenomCount(database string) (int64, error)
	QueryAssetCount(database string) (int64, error)
	QueryIdentityCount(database string) (int64, error)
	QueryServiceCount(database string) (int64, error)
	QueryConsensusValidatorCount(database string) (int64, error)
	QueryValidatorNumCount(database string) (int64, error)
	FindStatisticsRecord(database string, statisticName string) (*model.StatisticsType, error)
	InsertStatisticsRecord(database string, statisticsType model.StatisticsType) error
	UpdateStatisticsRecord(database string, statisticsType model.StatisticsType) error
	GetTaskStatus(database string, height int64, status string) (int64, error)
	QueryLatestHeight(database string, height int64) (*model.SyncTx, error)
	QueryTxCountWithHeight(database string, height int64) (int64, error)
}

type statisticRepo struct {
}

func (repo *statisticRepo) QueryTxCountWithHeight(database string, height int64) (int64, error) {
	q := bson.M{
		"height": height,
	}
	count, err := MgoCli.Database(database).Collection(CollectionNameSyncTx).Find(ctx, q).Count()
	return count, err
}

func (repo *statisticRepo) QueryLatestHeight(database string, height int64) (*model.SyncTx, error) {
	var syncTx model.SyncTx
	q := bson.M{
		"height": bson.M{
			"$gte": height,
		},
	}
	err := MgoCli.Database(database).Collection(CollectionNameSyncTx).Find(ctx, q).Sort("-height").One(&syncTx)
	return &syncTx, err
}

func (repo *statisticRepo) GetTaskStatus(database string, height int64, status string) (int64, error) {
	q := bson.M{
		"end_height": height,
		"status":     status,
	}

	count, err := MgoCli.Database(database).Collection(CollectionNameSyncTask).Find(ctx, q).Count()
	return count, err
}

func (repo *statisticRepo) UpdateStatisticsRecord(database string, statisticsType model.StatisticsType) error {
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
	err := MgoCli.Database(database).Collection(CollectionNameExStatistic).UpdateOne(ctx, q, update)
	return err
}

func (repo *statisticRepo) InsertStatisticsRecord(database string, statisticsType model.StatisticsType) error {
	_, err := MgoCli.Database(database).Collection(CollectionNameExStatistic).InsertOne(ctx, statisticsType)
	return err
}

func (repo *statisticRepo) FindStatisticsRecord(database string, statisticName string) (*model.StatisticsType, error) {
	var statistic model.StatisticsType
	q := bson.M{
		"statistics_name": statisticName,
	}
	err := MgoCli.Database(database).Collection(CollectionNameExStatistic).Find(ctx, q).One(&statistic)
	return &statistic, err
}

func (repo *statisticRepo) QueryValidatorNumCount(database string) (int64, error) {
	q := bson.M{
		"status": validatorStatus["bonded"],
		"jailed": false,
	}
	count, err := MgoCli.Database(database).Collection(CollectionNameExStakingValidator).Find(ctx, q).Count()
	return count, err
}

func (repo *statisticRepo) QueryConsensusValidatorCount(database string) (int64, error) {
	count, err := MgoCli.Database(database).Collection(CollectionNameExSyncValidator).Find(ctx, bson.M{}).Count()
	return count, err
}

func (repo *statisticRepo) QueryServiceCount(database string) (int64, error) {
	q := bson.M{
		"msg.type": enum.DefineService,
		"status":   enum.Success,
	}
	count, err := MgoCli.Database(database).Collection(CollectionNameSyncTx).Find(ctx, q).Count()
	return count, err
}

func (repo *statisticRepo) QueryIdentityCount(database string) (int64, error) {
	count, err := MgoCli.Database(database).Collection(CollectionNameExSyncIdentity).Find(ctx, bson.M{}).Count()
	return count, err
}

func (repo *statisticRepo) QueryAssetCount(database string) (int64, error) {
	count, err := MgoCli.Database(database).Collection(CollectionNameExSyncNft).Find(ctx, bson.M{}).Count()
	return count, err
}

func (repo *statisticRepo) QueryDenomCount(database string) (int64, error) {
	count, err := MgoCli.Database(database).Collection(CollectionNameExSyncDenom).Find(ctx, bson.M{}).Count()
	return count, err
}
