package controllers

import (
	"go-jenkins/service/database"
	"go-jenkins/utils"

	"github.com/astaxie/beego"
)

type OwnedProject struct {
	beego.Controller
}

func (c *OwnedProject) Get() {
	if user := utils.GetUserFromSession(c.Ctx); user != nil {
		count, projects, err := database.GetProjectsByAccount(user.Account)
		if err != nil {
			c.Ctx.WriteString("error")
			return
		}
		if count == 0 {
			c.Data["json"] = &projects
			c.ServeJSON()
			return
		}

		pqcon := database.NewPostGres(beego.AppConfig.String("postgressUrl"))
		gitProjects, err := pqcon.GetGitProjectDetail(projects)
		if err != nil {
			c.Ctx.WriteString("error")
			return
		}
		for _, gitPro := range *gitProjects {
			for i, _ := range *projects {
				project := &(*projects)[i]
				if gitPro.Id == project.ProjectId {
					project.ProjectName = gitPro.Namespace + "/" + gitPro.Name
				}
			}
		}
		//bdgitlab.RawProjects(projects)
		if err != nil {
			c.Ctx.WriteString("error")
			return
		}
		c.Data["json"] = &projects
		c.ServeJSON()
		return
	}
	c.Ctx.WriteString("error")
}
