package repository

import (
	"context"
	"fmt"

	conf "github.com/bianjieai/iobscan-explorer-backend/internal/app/config"
	"github.com/qiniu/qmgo"
)

var MgoCli *qmgo.Client

func InitQMgo(config *conf.MongoDB) {
	var maxPoolSize uint64 = 4096
	dsn := fmt.Sprintf("mongodb://%s:%s@%s:%d/?connect=direct&authSource=%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
	cli, err := qmgo.NewClient(context.Background(), &qmgo.Config{
		Uri:         dsn,
		MaxPoolSize: &maxPoolSize,
	})
	if err != nil {
		fmt.Println(err)
	}
	MgoCli = cli
}

func GetClient() *qmgo.Client {
	return MgoCli
}

func MongoDbStatus() bool {
	err := MgoCli.Ping(5)
	if err != nil {
		return false
	}
	return true
}
