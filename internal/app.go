package app

import (
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/api"
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/config"
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/global"
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/lcd"
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/repository"
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/task"
	"github.com/sirupsen/logrus"
)

func Serve(cfg *config.Config) {
	global.Config = cfg
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   constant.DefaultTimeFormat,
		DisableHTMLEscape: true,
	})
	if level, err := logrus.ParseLevel(cfg.App.LogLevel); err == nil {
		logrus.SetLevel(level)
	}

	repository.InitQMgo(&cfg.Mongodb)

	lcd.Init(&cfg.Lcd)

	task.RegisterCron(&task.StatisticTask{})
	task.CronStart()

	server := api.NewApiServer(&cfg.App)
	server.Start()
}
