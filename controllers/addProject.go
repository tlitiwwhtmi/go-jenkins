package controllers

import (
	"go-jenkins/models/bd"
	"go-jenkins/service/database"
	"go-jenkins/utils"
	"time"

	"github.com/astaxie/beego"
)

type AddProject struct {
	beego.Controller
}

func (c *AddProject) Post() {
	if user := utils.GetUserFromSession(c.Ctx); user != nil {
		projectId, err := c.GetInt64("pid")
		if err != nil {
			utils.GetLogger().Error("ADD PROJECT FAILED: missing params projectId")
			c.Ctx.WriteString("param wrong")
			return
		}
		project := new(bdmodels.BdProject)
		project.CreateTime = time.Now()
		project.Language = c.GetString("lang")
		project.ProjectId = projectId
		project.UserAccount = user.Account
		_, err = database.GetProjectByAccountAndPid(project)
		if err == nil {
			c.Ctx.WriteString("already exist")
			return
		}
		_, err = database.AddProject(project)
		if err != nil {
			utils.GetLogger().Error("ADD PROJECT FAILED: save project error " + err.Error())
			c.Ctx.WriteString("save error")
			return
		}
		utils.GetLogger().Info("ADD PROJECT SUCCESS")
		c.Ctx.WriteString("success")
		return
	}
	c.Ctx.WriteString("not login")
}
