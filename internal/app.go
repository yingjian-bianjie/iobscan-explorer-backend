package app

import (
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/api"
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/config"
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/constant"
	"github.com/bianjieai/iobscan-explorer-backend/internal/app/global"
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
	server := api.NewApiServer(&cfg.App)
	server.Start()
}
