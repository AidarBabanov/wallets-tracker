package logging

import (
	"github.com/beego/beego/v2/core/logs"
	"os"
)

const (
	logsDir = "logs"
)

const fileConfig = `{"filename":"logs/logs.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":30,"color":true}`

func Init() {
	logs.SetLogFuncCallDepth(3)
	logs.EnableFuncCallDepth(true)

	err := logs.SetLogger(logs.AdapterConsole)
	if err != nil {
		logs.Error(err)
	}
	if _, err = os.Stat(logsDir); os.IsNotExist(err) {
		err = os.Mkdir(logsDir, 0777)
		if err != nil {
			logs.Error(err)
		}
	}
	err = logs.SetLogger(logs.AdapterFile, fileConfig)
	if err != nil {
		logs.Error(err)
	}
}
