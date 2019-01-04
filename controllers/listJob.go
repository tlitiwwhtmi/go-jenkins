package controllers

import (
	"fmt"
	"go-jenkins/models/bd"
	"go-jenkins/service/database"
	"go-jenkins/service/gitlab"
	"go-jenkins/utils"

	"github.com/astaxie/beego"
	"github.com/xanzy/go-gitlab"
)

type ListJob struct {
	beego.Controller
}

func (c *ListJob) Get() {
	if user := utils.GetUserFromSession(c.Ctx); user != nil {
		pageIndex, err := c.GetInt("page")
		if err != nil {
			fmt.Println("1;", err.Error())
			c.Ctx.WriteString("error")
			return
		}
		perPage, err := c.GetInt("size")
		if err != nil {
			fmt.Println("2;", err.Error())
			c.Ctx.WriteString("error")
			return
		}
		projectId, err := c.GetInt64("pid")
		if err != nil {
			fmt.Println("3;", err.Error())
			c.Ctx.WriteString("error")
			return
		}
		branchName := c.GetString("branch")
		keyword := c.GetString("keyword")
		total, jobs, err := database.GetJobsByProjectId(projectId, pageIndex, perPage, branchName, keyword)
		if err != nil {
			fmt.Println("4;", err.Error())
			c.Ctx.WriteString("error")
			return
		}
		branches, err := bdgitlab.GetBranchByProjectId(c.GetString("pid"))
		if err != nil {
			fmt.Println("5;", err.Error())
			c.Ctx.WriteString("error")
			return
		}
		//service.RawJob(jobs)
		//service.GetJobStatus(jobs)
		c.Data["json"] = &pageJobs{total, jobs, branches, c.GetString("pid"), keyword}
		c.ServeJSON()
		return
	}
	c.Ctx.WriteString("error")
}

type pageJobs struct {
	Total     int64
	Jobs      *[]bdmodels.BdJob
	Branches  []*gitlab.Branch
	ProjectId string
	Keyword   string
}
