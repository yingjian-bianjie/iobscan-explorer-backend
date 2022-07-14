package task

import (
	"encoding/json"
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/enum"
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/lcd"
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/model"
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/repository"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

const (
	NoDocument             = "mongo: no documents in result"
	TxMsgsAll              = "tx_msgs_all"
	SyncTaskStatusUnderway = "underway"
	CommunityPool          = "community_pool"
	TxAll                  = "tx_all"
	ServiceAll             = "service_all"
	ValidatorAll           = "validator_all"
	ValidatorActive        = "validator_active"
	IdentityAll            = "identity_all"
	NftAll                 = "nft_all"
	DenomAll               = "denom_all"
	BondedTokens           = "bonded_tokens"
	TotalSupply            = "total_supply"
	AccountsAll            = "accounts_all"
)

var (
	StatisticsNames = []string{TxAll, ServiceAll, ValidatorAll,
		ValidatorActive, IdentityAll, NftAll,
		DenomAll, BondedTokens, TotalSupply,
		CommunityPool, AccountsAll}

	StatisticsNameMap = map[string]int64{
		TxAll:           0,
		ServiceAll:      0,
		ValidatorAll:    0,
		ValidatorActive: 0,
		IdentityAll:     0,
		NftAll:          0,
		DenomAll:        0,
		BondedTokens:    0,
		TotalSupply:     0,
		AccountsAll:     0,
	}
)

type StatisticTask struct {
}

func (s *StatisticTask) Cron() string {
	return StatisticsExecuteTime
}

func (s *StatisticTask) Name() string {
	return StatisticTaskName
}

func (s *StatisticTask) Run() {
	//计算交易增量
	updateIncreTxCount()

	//交易消息
	updateIncreMsgsCount()

	jobs := global.Config.Task.CronJobs
	if strings.Contains(jobs, string(enum.Denom)) {
		count, err := repository.ExSyncDenomRepo.QueryDenomCount()
		if err != nil {
			logrus.Errorf("statistic_task QueryDenomCount err:%s", err.Error())
		}
		StatisticsNameMap[DenomAll] = count
	}
	if strings.Contains(jobs, string(enum.Nft)) {
		count, err := repository.ExSyncNftRepo.QueryAssetCount()
		if err != nil {
			logrus.Errorf("statistic_task QueryAssetCount err:%s", err.Error())
		}
		StatisticsNameMap[NftAll] = count
	}

	if strings.Contains(jobs, string(enum.Identity)) {
		count, err := repository.ExSyncIdentityRepo.QueryIdentityCount()
		if err != nil {
			logrus.Errorf("statistic_task QueryIdentityCount err:%s", err.Error())
		}
		StatisticsNameMap[IdentityAll] = count
	}

	if strings.Contains(jobs, string(enum.TxServiceName)) {
		count, err := repository.SyncTxRepo.QueryServiceCount()
		if err != nil {
			logrus.Errorf("statistic_task QueryServiceCount err:%s", err.Error())
		}
		StatisticsNameMap[ServiceAll] = count
	}

	if strings.Contains(jobs, string(enum.Validaters)) {
		count, err := repository.ExSyncValidatorRepo.QueryConsensusValidatorCount()
		if err != nil {
			logrus.Errorf("statistic_task QueryConsensusValidatorCount err:%s", err.Error())
		}
		StatisticsNameMap[ValidatorAll] = count
	}

	if strings.Contains(jobs, string(enum.StakingSyncValidatorsInfo)) {
		count, err := repository.ExStakingValidatorRepo.QueryValidatorNumCount()
		if err != nil {
			logrus.Errorf("statistic_task QueryIdentityCount err:%s", err.Error())
		}
		StatisticsNameMap[ValidatorActive] = count
	}
	if strings.Contains(jobs, string(enum.Tokens)) {
		tokens, err := lcd.GetBondedTokens()
		if err != nil {
			logrus.Errorf("statistic_task GetBondedTokens err: %s", err.Error())
		}
		supply, err := lcd.GetTotalSupply()
		if err != nil {
			logrus.Errorf("statistic_task GetTotalSupply err: %s", err.Error())
		}
		token, err := repository.ExTokensRepo.QueryMainToken()
		if err != nil && err.Error() != NoDocument {
			logrus.Errorf("statistic_task QueryMainToken err: %s", err.Error())
		}
		if token.IsMainToken {
			if len(supply.Supply) > 0 {
				for _, value := range supply.Supply {
					if value.Denom == token.Denom {
						amount, _ := strconv.ParseInt(value.Amount, 10, 64)
						StatisticsNameMap[TotalSupply] = amount
					}
				}
			}
		}
		bondedTokens, _ := strconv.ParseInt(tokens.Pool.BondedTokens, 10, 64)
		StatisticsNameMap[BondedTokens] = bondedTokens
	}

	pool, err := lcd.GetCommunityPool()
	//accounts_all
	if err != nil {
		logrus.Errorf("statistic_task GetCommunityPool err:%s", err.Error())
	}
	var community lcd.CommunityPool
	if len(pool.Pool) > 0 {
		for i, value := range pool.Pool {
			amount := strings.Split(value.Amount, ".")
			community.Pool[i].Denom = value.Denom
			if len(amount) > 0 {
				community.Pool[i].Amount = amount[0]
			}
		}
	}

	for _, name := range StatisticsNames {
		var data string
		if name == CommunityPool && len(community.Pool) > 0 {
			bytes, _ := json.Marshal(community.Pool)
			data = string(bytes)
		}
		statisticsRecord, err := repository.StatisticRepo.FindStatisticsRecord(name)
		if err != nil && err.Error() != NoDocument {
			logrus.Errorf("statistic_task FindStatisticsRecord:%s", err.Error())
		}

		if statisticsRecord.StatisticsName != "" {
			//更新
			if value, ok := StatisticsNameMap[name]; ok {
				statisticsRecord.Count = value
			}
			statisticsRecord.UpdateAt = time.Now().Unix()
			statisticsRecord.Data = data
			err := repository.StatisticRepo.UpdateStatisticsRecord(*statisticsRecord)
			if err != nil {
				logrus.Errorf("statistic_task UpdateStatisticsRecord err:%s", err)
			}
		} else {
			statisticsType := model.StatisticsType{
				StatisticsName: name,
				Count:          StatisticsNameMap[name],
				Data:           data,
				StatisticsInfo: "",
				CreateAt:       time.Now().Unix(),
				UpdateAt:       time.Now().Unix(),
			}
			err := repository.StatisticRepo.InsertStatisticsRecord(statisticsType)
			if err != nil {
				logrus.Errorf("statistic_task InsertStatisticsRecord err:%s", err)
			}
		}

	}
}

func updateIncreMsgsCount() {
	taskCount, err := repository.SyncTaskRepo.GetTaskStatus(0, SyncTaskStatusUnderway)
	if err != nil {
		logrus.Errorf("updateIncreMsgsCount query taskStatus err: %s", err.Error())
	}
	if taskCount > 0 {
		statisticsRecord, err := repository.StatisticRepo.FindStatisticsRecord(TxMsgsAll)
		if err != nil {
			logrus.Errorf("statistic_task findStatisticsRecord err: %s, statisticName:%s", err.Error(), TxMsgsAll)
		}

		if statisticsRecord.StatisticsName == "" {
			//新增
			statisticsType := model.StatisticsType{
				StatisticsName: "tx_msgs_all",
				Count:          0,
				Data:           "",
				StatisticsInfo: "",
				CreateAt:       time.Now().Unix(),
				UpdateAt:       time.Now().Unix(),
			}

			statistics := handleTxMsgsIncre(&statisticsType)
			//插入
			err = repository.StatisticRepo.InsertStatisticsRecord(*statistics)
			if err != nil {
				logrus.Errorf("statistic_task updateIncreMsgsCount InsertStatisticsRecord err：%s, statisticsType:%v", err.Error(), statistics)
			}
		} else {
			var increMsgsCnt int64 = 0
			if statisticsRecord.StatisticsInfo != "" {
				var msgsInfo model.TxMsgsInfo
				err := json.Unmarshal([]byte(statisticsRecord.StatisticsInfo), &msgsInfo)
				if err != nil {
					logrus.Errorf("StatisticsInfo Unmarshal err: %s", err.Error())
				}
				//大于等于这个块高的所有txMsgs数量
				increData, err := repository.SyncTxRepo.QueryTxMsgsIncre(msgsInfo.RecordHeight)
				var incre int64 = 0
				if len(increData) > 0 {
					incre = increData[0].Count
				}
				if incre > 0 && incre > msgsInfo.RecordHeightTxMsgs {
					//统计增量
					increMsgsCnt = incre - msgsInfo.RecordHeightTxMsgs
				}
				//查询最新的块高
				latestOneTx, err := repository.SyncTxRepo.QueryLatestHeight(msgsInfo.RecordHeight)
				if err != nil {
					logrus.Errorf("StatisticsInfo QueryLatestHeight err: %s", err.Error())
				}
				if latestOneTx.Height != 0 {
					var msgsInfo model.TxMsgsInfo
					msgsInfo.RecordHeight = latestOneTx.Height
					//查询最新块高下单交易消息数量
					increData, err := repository.SyncTxRepo.QueryTxMsgsCountByHeight(latestOneTx.Height)
					if err != nil {
						logrus.Errorf("StatisticsInfo QueryTxMsgsCountByHeight err: %s", err.Error())
					}
					if len(increData) > 0 {
						msgsInfo.RecordHeightTxMsgs = increData[0].Count
					}

					bytes, _ := json.Marshal(msgsInfo)
					statisticsRecord.StatisticsInfo = string(bytes)
				}

				//交易消息总数 = 统计增量数据 + 历史统计总数
				if increMsgsCnt > 0 {
					statisticsRecord.Count = increMsgsCnt + statisticsRecord.Count
					statisticsRecord.UpdateAt = time.Now().Unix()
					err := repository.StatisticRepo.UpdateStatisticsRecord(*statisticsRecord)
					if err != nil {
						logrus.Errorf("updateIncreMsgsCount StatisticsInfo UpdateStatisticsRecord err: %s，statisticsRecord：%v", err.Error(), statisticsRecord)
					}
				}
			} else {
				statistic := handleTxMsgsIncre(statisticsRecord)
				err := repository.StatisticRepo.UpdateStatisticsRecord(*statistic)
				if err != nil {
					logrus.Errorf("statistic_task UpdateStatisticsRecord err：%s", err.Error())
				}
			}
		}
	}
}

func handleTxMsgsIncre(statisticsType *model.StatisticsType) *model.StatisticsType {
	latestHeight, err := repository.SyncTxRepo.QueryLatestHeight(1)
	if err != nil {
		logrus.Errorf("statistic_task QueryLatestHeight err：%s", err.Error())
	}
	if latestHeight.Height != 0 {
		//查询此块高下单txMsgs数量
		txMsgCount, err := repository.SyncTxRepo.QueryTxMsgsCountByHeight(latestHeight.Height)
		if err != nil {
			logrus.Errorf("statistic_task QueryTxMsgsCountByHeight err：%s", err.Error())
		}
		var msgsInfo model.TxMsgsInfo
		msgsInfo.RecordHeight = latestHeight.Height
		if len(txMsgCount) > 0 {
			msgsInfo.RecordHeightTxMsgs = txMsgCount[0].Count
		}
		bytes, _ := json.Marshal(msgsInfo)

		statisticsType.StatisticsInfo = string(bytes)

		//总txMsgs
		incre, err := repository.SyncTxRepo.QueryTxMsgsIncre(1)
		if err != nil {
			logrus.Errorf("statistic_task QueryTxMsgsIncre err：%s", err.Error())
		}
		if len(incre) > 0 {
			statisticsType.Count = incre[0].Count
		}
	}
	return statisticsType
}

func updateIncreTxCount() {
	taskCount, err := repository.SyncTaskRepo.GetTaskStatus(0, SyncTaskStatusUnderway)
	if err != nil {
		logrus.Errorf("statistic_task updateIncreTxCount query taskStatus err: %s", err.Error())
	}

	if taskCount > 0 {
		list, err := repository.ExTxTypeRepo.QueryTxTypeList()
		if err != nil {
			logrus.Errorf("statistic_task updateIncreTxCount QueryTxTypeList: %s", err.Error())
		}
		typeList := getExTypeList(list)

		statisticsRecord, err := repository.StatisticRepo.FindStatisticsRecord("tx_all")
		if err != nil && err.Error() != NoDocument {
			logrus.Errorf("statistic_task updateIncreTxCount findStatisticsRecord err: %s, statisticName:%s", err.Error(), "tx_all")
		}
		if statisticsRecord.StatisticsName == "" {
			statisticsType := model.StatisticsType{
				StatisticsName: "tx_all",
				Count:          0,
				Data:           "",
				StatisticsInfo: "",
				CreateAt:       time.Now().Unix(),
				UpdateAt:       time.Now().Unix(),
			}
			//查询最新高度
			latestOneTx, err := repository.SyncTxRepo.QueryLatestHeight(1)
			if err != nil {
				logrus.Errorf("statistic_task updateIncreTxCount QueryLatestHeight err：%s", err.Error())
			}
			if latestOneTx.Height != 0 {
				//块高下的交易数
				txCount, err := repository.SyncTxRepo.QueryTxCountWithHeight(latestOneTx.Height)
				if err != nil {
					logrus.Errorf("statistic_task updateIncreTxCount QueryTxCountWithHeight err:%s,height:%v", err.Error(), latestOneTx.Height)
				}
				txAllInfo := model.AllTxStatisticsInfoType{
					RecordHeight:         latestOneTx.Height,
					RecordHeightBlockTxs: txCount,
				}
				bytes, _ := json.Marshal(txAllInfo)
				statisticsType.StatisticsInfo = string(bytes)
			}
			statistics, err := repository.SyncTxRepo.QueryTxCountStatistics(typeList)
			if err != nil {
				logrus.Errorf("statistic_task updateIncreTxCount QueryTxCountStatistics: %s", err.Error())
			}
			statisticsType.Count = statistics
			err = repository.StatisticRepo.InsertStatisticsRecord(statisticsType)
			if err != nil {
				logrus.Errorf("statistic_task updateIncreTxCount InsertStatisticsRecord err: %s, statisticsType:%v", err.Error(), statisticsType)
			}
		} else {
			var increTxCnt int64 = 0
			if statisticsRecord.StatisticsInfo != "" {
				var txAllInfo model.AllTxStatisticsInfoType
				err := json.Unmarshal([]byte(statisticsRecord.StatisticsInfo), &txAllInfo)
				if err != nil {
					logrus.Errorf("StatisticsInfo Unmarshal err: %s", err.Error())
				}
				if txAllInfo.RecordHeight != 0 && txAllInfo.RecordHeightBlockTxs != 0 {
					incre, err := repository.SyncTxRepo.QueryIncreTxCount(typeList, txAllInfo.RecordHeight)
					if err != nil {
						logrus.Errorf("statistic_task updateIncreTxCount QueryIncreTxCount err: %s", err.Error())
					}
					latestOneTx, err := repository.SyncTxRepo.QueryLatestHeight(txAllInfo.RecordHeight)
					if err != nil {
						logrus.Errorf("statistic_task updateIncreTxCount QueryLatestHeight err:%s", err.Error())
					}
					if incre > txAllInfo.RecordHeightBlockTxs {
						//统计增量数 = incre - record_height_block_txs
						increTxCnt = incre - txAllInfo.RecordHeightBlockTxs
					}
					if latestOneTx.Height != 0 {
						txAllInfo.RecordHeight = latestOneTx.Height
						txCount, err := repository.SyncTxRepo.QueryTxCountWithHeight(latestOneTx.Height)
						if err != nil {
							logrus.Errorf("statistic_task QueryTxCountWithHeight err:%s", err.Error())
						}
						txAllInfo.RecordHeightBlockTxs = txCount
						bytes, _ := json.Marshal(txAllInfo)
						statisticsRecord.StatisticsInfo = string(bytes)
					}
					//交易总数 = 统计增量数+历史统计总数
					if increTxCnt > 0 {
						statisticsRecord.Count = statisticsRecord.Count + increTxCnt
						statisticsRecord.UpdateAt = time.Now().Unix()
						err := repository.StatisticRepo.UpdateStatisticsRecord(*statisticsRecord)
						if err != nil {
							logrus.Errorf("statistic_task UpdateStatisticsRecord err: %s", err.Error())
						}
					}
				}
			} else {
				latestOneTx, err := repository.SyncTxRepo.QueryLatestHeight(1)
				if err != nil {
					logrus.Errorf("statistic_task UpdateStatisticsRecord err: %s", err.Error())
				}
				if latestOneTx.Height != 0 {
					txCount, err := repository.SyncTxRepo.QueryTxCountWithHeight(latestOneTx.Height)
					if err != nil {
						logrus.Errorf("statistic_task QueryTxCountWithHeight err: %s", err.Error())
					}
					infoType := model.AllTxStatisticsInfoType{
						RecordHeight:         latestOneTx.Height,
						RecordHeightBlockTxs: txCount,
					}
					bytes, _ := json.Marshal(infoType)
					statisticsRecord.StatisticsInfo = string(bytes)
					statistics, err := repository.SyncTxRepo.QueryTxCountStatistics(typeList)
					if err != nil {
						logrus.Errorf("statistic_task QueryTxCountStatistics err: %s", err.Error())
					}
					statisticsRecord.Count = statistics
					err = repository.StatisticRepo.UpdateStatisticsRecord(*statisticsRecord)
					if err != nil {
						logrus.Errorf("statistic_task UpdateStatisticsRecord err:%s", err.Error())
					}
				}
			}
		}
	} else {
		logrus.Info("ex_statistics:Catch-up status task suspended")
	}
}

func getExTypeList(txType []model.ExTxType) []string {
	res := make([]string, 0, len(txType))
	for _, v := range txType {
		res = append(res, v.TypeName)
	}
	return res
}
