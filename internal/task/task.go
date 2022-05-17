package task

import (
	"fmt"
	"github.com/bianjieai/iobscan-explorer-backend/internal/common/util"
	"github.com/bianjieai/iobscan-explorer-backend/pkg/logger"
	"time"
)

func init() {
	RegisterTasks(&SyncDdcTask{})
}

type Task interface {
	Name() string
	Cron() int // second of Intervals
	Start()
	DoTask(fn func(string) chan bool) error
}

var (
	tasks []Task
)

func RegisterTasks(task ...Task) {
	tasks = append(tasks, task...)
}

// GetTasks get all the task
func GetTasks() []Task {
	return tasks
}

func Start() {
	if len(GetTasks()) == 0 {
		return
	}

	for _, one := range GetTasks() {
		var taskId = fmt.Sprintf("%s[%s]", one.Name(), util.FmtTime(time.Now(), util.DateFmtYYYYMMDD))
		logger.Info("timerTask begin to work", logger.String("taskId", taskId))
		go one.Start()
	}

}
