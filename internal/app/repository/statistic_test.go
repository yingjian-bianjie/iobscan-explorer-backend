package repository

import (
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/config"
	"testing"
)

func TestMain(m *testing.M) {
	cfg := config.MongoDB{
		Url:      "mongodb://wcchain:wcchainPassword@192.168.150.40:27017/?authSource=iobscan_wcchain",
		Database: "iobscan_wcchain",
	}
	InitQMgo(&cfg)
	m.Run()
}

func Test_statisticRepo_FindStatisticsRecord(t *testing.T) {
	record, err := StatisticRepo.FindStatisticsRecord("tx_all")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(record)
}
