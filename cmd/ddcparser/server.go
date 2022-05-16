package main

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/bianjieai/iobscan-explorer-backend/internal/common/configs"
	"github.com/bianjieai/iobscan-explorer-backend/internal/repository"
	"github.com/bianjieai/iobscan-explorer-backend/internal/task"
	ddc_sdk "github.com/bianjieai/iobscan-explorer-backend/pkg/libs/ddc-sdk"
	"github.com/bianjieai/iobscan-explorer-backend/pkg/logger"
)

func start(config *configs.Config) {
	runtime.GOMAXPROCS(runtime.NumCPU() / 2)
	c := make(chan os.Signal)
	defer func() {
		logger.Info("System Exit")

		repository.Stop()

		if err := recover(); err != nil {
			logger.Error("occur error", logger.Any("err", err))
			os.Exit(1)
		}
	}()

	logger.Debug("config info", logger.Any("conf", config))

	repository.Start(config)
	repository.EnsureIndexes()

	ddc_sdk.InitDDCSDKClient(config.DdcClient)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	task.Start(config)
	<-c
}
