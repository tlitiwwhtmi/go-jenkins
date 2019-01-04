package controllers

import (
	"github.com/astaxie/beego"
	"go-jenkins/service/database"
	"go-jenkins/utils"
)

type BuildsDay struct {
	beego.Controller
}

func (c *BuildsDay) Get() {
	if user := utils.GetUserFromSession(c.Ctx); user != nil {
		buildsDays, err := database.GetBuildsPerday()
		if err != nil {
			c.Ctx.WriteString("error")
			return
		}
		c.Data["json"] = buildsDays
		c.ServeJSON()
		return
	}
	c.Ctx.WriteString("error")
}
