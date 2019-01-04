package controllers

import (
	"go-jenkins/service"

	"github.com/astaxie/beego"
	"go-jenkins/service/database"
)

type GetLogController struct {
	beego.Controller
}

func (c *GetLogController) Get() {
	jobId, err := c.GetInt64("jobid")
	if err != nil {
		c.Ctx.WriteString("error")
		return
	}
	log := ""
	outCome := ""
	isRunning := false
	needRead := true
	nLogName := ""
	logName := c.GetString("logname")
	history, err := database.GetLatestHistory(jobId)
	if history != nil {
		if logName == "" {
			log, outCome, isRunning = service.GetJobLog(jobId)
		} else {
			if logName == history.Log {
				needRead = false
			} else {
				log, outCome, isRunning = service.GetJobLog(jobId)
			}
		}
		nLogName = history.Log
	}

	c.Data["json"] = &logMsg{log, outCome, isRunning, needRead, nLogName}
	c.ServeJSON()
	return
}

type logMsg struct {
	Log       string
	OutCome   string
	IsRunning bool
	NeedRead  bool
	LogName   string
}
