package controllers

import (
	"github.com/astaxie/beego"
	"go-jenkins/service"
)

type StartBuild struct {
	beego.Controller
}

func (c *StartBuild) Get() {
	jobId, err := c.GetInt64("id")
	if err != nil {
		c.Ctx.WriteString("error")
		return
	}
	service.StartBuild(jobId, nil)
	c.Ctx.WriteString("success")
}
