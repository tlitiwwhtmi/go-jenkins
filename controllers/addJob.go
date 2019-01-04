package controllers

import (
	"go-jenkins/models/bd"
	"go-jenkins/service/database"
	"go-jenkins/utils"
	"time"

	"github.com/astaxie/beego"
	"go-jenkins/service/schedule"
)

type AddJob struct {
	beego.Controller
}

func (c *AddJob) Post() {
	if user := utils.GetUserFromSession(c.Ctx); user != nil {
		projectId, err := c.GetInt64("pid")
		if err != nil {
			utils.GetLogger().Error("CREATE JOB FAILED: missing params projectId")
			c.Ctx.WriteString("param wrong")
			return
		}
		isGitLab, err := c.GetInt("gitlab")
		if err != nil {
			utils.GetLogger().Error("CREATE JOB FAILED: missing params isGitlab")
			c.Ctx.WriteString("param wrong")
			return
		}
		autoBuild, err := c.GetInt("autobuild")
		if err != nil {
			utils.GetLogger().Error("CREATE JOB FAILED: missing params autoBuild")
			c.Ctx.WriteString("param wrong")
		}
		deployType, err := c.GetInt("deploytype")
		if err != nil {
			utils.GetLogger().Error("CREATE JOB FAILED: missing params deployType")
			c.Ctx.WriteString("param wrong")
		}

		//check if is existed
		jobs, err := database.GetJobByPIdAndBranch(projectId, c.GetString("branch"))
		if c.GetString("build") == "1" {
			if err == nil {
				for _, tempJob := range *jobs {
					if tempJob.Build == c.GetString("build") && tempJob.PomPath == c.GetString("pompath") && tempJob.Deploy == deployType {
						c.Ctx.WriteString("already exist:" + tempJob.JobName)
						return
					}
				}
			}
		}
		if c.GetString("build") == "2" {
			if err == nil {
				for _, tempJob := range *jobs {
					if tempJob.Shell == c.GetString("shell") && tempJob.Build == c.GetString("build") && tempJob.PomPath == c.GetString("pompath") {
						c.Ctx.WriteString("already exist:" + tempJob.JobName)
						return
					}
				}
			}
		}

		job := new(bdmodels.BdJob)
		job.ProjectId = projectId
		job.BranchName = c.GetString("branch")
		job.JobName = c.GetString("name")
		job.IsGitlab = isGitLab
		job.Email = c.GetString("emails")
		job.Shell = c.GetString("shell")
		job.Build = c.GetString("build")
		job.PomPath = c.GetString("pompath")
		job.Profile = c.GetString("profile")
		job.Deploy = deployType
		job.AutoBuild = autoBuild
		job.AutoTime = c.GetString("autotime")
		job.CreateTime = time.Now()
		job.CreatorAccount = user.Account
		job.ModifyTime = time.Now()
		job.Modifier = user.Account
		jobId, err := database.AddJob(job)
		if err != nil {
			utils.GetLogger().Error("CREATE JOB FAILED: save job error  " + err.Error())
			c.Ctx.WriteString("save wrong")
		}
		job.Id = jobId
		utils.GetLogger().Info("CREATE JOB SUCCESS")
		schedule.ScheduledJob(job)
		c.Ctx.WriteString("success")
		return
	}
	c.Ctx.WriteString("not login")
}
