package utils

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
)

var log *logs.BeeLogger

func InitLogger() {
	//调整日志输出缓存
	log = logs.NewLogger(1000)
	//使用文件输出方式
	logPath := beego.AppConfig.String("sys_log_path")+"go-jenkins.log"
	log.SetLogger("file", `{"filename":"`+logPath+`"}`)
}

func GetLogger() *logs.BeeLogger {
	if log == nil {
		InitLogger()
	}
	return log
}
