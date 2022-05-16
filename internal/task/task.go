package task

import (
	"github.com/bianjieai/iobscan-explorer-backend/internal/common/configs"
	"github.com/robfig/cron"
)

var (
	confSrv *configs.Config
)

func Start(cfg *configs.Config) {
	confSrv = cfg
	// tasks manager by cron job
	c := cron.New()
	// add cronjob
	if cfg.Server.DdcCronjobExectime == "" {
		cfg.Server.DdcCronjobExectime = "1 */1 * * * ?"
	}
	c.AddFunc(cfg.Server.DdcCronjobExectime, func() {
		new(SyncDdcTask).Start()
	})
	c.Start()

}
