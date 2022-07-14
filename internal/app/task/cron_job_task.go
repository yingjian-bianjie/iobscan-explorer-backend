package task

import (
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type Cron interface {
	Name() string
	Run()
	Cron() string
}

var cronList []Cron

func RegisterCron(cron ...Cron) {
	cronList = append(cronList, cron...)
}

func CronStart() {
	if len(cronList) == 0 {
		return
	}

	c := cron.New(cron.WithSeconds())
	for _, v := range cronList {
		task := v

		task.Run()
		_, err := c.AddFunc(task.Cron(), func() {
			logrus.Infof("======== cron %s start", task.Name())
			task.Run()
			logrus.Infof("======== cron %s end", task.Name())
		})
		if err != nil {
			logrus.Fatal("cron job err", err)
		}
	}
	c.Start()
}
