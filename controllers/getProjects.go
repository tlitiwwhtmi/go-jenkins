package controllers

import (
	"go-jenkins/service/gitlab"
	"go-jenkins/utils"

	"github.com/astaxie/beego"
)

type GetProjects struct {
	beego.Controller
}

func (c *GetProjects) Get() {
	if user := utils.GetUserFromSession(c.Ctx); user != nil {
		projects := bdgitlab.GetProjectsByUser(user)
		c.Data["json"] = &projects
		c.ServeJSON()
	}
}
