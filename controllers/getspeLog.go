package controllers

import (
	"github.com/astaxie/beego"
	"go-jenkins/service/bdlog"
	"go-jenkins/service/database"
)

type GetspeLogController struct {
	beego.Controller
}

func (c *GetspeLogController) Get() {
	hId, err := c.GetInt64("hid")
	if err != nil {
		c.Ctx.WriteString("error")
		return
	}
	history, err := database.GetHistoryById(hId)
	if err != nil {
		c.Ctx.WriteString("error")
		return
	}
	log, err := bdlog.ReadLog(history.Log)
	if err != nil {
		c.Data["json"] = &spelogMsg{"", history.Status, ""}
	} else {
		c.Data["json"] = &spelogMsg{*log, history.Status, history.Log}
	}
	c.ServeJSON()
	return
}

type spelogMsg struct {
	Log      string
	OutCome  string
	FileName string
}
