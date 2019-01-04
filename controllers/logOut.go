package controllers

import (
	"github.com/astaxie/beego"
)

type LogOutController struct {
	beego.Controller
}

func (c *LogOutController) Get() {
	jobId, err := c.GetInt64("jobid")
	if err != nil {
		c.Ctx.WriteString("error")
		return
	}
	c.Data["jobId"] = jobId
	c.TplName = "log.html"
}
