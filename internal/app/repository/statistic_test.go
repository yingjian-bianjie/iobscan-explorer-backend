package repository

import (
	"testing"
	"time"

	"github.com/bianjieai/iobscan-explorer-backend/internal/app/config"
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/model"
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/utils"
)

func TestMain(m *testing.M) {
	cfg := config.MongoDB{
		Host:     "192.168.150.40",
		Port:     27017,
		Username: "wcchain",
		Password: "wcchainPassword",
		Database: "iobscan_wcchain",
	}
	InitQMgo(&cfg)
	m.Run()
}

func Test_statisticRepo_QueryDenomCount(t *testing.T) {

	count, err := StatisticRepo.QueryDenomCount("iobscan_wcchain")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(count)
}

func Test_statisticRepo_QueryAssetCount(t *testing.T) {
	count, err := StatisticRepo.QueryAssetCount("iobscan_wcchain")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(count)
}

func Test_statisticRepo_QueryIdentityCount(t *testing.T) {
	count, err := StatisticRepo.QueryIdentityCount("iobscan_wcchain")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(count)
}

func Test_statisticRepo_QueryConsensusValidatorCount(t *testing.T) {
	count, err := StatisticRepo.QueryConsensusValidatorCount("iobscan_wcchain")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(count)
}

func Test_statisticRepo_QueryValidatorNumCount(t *testing.T) {
	count, err := StatisticRepo.QueryValidatorNumCount("iobscan_wcchain")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(count)
}

func Test_statisticRepo_FindStatisticsRecord(t *testing.T) {
	record, err := StatisticRepo.FindStatisticsRecord("iobscan_wcchain", "tx_all")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(utils.MustMarshalJsonToStr(record))
}

func Test_statisticRepo_InsertStatisticsRecord(t *testing.T) {
	statisticsType := model.StatisticsType{
		StatisticsName: "test",
		Count:          222,
		Data:           "",
		StatisticsInfo: "",
		CreateAt:       time.Now().Unix(),
		UpdateAt:       time.Now().Unix(),
	}
	err := StatisticRepo.InsertStatisticsRecord("iobscan_wcchain", statisticsType)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_statisticRepo_QueryLatestHeight(t *testing.T) {
	height, err := StatisticRepo.QueryLatestHeight("iobscan_wcchain", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(height.Height)
}
