package repository

import "github.com/qiniu/qmgo"

var (
	SyncTxRepo             ISyncTxRepo
	ExTxTypeRepo           IExTxTypeRepo
	ExSyncNftRepo          IExSyncNftRepo
	ExSyncValidatorRepo    IExSyncValidatorRepo
	ExSyncIdentityRepo     IExSyncIdentityRepo
	ExSyncDenomRepo        IExSyncDenomRepo
	ExStakingValidatorRepo IExStakingValidatorRepo
	SyncTaskRepo           ISyncTaskRepo
	StatisticRepo          IStatisticRepo
	ExTokensRepo           IExTokensRepo
	SyncBlockRepo          ISyncBlockRepo
)

func InitRepo(cli *qmgo.Client, database string) {
	SyncTxRepo = NewSyncTxRepo(cli, database)
	ExTxTypeRepo = NewExTxTypeRepo(cli, database)
	ExSyncNftRepo = NewExSyncNftRepo(cli, database)
	ExSyncValidatorRepo = NewExSyncValidatorRepo(cli, database)
	ExSyncIdentityRepo = NewExSyncIdentityRepo(cli, database)
	ExSyncDenomRepo = NewExSyncDenomRepo(cli, database)
	ExStakingValidatorRepo = NewExStakingValidatorRepo(cli, database)
	SyncTaskRepo = NewSyncTaskRepo(cli, database)
	StatisticRepo = NewStatisticRepo(cli, database)
	ExTokensRepo = NewExTokensRepo(cli, database)
	SyncBlockRepo = NewSyncBlockRepo(cli, database)
}
