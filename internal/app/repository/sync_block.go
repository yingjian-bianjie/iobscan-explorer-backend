package repository

import (
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/model"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	CollectionNameSyncBlock = "sync_block"
)

func NewSyncBlockRepo(cli *qmgo.Client, database string) ISyncBlockRepo {
	return &syncBlockRepo{coll: cli.Database(database).Collection(CollectionNameSyncBlock)}
}

type ISyncBlockRepo interface {
	QueryLatestBlockHeight(height int64) (*model.SyncBlock, error)
}

type syncBlockRepo struct {
	coll *qmgo.Collection
}

func (repo *syncBlockRepo) QueryLatestBlockHeight(height int64) (*model.SyncBlock, error) {
	var block model.SyncBlock
	q := bson.M{
		"height": bson.M{
			"$gte": height,
		},
	}

	err := repo.coll.Find(ctx, q).Sort("-height").One(&block)
	return &block, err
}
