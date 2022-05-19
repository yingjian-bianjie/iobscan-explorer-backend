package constant

import "errors"

var (
	SkipErrmsgABIMethodNoFound  = errors.New("no found method in ABI data")
	SkipErrmsgABITypeNoFound    = errors.New("no found DdcType")
	SkipErrmsgNoSupport         = errors.New("no implement method parse code")
	SkipErrmsgNoSupportContract = errors.New("no support contract")
	ErrDbExist                  = errors.New("save failed for already exist")
)

const (
	SyncDdcTaskName     = "SyncDdcTaskName"
	CronTimeSyncDdcTask = 30
)
