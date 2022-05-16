package repository

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	CollectionNameSyncTask = "sync_task"

	SyncTaskStatusUnderway = "underway"
)

type (
	SyncTask struct {
		ID             bson.ObjectId `bson:"_id"`
		StartHeight    int64         `bson:"start_height"`     // task start height
		EndHeight      int64         `bson:"end_height"`       // task end height
		CurrentHeight  int64         `bson:"current_height"`   // task current height
		Status         string        `bson:"status"`           // task status
		WorkerId       string        `bson:"worker_id"`        // worker id
		WorkerLogs     []WorkerLog   `bson:"worker_logs"`      // worker logs
		LastUpdateTime int64         `bson:"last_update_time"` // unix timestamp
	}

	WorkerLog struct {
		WorkerId  string    `bson:"worker_id"`  // worker id
		BeginTime time.Time `bson:"begin_time"` // time which worker begin handle this task
	}
)

func (d SyncTask) Name() string {
	return CollectionNameSyncTask
}

func (d SyncTask) EnsureIndexes() {

}

func (d SyncTask) PkKvPair() map[string]interface{} {
	return bson.M{"start_height": d.CurrentHeight, "end_height": d.EndHeight}
}

// query valid follow way
func (d SyncTask) QueryValidFollowTasks() (bool, error) {
	var syncTasks []SyncTask
	q := bson.M{}

	q["status"] = SyncTaskStatusUnderway

	q["end_height"] = bson.M{
		"$eq": 0,
	}

	fn := func(c *mgo.Collection) error {
		return c.Find(q).All(&syncTasks)
	}

	err := ExecCollection(d.Name(), fn)

	if err != nil {
		return false, err
	}

	if len(syncTasks) == 1 {
		return true, nil
	}

	return false, nil
}
