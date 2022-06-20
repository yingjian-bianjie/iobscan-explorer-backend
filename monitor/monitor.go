package monitor

import (
	"context"
	"github.com/bianjieai/ddc-sdk-go/ddc-sdk-platform-go/config"
	"github.com/bianjieai/iobscan-explorer-backend/internal/repository"
	"github.com/bianjieai/iobscan-explorer-backend/monitor/metrics"
	"github.com/bianjieai/iobscan-explorer-backend/pkg/logger"
	"os"
	"time"
)

var (
	sdkClientStatusMetric metrics.Guage
	cronTaskStatusMetric  metrics.Guage
)

func NewMetricCronWorkStatus() metrics.Guage {
	syncDdcWorkStatusMetric := metrics.NewGuage(
		"ddc_parser",
		"cron_task",
		"working_status",
		"ddc_parser_cron_task_working_status cron task working status (1:Normal  -1:UNormal)",
		[]string{"taskName"},
	)
	syncDdcWorkStatus, _ := metrics.CovertGuage(syncDdcWorkStatusMetric)
	return syncDdcWorkStatus
}

func NewMetricDdcSdkClientStatus() metrics.Guage {
	ddcSdkClientStatusMetric := metrics.NewGuage(
		"ddc_parser",
		"ddc_sdk_go",
		"client_status",
		"ddc_parser_ddc_sdk_go_client_status ddc-sdk-go client connection  status (1:Normal  -1:UNormal)",
		nil,
	)
	ddcSdkClientStatus, _ := metrics.CovertGuage(ddcSdkClientStatusMetric)
	return ddcSdkClientStatus
}

func sdkClientStatus() {
	for {
		t := time.NewTimer(time.Duration(30) * time.Second)
		select {
		case <-t.C:
			getSdkClientStatus()
		}
	}
}
func getSdkClientStatus() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()
	_, err := config.Info.Conn().NetworkID(ctx)
	if err != nil {
		logger.Warn(err.Error(), logger.String("funcName", "getSdkClientStatus"))
		sdkClientStatusMetric.Set(float64(-1))
	} else {
		sdkClientStatusMetric.Set(float64(1))
	}

	return
}

func SetCronTaskStatusMetricValue(taskName string, value float64) {
	cronTaskStatusMetric.With("taskName", taskName).Set(value)
}

func Start() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("monitor server occur error", logger.Any("err", err))
			os.Exit(1)
		}
	}()
	logger.Info("monitor server start")
	// start monitor
	server := metrics.NewMonitor(repository.GetSrvConf().Prometheus)
	sdkClientStatusMetric = NewMetricDdcSdkClientStatus()
	cronTaskStatusMetric = NewMetricCronWorkStatus()

	server.Report(func() {
		sdkClientStatus()
	})
}
