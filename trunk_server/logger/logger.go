package logger

import (
	"go-admin/config"

	"github.com/sohaha/zlsgo/zlog"
)

func OnInit() error {

	logPath := config.Config.Logger.Path
	if logPath == "" {
		config.Config.Logger.Path = config.WorkPath + "\\logs.log"
	}
	zlog.SetSaveFile(config.Config.Logger.Path, true)
	zlog.LogMaxDurationDate = config.Config.Logger.SaveDay
	return nil
}
