package repository

import (
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	CollectionNameSyncTask = "sync_task"
)

func NewSyncTaskRepo(cli *qmgo.Client, database string) ISyncTaskRepo {
	return &syncTaskRepo{coll: cli.Database(database).Collection(CollectionNameSyncTask)}
}

type ISyncTaskRepo interface {
	GetTaskStatus(height int64, status string) (int64, error)
}

type syncTaskRepo struct {
	coll *qmgo.Collection
}

func (repo *syncTaskRepo) GetTaskStatus(height int64, status string) (int64, error) {
	q := bson.M{
		"end_height": height,
		"status":     status,
	}

	count, err := repo.coll.Find(ctx, q).Count()
	return count, err
}
