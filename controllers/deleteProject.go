package controllers

import (
	"fmt"
	"go-jenkins/models/bd"
	"go-jenkins/service/database"
	"go-jenkins/utils"

	"github.com/astaxie/beego"
)

type DeleteProject struct {
	beego.Controller
}

func (c *DeleteProject) Get() {
	if user := utils.GetUserFromSession(c.Ctx); user != nil {
		projectId, err := c.GetInt64("pid")
		if err != nil {
			fmt.Println(err.Error())
			utils.GetLogger().Info("DELETE PROJECT FAILED:missing ")
			c.Ctx.WriteString("param wrong")
			return
		}
		project := new(bdmodels.BdProject)
		project.UserAccount = user.Account
		project.ProjectId = projectId
		pro, err := database.GetProjectByAccountAndPid(project)
		if err != nil {
			fmt.Println(err.Error())
			c.Ctx.WriteString("not exist")
			return
		}
		err = database.DeleteProject(pro)
		if err != nil {
			fmt.Println(err.Error())
			c.Ctx.WriteString("delete failed")
			return
		}
		c.Ctx.WriteString("success")
		return
	}
	c.Ctx.WriteString("not login")
}
