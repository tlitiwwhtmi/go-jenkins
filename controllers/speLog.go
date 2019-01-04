package controllers

import (
	"github.com/astaxie/beego"
)

type SpeLogController struct {
	beego.Controller
}

func (c *SpeLogController) Get() {
	hId, err := c.GetInt64("id")
	if err != nil {
		c.Ctx.WriteString("error")
		return
	}
	c.Data["hId"] = hId
	c.TplName = "spelog.html"
}
